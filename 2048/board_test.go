package main

import "testing"

func compare(b1, b2 *Board) bool {
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if b1.field[x][y] != b2.field[x][y] {
				return false
			}
		}
	}
	return true
}

func TestMoveUp(t *testing.T) {
	before := Board{[4][4]uint16{
		[4]uint16{2, 4, 4, 0},
		[4]uint16{2, 4, 4, 2},
		[4]uint16{2, 4, 8, 0},
		[4]uint16{2, 8, 0, 2},
	}}
	after := Board{[4][4]uint16{
		[4]uint16{4, 8, 8, 4},
		[4]uint16{4, 4, 8, 0},
		[4]uint16{0, 8, 0, 0},
		[4]uint16{0, 0, 0, 0},
	}}

	err := before.moveUp()
	if err != nil {
		t.Errorf("could not move up: %v", err)
	}
	if !compare(&before, &after) {
		t.Errorf("move not ok except %s, got %s", after, before)
	}
}

func TestColumnUpNotValid(t *testing.T) {
	b := Board{[4][4]uint16{
		[4]uint16{2, 4, 4, 0},
		[4]uint16{0, 0, 0, 0},
		[4]uint16{0, 0, 0, 0},
		[4]uint16{0, 0, 0, 0},
	}}
	org := b
	err := b.Move(UP)
	if err == nil {
		t.Errorf("expect error moving up board %s, got no error, got %s", org, b)
	}
}

func TestMoveDown(t *testing.T) {
	before := Board{[4][4]uint16{
		[4]uint16{2, 4, 4, 0},
		[4]uint16{2, 4, 4, 2},
		[4]uint16{2, 4, 8, 0},
		[4]uint16{2, 8, 0, 2},
	}}
	after := Board{[4][4]uint16{
		[4]uint16{0, 0, 0, 0},
		[4]uint16{0, 4, 0, 0},
		[4]uint16{4, 8, 8, 0},
		[4]uint16{4, 8, 8, 4},
	}}

	err := before.moveDown()
	if err != nil {
		t.Errorf("could not move down: %v", err)
	}
	if !compare(&before, &after) {
		t.Errorf("move not ok except %s, got %s", after, before)
	}
}

func TestColumnDownNotValid(t *testing.T) {
	b := Board{[4][4]uint16{
		[4]uint16{2, 4, 4, 0},
		[4]uint16{0, 0, 0, 0},
		[4]uint16{0, 0, 0, 0},
		[4]uint16{0, 0, 0, 0},
	}}
	org := b
	err := b.Move(DOWN)
	if err == nil {
		t.Errorf("expect error moving down board %s, got no error, got %s", org, b)
	}
}

func TestMoveRight(t *testing.T) {
	before := Board{[4][4]uint16{
		[4]uint16{2, 2, 2, 2},
		[4]uint16{8, 4, 4, 4},
		[4]uint16{0, 8, 4, 4},
		[4]uint16{2, 0, 2, 0},
	}}
	after := Board{[4][4]uint16{
		[4]uint16{0, 0, 4, 4},
		[4]uint16{0, 8, 4, 8},
		[4]uint16{0, 0, 8, 8},
		[4]uint16{0, 0, 0, 4},
	}}

	err := before.moveRight()
	if err != nil {
		t.Errorf("could not move right: %v", err)
	}
	if !compare(&before, &after) {
		t.Errorf("move not ok except %s, got %s", after, before)
	}
}

func TestMoveLeft(t *testing.T) {
	before := Board{[4][4]uint16{
		[4]uint16{2, 2, 2, 2},
		[4]uint16{8, 4, 4, 4},
		[4]uint16{0, 8, 4, 4},
		[4]uint16{2, 0, 2, 0},
	}}
	after := Board{[4][4]uint16{
		[4]uint16{4, 4, 0, 0},
		[4]uint16{8, 8, 4, 0},
		[4]uint16{8, 8, 0, 0},
		[4]uint16{4, 0, 0, 0},
	}}

	err := before.moveLeft()
	if err != nil {
		t.Errorf("could not move left: %v", err)
	}
	if !compare(&before, &after) {
		t.Errorf("move not ok except %s, got %s", after, before)
	}
}
