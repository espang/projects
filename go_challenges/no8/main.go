package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	f := Field{}

	// read the sudoku
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		err := f.AddLine(line)
		if err != nil {
			// wrong input -> crash
			panic(err)
		}
	}

	fmt.Println("Start field:")
	fmt.Println(f)

	// solve the sudoku
	solved, err := f.Solve()
	if err != nil {
		panic(err)
	}

	if !solved.Valid(false) {
		fmt.Println("Got invalid solution!")
		fmt.Println(solved)

	} else {
		fmt.Println("Solution:")
		fmt.Println(solved)
	}
}
