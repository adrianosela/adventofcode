package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/adrianosela/adventofcode/utils/grid"
	"github.com/adrianosela/adventofcode/utils/set"
)

const (
	guardIndicator    = '^'
	spaceIndicator    = '.'
	obstacleIndicator = '#'

	directionUp    = "UP"
	directionDown  = "DOWN"
	directionLeft  = "LEFT"
	directionRight = "RIGHT"
)

var (
	dirToMovement = map[string]grid.Coordinate{
		directionUp:    {X: 0, Y: -1},
		directionDown:  {X: 0, Y: 1},
		directionLeft:  {X: -1, Y: 0},
		directionRight: {X: 1, Y: 0},
	}
	dirToNextDir = map[string]string{
		directionUp:    directionRight,
		directionDown:  directionLeft,
		directionLeft:  directionUp,
		directionRight: directionDown,
	}
)

func main() {
	filename := flag.String("filename", "", "The path to the input file")
	flag.Parse()

	locations, err := part1(*filename)
	if err != nil {
		log.Fatalf("failed to solve part 1: %v", err)
	}
	log.Printf("[Answer to Part 1] The number of unique visited coordinates is %d", locations)

	locations, err = part2(*filename)
	if err != nil {
		log.Fatalf("failed to solve part 2: %v", err)
	}
	log.Printf("[Answer to Part 2] The number of unique locations to place obstruction is %d", locations)
}

func findGuard(g grid.Grid[byte]) (*grid.Coordinate, bool) {
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			if g[y][x] == guardIndicator {
				return &grid.Coordinate{Y: y, X: x}, true
			}
		}
	}
	return nil, false
}

func part1(filename string) (int, error) {
	g, err := grid.LoadByte(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to load grid: %v", err)
	}

	guardPosition, ok := findGuard(g)
	if !ok {
		return 0, fmt.Errorf("grid did not contain guard indicator (%c)", guardIndicator)
	}

	visited := set.New(coordsToKey(guardPosition.Y, guardPosition.X))
	walk(g, visited, guardPosition.Y, guardPosition.X, directionUp)
	return visited.Size(), nil
}

func part2(filename string) (int, error) {
	g, err := grid.LoadByte(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to load grid: %v", err)
	}

	guardPosition, ok := findGuard(g)
	if !ok {
		return 0, fmt.Errorf("grid did not contain guard indicator (%c)", guardIndicator)
	}

	sum := 0
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			// this is the guard position, can't put it here
			if y == guardPosition.Y && x == guardPosition.X {
				continue
			}
			// this is already an obstacle, can't put it here
			if g[y][x] == obstacleIndicator {
				continue
			}

			// put obstacle
			g[y][x] = obstacleIndicator
			// count
			sum += count(
				g,
				set.New(coordsAndDirToKey(guardPosition.Y, guardPosition.X, directionUp)),
				guardPosition.Y,
				guardPosition.X,
				directionUp,
			)
			// remove obstacle
			g[y][x] = spaceIndicator
		}
	}
	return sum, nil
}

func walk(
	g grid.Grid[byte],
	visited set.Set[string],
	y int,
	x int,
	dir string,
) {
	movement := dirToMovement[dir]
	newY := y + movement.Y
	newX := x + movement.X
	if newY < 0 || newY >= len(g) || newX < 0 || newX >= len(g[y]) {
		return
	}
	if g[newY][newX] == obstacleIndicator {
		// keep the same coordinates, just change direction
		walk(g, visited, y, x, dirToNextDir[dir])
		return
	}
	visited.Put(coordsToKey(newY, newX))
	walk(g, visited, newY, newX, dir)
}

func coordsToKey(y int, x int) string {
	return fmt.Sprintf("(%d,%d)", y, x)
}

func count(
	g grid.Grid[byte],
	visited set.Set[string],
	y int,
	x int,
	dir string,
) int {
	movement := dirToMovement[dir]
	newY := y + movement.Y
	newX := x + movement.X
	if newY < 0 || newY >= len(g) || newX < 0 || newX >= len(g[y]) {
		return 0
	}

	if g[newY][newX] == obstacleIndicator {
		// turn direction
		return count(g, visited, y, x, dirToNextDir[dir])
	}

	key := coordsAndDirToKey(newY, newX, dir)
	if visited.Has(key) {
		// already visited, so this counts as one!
		return 1
	}
	visited.Put(key)
	return count(g, visited, newY, newX, dir)
}

func coordsAndDirToKey(y int, x int, dir string) string {
	return fmt.Sprintf("%s(%d,%d)", dir, y, x)
}
