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

var directions = [...]string{
	"up",
	"down",
	"left",
	"right",
}

func (dir Direction) String() string {
	return directions[dir]
}

// Board represents the 2048-4x4-board
// with numbers between 2 and 2**16
type Board struct {
	field [4][4]uint16
}

func (b Board) Equals(cmp Board) bool {
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if b.field[x][y] != cmp.field[x][y] {
				return false
			}
		}
	}
	return true
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

func (b *Board) Max() uint16 {
	var max uint16 = 0

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if b.field[row][col] > max {
				max = b.field[row][col]
			}
		}
	}
	return max
}

func (b *Board) Won() bool {
	return b.Max() >= 2048
}

func (b *Board) Lost() bool {
	//not lost when board has at least one zero
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if b.field[r][c] == 0 {
				return false
			}
		}
	}

	// need two same numbers in
	for i := 0; i < 4; i++ {
		// 4 2 2 4
		//column i:
		if b.field[0][i] == b.field[1][i] ||
			b.field[1][i] == b.field[2][i] ||
			b.field[2][i] == b.field[3][i] {
			return false
		}
		//row i:
		if b.field[i][0] == b.field[i][1] ||
			b.field[i][1] == b.field[i][2] ||
			b.field[i][2] == b.field[i][3] {
			return false
		}
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
	cpy := *b
	for col := 0; col < 4; col++ {
		b.handleColumn(col, UP)
	}
	if b.Equals(cpy) {
		return errors.New("Move up not valid")
	}
	return nil
}

func (b *Board) moveDown() error {
	cpy := *b
	for col := 0; col < 4; col++ {
		b.handleColumn(col, DOWN)
	}
	if b.Equals(cpy) {
		return errors.New("Move up not valid")
	}
	return nil
}

func (b *Board) handleColumn(col int, dir Direction) {
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
}

func (b *Board) moveLeft() error {
	cpy := *b
	for row := 0; row < 4; row++ {
		b.handleRow(row, LEFT)
	}
	if b.Equals(cpy) {
		return errors.New("Move up not valid")
	}
	return nil
}

func (b *Board) moveRight() error {
	cpy := *b
	for row := 0; row < 4; row++ {
		b.handleRow(row, RIGHT)
	}
	if b.Equals(cpy) {
		return errors.New("Move up not valid")
	}
	return nil
}

func (b *Board) handleRow(row int, dir Direction) {
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
}
