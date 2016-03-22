package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	repeat = flag.Int("repeat", 1e4, "how many games to play")
	seed   = flag.Int64("seed", 0, "Random seed, 0 to choose random")

	//random  = flag.Bool("random", false, "test with random solver")
	pattern = flag.String("pattern", "UUDLLR", "repeating pattern: (U)p, (D)own, (L)eft and (R)ight")
)

type Score uint16

const (
	S0 Score = 1 << iota
	S2
	S4
	S8
	S16
	S32
	S64
	S128
	S256
	S512
	S1024
	S2048_AND_MORE
)

var sortedScores = [...]Score{
	S0, S2, S4, S8, S16, S32, S64, S128, S256, S512, S1024, S2048_AND_MORE,
}

type Stats struct {
	count int
	hist  map[Score]int
	won   int
	times []time.Duration
}

func (s *Stats) Add(max uint16, won bool, dur time.Duration) {
	s.count++
	if max > uint16(S2048_AND_MORE) {
		s.hist[S2048_AND_MORE]++
	} else {
		s.hist[Score(max)]++
	}
	if won {
		s.won++
	}
	s.times = append(s.times, dur)
}

func (s *Stats) Print() {
	digits := math.Ceil(math.Log10(float64(s.count)))
	fmt.Println(digits, " digits")

	fmt.Println("Played ", s.count, " games")
	fmt.Println("From which ", s.won, " were won")
	var duration time.Duration
	for _, dur := range s.times {
		duration += dur
	}
	fmt.Println("Took ", duration, " or ", duration/time.Duration(s.count),
		"per game")

	maxV := 0
	for _, v := range s.hist {
		if v > maxV {
			maxV = v
		}
	}

	factor := 30 / float64(maxV)

	countPattern := "%" + fmt.Sprintf("%d", int(digits)) + "d: "
	fmt.Println("Score|Count")
	for _, score := range sortedScores {
		fmt.Printf("%5d|", score)
		fmt.Printf(countPattern, s.hist[score])
		points := int(math.Ceil(float64(s.hist[score]) * factor))
		for i := 0; i < points; i++ {
			fmt.Print("*")
		}
		fmt.Println("")
	}

}

func play(times int) {
	stats := &Stats{
		hist:  map[Score]int{},
		times: []time.Duration{},
	}

	var solver Solver
	solver = NewPatternSolver(*pattern)
	var t time.Time
	for i := 0; i < times; i++ {
		b := StartBoard()

		var newDir Direction
		var err error

		lost := b.Lost()
		won := b.Won()

		t = time.Now()

		for won || !lost {
			newDir = solver.NextStep(b)
			err = b.Move(newDir)
			if err != nil {
				continue
			}
			lost = b.Lost()
			won = b.Won()
		}
		stats.Add(b.Max(), won, time.Since(t))
	}

	stats.Print()
}

func main() {
	flag.Parse()
	if *seed == 0 {
		rand.Seed(time.Now().Unix())
	} else {
		rand.Seed(*seed)
	}

	play(*repeat)
}
