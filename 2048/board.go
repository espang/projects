package main

import (
	"errors"
	"fmt"
	"math/rand"
)

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

// Board represents the 2048-4x4-board
// with numbers between 2 and 2**16
type Board struct {
	field [4][4]uint16
}

func StartBoard() *Board {
	b := &Board{}
	b.add()
	b.add()
	return b
}

func (b Board) String() string {
	return fmt.Sprint(b.field)
}

func (b Board) Print() {

	for row := 0; row < 4; row++ {

		for col := 0; col < 4; col++ {
			fmt.Printf("%5d ", b.field[row][col])
		}
		fmt.Println()

	}

}

func (b *Board) Lost() bool {
	cpy := *b

	err := cpy.moveUp()
	if err == nil {
		return false
	}

	cpy = *b
	err = cpy.moveDown()
	if err == nil {
		return false
	}

	cpy = *b
	err = cpy.moveRight()
	if err == nil {
		return false
	}

	cpy = *b
	err = cpy.moveLeft()
	if err == nil {
		return false
	}

	return true
}

func (b *Board) add() error {
	var newValue uint16 = 2
	if rand.Float32() < 0.1 {
		newValue = 4
	}

	zeroIndizes := []int{}
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if b.field[row][col] == 0 {
				zeroIndizes = append(zeroIndizes, row*10+col)
			}
		}
	}

	randIndex := rand.Intn(len(zeroIndizes))
	value := zeroIndizes[randIndex]
	row := value / 10
	col := value % 10

	b.field[row][col] = newValue
	return nil
}

func (b *Board) Move(dir Direction) error {
	var err error
	switch dir {
	case UP:
		err = b.moveUp()
	case DOWN:
		err = b.moveDown()
	case LEFT:
		err = b.moveLeft()
	case RIGHT:
		err = b.moveRight()
	}
	if err != nil {
		return err
	}
	return b.add()
}

// from up to down
//
//
func (b *Board) moveUp() error {
	hasMoved := false
	for col := 0; col < 4; col++ {
		err := b.handleColumn(col, UP)
		if err == nil {
			hasMoved = true
		}
	}
	if hasMoved {
		return nil
	}
	return errors.New("Move up not valid")
}

func (b *Board) moveDown() error {
	hasMoved := false
	for col := 0; col < 4; col++ {
		err := b.handleColumn(col, DOWN)
		if err == nil {
			hasMoved = true
		}
	}
	if hasMoved {
		return nil
	}
	return errors.New("Move down not valid")
}

func (b *Board) handleColumn(col int, dir Direction) error {
	getRowIndex := func(row int) int {
		return row
	}
	if dir == DOWN {
		getRowIndex = func(row int) int {
			return (3 - row) % 4
		}
	}

	nonZeros := []uint16{}
	var row int
	for i := 0; i < 4; i++ {
		row = getRowIndex(i)
		if b.field[row][col] != 0 {
			nonZeros = append(nonZeros, b.field[row][col])
		}
	}
	result := []uint16{}
	for i := 0; i < len(nonZeros); i++ {
		if i == len(nonZeros)-1 {
			//last element cannot be merged
			result = append(result, nonZeros[i])
			break
		}
		if nonZeros[i] == nonZeros[i+1] {
			//merge
			result = append(result, nonZeros[i]<<1)
			i += 1
		} else {
			result = append(result, nonZeros[i])
		}
	}

	for i := 0; i < 4; i++ {
		row = getRowIndex(i)
		if i < len(result) {
			b.field[row][col] = result[i]
		} else {
			b.field[row][col] = 0
		}
	}

	if len(result) == len(nonZeros) {
		return errors.New("No change in column")
	}
	return nil
}

func (b *Board) moveLeft() error {
	hasMoved := false
	for row := 0; row < 4; row++ {
		err := b.handleRow(row, LEFT)
		if err == nil {
			hasMoved = true
		}
	}
	if hasMoved {
		return nil
	}
	return errors.New("Move left not valid")
}

func (b *Board) moveRight() error {
	hasMoved := false
	for row := 0; row < 4; row++ {
		err := b.handleRow(row, RIGHT)
		if err == nil {
			hasMoved = true
		}
	}
	if hasMoved {
		return nil
	}
	return errors.New("Move right not valid")
}

func (b *Board) handleRow(row int, dir Direction) error {
	getColIndex := func(col int) int {
		return col
	}
	if dir == RIGHT {
		getColIndex = func(col int) int {
			return (3 - col) % 4
		}
	}

	nonZeros := []uint16{}
	var col int
	for i := 0; i < 4; i++ {
		col = getColIndex(i)
		if b.field[row][col] != 0 {
			nonZeros = append(nonZeros, b.field[row][col])
		}
	}
	result := []uint16{}
	for i := 0; i < len(nonZeros); i++ {
		if i == len(nonZeros)-1 {
			//last element cannot be merged
			result = append(result, nonZeros[i])
			break
		}
		if nonZeros[i] == nonZeros[i+1] {
			//merge
			result = append(result, nonZeros[i]<<1)
			i += 1
		} else {
			result = append(result, nonZeros[i])
		}
	}

	for i := 0; i < 4; i++ {
		col = getColIndex(i)
		if i < len(result) {
			b.field[row][col] = result[i]
		} else {
			b.field[row][col] = 0
		}
	}

	if len(result) == len(nonZeros) {
		return errors.New("No change in column")
	}
	return nil
}
