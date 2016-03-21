package main

import "math/rand"

type randomSolver struct {
	steps [4]Direction
}

func NewRandomSolver() *randomSolver {
	r := &randomSolver{
		steps: [4]Direction{UP, DOWN, LEFT, RIGHT},
	}
	return r
}

func (r *randomSolver) NextStep(b *Board) Direction {
	return r.steps[rand.Intn(4)]
}
