package grid

import "fmt"

type Grid[T any] [][]T

func New[T any]() Grid[T] {
	return make(Grid[T], 0)
}

func (g Grid[T]) String() string {
	s := ""
	for y := 0; y < len(g); y++ {
		s = fmt.Sprintf("%s%v\n", s, g[y])
	}
	return s
}
