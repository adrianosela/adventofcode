package grid

import "fmt"

type Coordinate struct {
	X int
	Y int
}

func (c *Coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.Y, c.X)
}

func (c *Coordinate) Equal(c2 *Coordinate) bool {
	return c.X == c2.X && c.Y == c2.Y
}
