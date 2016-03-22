package main

type Solver interface {
	NextStep(*Board) Direction
}
