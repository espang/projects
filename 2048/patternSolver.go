package main

type patternSolver struct {
	pattern []Direction
	head    int
}

func NewPatternSolver(pattern string) *patternSolver {

	runeSlice := []rune(pattern)
	directions := []Direction{}
	for _, r := range runeSlice {
		switch r {
		case 'U':
			directions = append(directions, UP)
		case 'D':
			directions = append(directions, DOWN)
		case 'L':
			directions = append(directions, LEFT)
		case 'R':
			directions = append(directions, RIGHT)
		}
	}

	s := &patternSolver{
		pattern: directions,
		head:    0,
	}
	return s
}

func (s *patternSolver) NextStep(b *Board) Direction {
	// get next direction
	dir := s.pattern[s.head]
	// increase the head
	s.head = (s.head + 1) % len(s.pattern)
	return dir
}
