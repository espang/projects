package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Field struct {
	field     [9][9]int
	lineCount int
}

const empty = 0

func (f *Field) AddLine(line string) error {
	if f.lineCount == 9 {
		return fmt.Errorf("already 9 lines added, wrong input format")
	}

	numbers := strings.Split(line, " ")

	if len(numbers) != 9 {
		return fmt.Errorf(
			"%d and not 9 numbers in line '%s'",
			len(numbers),
			line,
		)
	}

	for i, number := range numbers {
		if number == "_" {
			f.field[f.lineCount][i] = empty
			continue
		}
		n, err := strconv.Atoi(number)
		if err != nil {
			return err
		}
		f.field[f.lineCount][i] = n
	}
	f.lineCount++
	return nil
}

func (f Field) String() string {
	var buffer bytes.Buffer
	for _, row := range f.field {
		for _, cell := range row {
			if cell == empty {
				buffer.WriteString("_ ")
				continue
			}
			buffer.WriteString(
				fmt.Sprintf("%d ", cell),
			)
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (f Field) Valid(ignoreEmpty bool) bool {
	var val int
	// check rows and cols
	for i := 0; i < 9; i++ {
		rowVals := map[int]bool{}
		colVals := map[int]bool{}
		blockVals := map[int]bool{}

		// 0,0 | 0,3 | 0,6
		// 3,0 | 3,3 | 3,6
		// 6,0 | 6,3 | 6,6
		blockRow := 3 * (i / 3)
		blockCol := 3 * (i % 3)

		for j := 0; j < 9; j++ {
			val = f.field[i][j]
			if (val == empty && !ignoreEmpty) || rowVals[val] {
				return false
			}
			rowVals[val] = true

			val = f.field[j][i]
			if (val == empty && !ignoreEmpty) || colVals[val] {
				return false
			}
			colVals[val] = true

			val = f.field[blockRow+j%3][blockCol+j/3]
			if (val == empty && !ignoreEmpty) || blockVals[val] {
				return false
			}
			blockVals[val] = true
		}
	}

	return true
}

func (f Field) Possible(irow, icol int) []int {
	notPossible := map[int]bool{}
	var val int

	//block:
	// row 0, 1, 2, --> 0
	// row 3, 4, 5, --> 3
	// row 6, 7, 8, --> 6
	blockRow := 3 * (irow / 3)
	blockCol := 3 * (icol / 3)

	for i := 0; i <= 8; i++ {
		//row:
		val = f.field[i][icol]
		notPossible[val] = true
		//column:
		val = f.field[irow][i]
		notPossible[val] = true
		//block:
		val = f.field[blockRow+i%3][blockCol+i/3]
		notPossible[val] = true
	}

	result := []int{}
	for i := 1; i <= 9; i++ {
		if notPossible[i] {
			continue
		}
		result = append(result, i)
	}
	return result
}

func (f Field) Solve() (Field, error) {
	if !f.Valid(true) {
		return f, errors.New("Invalid solution")
	}

	// get all empty cells and there possible numbers
	emptyMap := map[[2]int][]int{}

	for irow, row := range f.field {
		for icol, cell := range row {
			if cell == empty {
				emptyMap[[2]int{irow, icol}] = []int{}
			}
		}
	}

	if len(emptyMap) == 0 {
		return f, nil
	}

	for key := range emptyMap {
		irow := key[0]
		icol := key[1]
		emptyMap[key] = f.Possible(irow, icol)
	}

	changed := false
	for key, values := range emptyMap {
		if len(values) == 1 {
			f.field[key[0]][key[1]] = values[0]
			changed = true
		}
		if len(values) == 0 {
			return f, errors.New("Not solvable")
		}
	}
	if changed {
		return f.Solve()
	}

	// not changed --> guess!
	k := [2]int{}
	length := 10
	for key, values := range emptyMap {
		if len(values) < length {
			length = len(values)
			k = key
		}
	}
	fmt.Println("Need to guess")
	fmt.Println(f)
	for _, guess := range emptyMap[k] {
		f.field[k[0]][k[1]] = guess
		solution, err := f.Solve()
		if err == nil {
			return solution, nil
		}
	}
	return f, errors.New("Not solvable yet")
}
