package main

import (
	"errors"
	"sync"
	"time"
)

// A repacker repacks trucks.
type repacker struct {
	boxes   chan box
	pallets chan pallet
	trucks  chan *truck
}

const areaPallet = 16

func (r *repacker) handleIncomings(in <-chan *truck) {
	for t := range in {
		r.trucks <- r.handleIncomingTruck(t)
	}
	close(r.boxes)
	close(r.trucks)
}

func (r *repacker) handleIncomingTruck(t *truck) (out *truck) {
	out = &truck{id: t.id}
	var area uint8
	for _, p := range t.pallets {
		area = 0
		for _, b := range p.boxes {
			area += b.w * b.l
		}
		if area == areaPallet && p.IsValid() == nil {
			out.pallets = append(out.pallets, p)
			continue
		}
		//remove all boxes from the pallet
		for _, b := range p.boxes {
			r.boxes <- b
		}
	}
	return out
}

type BoxRack struct {
	sync.Mutex
	content   map[uint8]map[uint8][]box
	boxes     int
	boxesDone bool
	done      bool
}

func NewBoxRack() *BoxRack {
	br := &BoxRack{}
	br.content = map[uint8]map[uint8][]box{}

	// 1,1 | 2,1 | 2,2 | 3,1 | 3,2 | 3,3 |
	// 4,1 | 4,2 | 4,3 | 4,4
	var w uint8 = 1
	var l uint8 = 1
	for ; w <= palletWidth; w++ {
		for ; l <= w; l++ {
			br.content[w] = map[uint8][]box{l: []box{}}
		}
	}
	return br
}

func (b *BoxRack) Add(b1 box) {
	b.Lock()
	defer b.Unlock()

	w, l := b1.w, b1.l
	if w < l {
		w, l = l, w
	}
	b.content[w][l] = append(b.content[w][l], b1)
	b.boxes++
}

func (b *BoxRack) BoxesDone() {
	b.Lock()
	defer b.Unlock()

	b.boxesDone = true
}

func (b *BoxRack) pop(w, l uint8) box {
	var x box
	x, b.content[w][l] = b.content[w][l][len(b.content[w][l])-1], b.content[w][l][:len(b.content[w][l])-1]
	b.boxes--
	return x
}

func (b *BoxRack) Pack() (pallet, error) {
	b.Lock()
	defer b.Unlock()

	if len(b.content[4][4]) > 0 {
		return pallet{boxes: []box{
			b.pop(4, 4),
		}}, nil
	}

	if len(b.content[4][3]) > 0 {
		if len(b.content[4][1]) > 0 {
			return pallet{boxes: []box{
				b.pop(4, 3),
				b.pop(4, 1),
			}}, nil
		}
		if len(b.content[2][1]) > 1 {
			return pallet{boxes: []box{
				b.pop(4, 3),
				b.pop(2, 1),
				b.pop(2, 1),
			}}, nil
		}
		if len(b.content[1][1]) > 0 && len(b.content[3][1]) > 0 {
			return pallet{boxes: []box{
				b.pop(4, 3),
				b.pop(3, 1),
				b.pop(1, 1),
			}}, nil
		}
	}

	if b.boxes == 0 && b.boxesDone {
		b.done = true
	}

	if b.boxesDone {
		for w, m1 := range b.content {
			for l, boxes := range m1 {
				if len(boxes) > 0 {
					return pallet{boxes: []box{b.pop(w, l)}}, nil
				}
			}
		}
	}
	return pallet{}, errors.New("nothing to pack")
}

func (r *repacker) repackPallets() {
	rack := NewBoxRack()

	go func() {
		for b := range r.boxes {
			rack.Add(b)
		}
		rack.BoxesDone()

	}()

	for !rack.done {
		p, err := rack.Pack()
		if err != nil {
			//nothing to pack...
			time.Sleep(50 * time.Millisecond)
			continue
		}
		r.pallets <- p
	}

	close(r.pallets)

}

func (r *repacker) handleOutgoings(out chan<- *truck) {
	//block till first trucks arrives
	current := <-r.trucks

	trucksDone := false
	palletsDone := false

	for !trucksDone && !palletsDone {
		select {
		case t, ok := <-r.trucks:
			if !ok {
				trucksDone = true
				continue
			}
			out <- current
			current = t
		case p, ok := <-r.pallets:
			if !ok {
				palletsDone = true
				continue
			}
			current.pallets = append(current.pallets, p)
		}
	}
	//send last truck with all pallets
	out <- current
	time.Sleep(50 * time.Millisecond)
	// The repacker must close channel out after it detects that
	// channel in is closed so that the driver program will finish
	// and print the stats.
	close(out)
}

func newRepacker(in <-chan *truck, out chan<- *truck) *repacker {
	repacker := &repacker{
		make(chan box, 1000),
		make(chan pallet, 1000),
		make(chan *truck, 1000),
	}

	go repacker.handleIncomings(in)
	go repacker.repackPallets()
	go repacker.handleOutgoings(out)

	return repacker
}
