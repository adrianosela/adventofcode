package grid

import "fmt"

type Coordinate struct {
	X int
	Y int
}

func (c *Coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.Y, c.X)
}
