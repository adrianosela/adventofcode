package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/adrianosela/adventofcode/utils/grid"
)

const (
	moveLeft  = '<'
	moveRight = '>'
	moveUp    = '^'
	moveDown  = 'v'

	indicatorRobot = '@'
	indicatorWall  = '#'
	indicatorEmpty = '.'
	indicatorBox   = 'O'
)

func main() {
	smallSampleSoln, err := solve("sample-input-small.txt")
	if err != nil {
		log.Fatalf("failed to solve for small sample input: %v", err)
	}
	log.Printf("[Answer to Small Sample in Part 1] The sum is %d (should be 2028)", smallSampleSoln)

	largeSampleSoln, err := solve("sample-input-large.txt")
	if err != nil {
		log.Fatalf("failed to solve for large sample input: %v", err)
	}
	log.Printf("[Answer to Large Sample in Part 1] The sum is %d (should be 10092)", largeSampleSoln)

	inputSoln, err := solve("input.txt")
	if err != nil {
		log.Fatalf("failed to solve for input: %v", err)
	}
	log.Printf("[Answer to Part 1] The sum is %d", inputSoln)
}

func solve(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open grid file: %v", err)
	}
	defer file.Close()

	g := grid.New[byte]()
	moves := []byte{}

	gridDone := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			gridDone = true
			continue
		}
		if !gridDone {
			g = append(g, []byte(line))
			continue
		}
		moves = append(moves, []byte(line)...)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed to scan file contents: %v", err)
	}

	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			if g[y][x] == indicatorRobot {
				return solveAtCoord(g, grid.Coordinate{Y: y, X: x}, moves), nil
			}
		}
	}
	return 0, errors.New("grid did not contain robot")
}

func solveAtCoord(g grid.Grid[byte], robot grid.Coordinate, moves []byte) int {
	// printGrid(g)
	for i := 0; i < len(moves); i++ {
		robot = doMove(g, robot, moves[i])
		// printGrid(g)
	}

	sum := 0
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			if g[y][x] == indicatorBox {
				sum += (100*y + x)
			}
		}
	}
	return sum
}

func doMove(g grid.Grid[byte], robot grid.Coordinate, move byte) grid.Coordinate {
	switch move {
	case moveLeft:
		for x := robot.X - 1; x >= 1; x-- {
			if g[robot.Y][x] == indicatorWall {
				// no empty space, just return original
				// robot location without moving anything
				return robot
			}
			if g[robot.Y][x] == indicatorEmpty {
				// shift elements to the left
				for shift := x; shift < robot.X; shift++ {
					g[robot.Y][shift] = g[robot.Y][shift+1]
				}
				g[robot.Y][robot.X] = indicatorEmpty
				return grid.Coordinate{Y: robot.Y, X: robot.X - 1}
			}
		}
	case moveRight:
		for x := robot.X + 1; x < len(g[robot.Y])-1; x++ {
			if g[robot.Y][x] == indicatorWall {
				// no empty space, just return original
				// robot location without moving anything
				return robot
			}
			if g[robot.Y][x] == indicatorEmpty {
				// shift elements to the right
				for shift := x; shift > robot.X; shift-- {
					g[robot.Y][shift] = g[robot.Y][shift-1]
				}
				g[robot.Y][robot.X] = indicatorEmpty
				return grid.Coordinate{Y: robot.Y, X: robot.X + 1}
			}
		}
	case moveUp:
		for y := robot.Y - 1; y >= 1; y-- {
			if g[y][robot.X] == indicatorWall {
				// no empty space, just return original
				// robot location without moving anything
				return robot
			}
			if g[y][robot.X] == indicatorEmpty {
				// shift elements upward
				for shift := y; shift < robot.Y; shift++ {
					g[shift][robot.X] = g[shift+1][robot.X]
				}
				g[robot.Y][robot.X] = indicatorEmpty
				return grid.Coordinate{Y: robot.Y - 1, X: robot.X}
			}
		}
	case moveDown:
		for y := robot.Y + 1; y < len(g)-1; y++ {
			if g[y][robot.X] == indicatorWall {
				// no empty space, just return original
				// robot location without moving anything
				return robot
			}
			if g[y][robot.X] == indicatorEmpty {
				// shift elements downward
				for shift := y; shift > robot.Y; shift-- {
					g[shift][robot.X] = g[shift-1][robot.X]
				}
				g[robot.Y][robot.X] = indicatorEmpty
				return grid.Coordinate{Y: robot.Y + 1, X: robot.X}
			}
		}
	}

	// no empty space, just return original
	// robot location without moving anything
	return robot
}

// func printGrid(g grid.Grid[byte]) {
// 	s := ""
// 	for y := 0; y < len(g); y++ {
// 		s += fmt.Sprintf("%s\n", string(g[y]))
// 	}
// 	fmt.Println(strings.TrimPrefix(s, "\n"))
// }
