package grid

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/adrianosela/adventofcode/utils/sliceconv"
)

type Grid[T any] [][]T

func New[T any]() Grid[T] {
	return make(Grid[T], 0)
}

func LoadByte(filename string) (Grid[byte], error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open grid file: %v", err)
	}
	grid := New[byte]()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file contents: %v", err)
	}
	return grid, nil
}

func LoadInt(filename string, delim string) (Grid[int], error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open grid file: %v", err)
	}
	grid := New[int]()
	scanner := bufio.NewScanner(file)
	for lineNo := 0; scanner.Scan(); lineNo++ {
		slice, err := sliceconv.StringsToInts(strings.Split(scanner.Text(), delim))
		if err != nil {
			return nil, fmt.Errorf("failed to convert line %d to integers slice: %v", lineNo, err)
		}
		grid = append(grid, slice)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file contents: %v", err)
	}
	return grid, nil
}

func (g Grid[T]) String() string {
	s := ""
	for y := 0; y < len(g); y++ {
		s = fmt.Sprintf("%s%v\n", s, g[y])
	}
	return s
}
