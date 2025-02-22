package main

import (
	"fmt"
	"log"

	"github.com/adrianosela/adventofcode/utils/grid"
	"github.com/adrianosela/adventofcode/utils/set"
)

func main() {
	smallSampleSoln, err := solve("sample-input-small.txt")
	if err != nil {
		log.Fatalf("failed to solve for small sample input: %v", err)
	}
	log.Printf("[Answer to Small Sample in Part 1] The price is %d (should be 140)", smallSampleSoln)

	largeSampleSoln, err := solve("sample-input-large.txt")
	if err != nil {
		log.Fatalf("failed to solve for large sample input: %v", err)
	}
	log.Printf("[Answer to Large Sample in Part 1] The price is %d (should be 1930)", largeSampleSoln)

	inputSoln, err := solve("input.txt")
	if err != nil {
		log.Fatalf("failed to solve for input: %v", err)
	}
	log.Printf("[Answer to Part 1] The price is %d", inputSoln)
}

func solve(filename string) (int, error) {
	g, err := grid.LoadByte(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to load byte grid from file: %v", err)
	}
	memo := set.New[grid.Coordinate]()
	sum := 0
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			coord := grid.Coordinate{Y: y, X: x}
			if !memo.Has(coord) {
				area, perim := measureRegion(g, memo, g[y][x], coord)
				sum += (area * perim)
			}
		}
	}
	return sum, nil
}

func measureRegion(g grid.Grid[byte], memo set.Set[grid.Coordinate], flower byte, coord grid.Coordinate) (int, int) {
	// if the current coordinate exceeds the bounds of the grid, return
	if coord.X < 0 || coord.X >= len(g[coord.Y]) || coord.Y < 0 || coord.Y >= len(g) {
		return 0, 0
	}
	// if the current coordinate does not match the current flower, return
	if g[coord.Y][coord.X] != flower {
		return 0, 0
	}
	// if the current coordiante was already visited, return
	if memo.Has(coord) {
		return 0, 0
	}
	// this is a flower, visit it
	memo.Put(coord)
	area := 1
	sides := 0
	for _, n := range []grid.Coordinate{
		{Y: coord.Y + 1, X: coord.X},
		{Y: coord.Y - 1, X: coord.X},
		{Y: coord.Y, X: coord.X + 1},
		{Y: coord.Y, X: coord.X - 1},
	} {
		// add a side for every "neighbor" that exceeds the grid's bounds
		if n.Y < 0 || n.Y >= len(g) || n.X < 0 || n.X >= len(g[n.Y]) {
			sides++
			continue
		}
		// add a side for every neighbor that is not of the current flower
		if g[n.Y][n.X] != flower {
			sides++
			continue
		}
		// don't count already visited neighbors
		if memo.Has(n) {
			continue
		}
		// visit and measure neighbor
		innerArea, innerSides := measureRegion(g, memo, flower, n)
		area += innerArea
		sides += innerSides
	}
	return area, sides
}
