package grid

import "fmt"

type Grid[T any] [][]T

type Coordinate struct {
	X int
	Y int
}

func (g Grid[T]) String() string {
	s := ""
	for y := 0; y < len(g); y++ {
		s = fmt.Sprintf("%s%v\n", s, g[y])
	}
	return s
}
