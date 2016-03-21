package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	b := StartBoard()
	b.Print()

	var solver Solver
	solver = NewRandomSolver()

	var newDir Direction
	var err error

	lost := b.Lost()
	won := b.Won()

	for won || !lost {
		newDir = solver.NextStep(b)
		err = b.Move(newDir)
		if err != nil {
			continue
		}
		lost = b.Lost()
		won = b.Won()
	}

	fmt.Printf("Done. Max=%d, Lost=%v, Won=%v\n", b.Max(), lost, won)
}
