package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/adrianosela/adventofcode/utils/grid"
)

type robot struct {
	position grid.Coordinate
	velocity grid.Coordinate
}

func main() {
	sampleRobots, err := loadInput("sample-input-12.txt")
	if err != nil {
		log.Fatalf("failed to load robots data: %v", err)
	}
	log.Printf("[Answer to Sample] The result is: %d", part1(sampleRobots, grid.Coordinate{X: 11, Y: 7}, 100))

	robots, err := loadInput("input.txt")
	if err != nil {
		log.Fatalf("failed to load robots data: %v", err)
	}
	log.Printf("[Answer to Part 1] The result is: %d", part1(robots, grid.Coordinate{X: 101, Y: 103}, 100))

	log.Printf("[Answer to Part 2] The number of seconds is: %d", part2(robots, grid.Coordinate{X: 101, Y: 103}, false))
}

func loadInput(filename string) ([]robot, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	robots := []robot{}

	scanner := bufio.NewScanner(file)
	for lineNo := 0; scanner.Scan(); lineNo++ {
		// space separated position and velocity
		positionPart, velocityPart, ok := strings.Cut(scanner.Text(), " ")
		if !ok {
			return nil, fmt.Errorf("invalid input on line %d (no space)", lineNo)
		}

		_, positionStr, ok := strings.Cut(positionPart, "p=")
		if !ok {
			return nil, fmt.Errorf("invalid input on line %d (position part does not start with p=)", lineNo)
		}
		position, err := parseCoord(positionStr)
		if err != nil {
			return nil, fmt.Errorf("invalid input on line %d (invalid position coordinate): %v", lineNo, err)
		}

		_, velocityStr, ok := strings.Cut(velocityPart, "v=")
		if !ok {
			return nil, fmt.Errorf("invalid input on line %d (velocity part does not start with v=)", lineNo)
		}
		velocity, err := parseCoord(velocityStr)
		if err != nil {
			return nil, fmt.Errorf("invalid input on line %d (invalid velocity coordinate): %v", lineNo, err)
		}

		robots = append(robots, robot{
			position: position,
			velocity: velocity,
		})
	}

	return robots, nil
}

// parses comma separated coordinates of the form "X,Y" e.g. "11,7"
func parseCoord(coords string) (grid.Coordinate, error) {
	xStr, yStr, ok := strings.Cut(coords, ",")
	if !ok {
		return grid.Coordinate{}, fmt.Errorf("coordinates not comma separated: \"%s\"", coords)
	}
	x, err := strconv.Atoi(xStr)
	if err != nil {
		return grid.Coordinate{}, fmt.Errorf("x coordinate not an integer: \"%s\"", xStr)
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return grid.Coordinate{}, fmt.Errorf("y coordinate not an integer: \"%s\"", xStr)
	}
	return grid.Coordinate{X: x, Y: y}, nil
}

func part1(robots []robot, gridDims grid.Coordinate, seconds int) int {
	quadNW := 0
	quadNE := 0
	quadSW := 0
	quadSE := 0

	for _, robot := range robots {
		afterX := (robot.position.X + robot.velocity.X*seconds) % gridDims.X
		afterY := (robot.position.Y + robot.velocity.Y*seconds) % gridDims.Y

		// when negative, we simply add the grid dimensions to correct
		if afterX < 0 {
			afterX += gridDims.X
		}
		if afterY < 0 {
			afterY += gridDims.Y
		}

		// is in northwest quadrant
		if afterX < gridDims.X/2 && afterY < gridDims.Y/2 {
			quadNW++
			continue
		}
		// is in northeast quadrant
		if afterX > gridDims.X/2 && afterY < gridDims.Y/2 {
			quadNE++
			continue
		}
		// is in southwest quadrant
		if afterX < gridDims.X/2 && afterY > gridDims.Y/2 {
			quadSW++
			continue
		}
		// is in southeast quadrant
		if afterX > gridDims.X/2 && afterY > gridDims.Y/2 {
			quadSE++
			continue
		}
	}

	return quadNW * quadNE * quadSW * quadSE
}

func part2(robots []robot, gridDims grid.Coordinate, debug bool) int {
	return part2TryLowestSafetyFactor(robots, gridDims, debug)
}

// // this did not return after 10 minutes of running...
// func part2TrySymmetricQuadrants(robots []robot, gridDims grid.Coordinate) int {
// 	for seconds := 0; true; seconds++ {

// 		// move all robots
// 		for _, robot := range robots {
// 			robot.position.X = (robot.position.X + robot.velocity.X) % gridDims.X
// 			robot.position.Y = (robot.position.Y + robot.velocity.Y) % gridDims.Y

// 			// when negative, we simply add the grid dimensions to correct
// 			if robot.position.X < 0 {
// 				robot.position.X += gridDims.X
// 			}
// 			if robot.position.Y < 0 {
// 				robot.position.Y += gridDims.Y
// 			}
// 		}

// 		// compute quadrants for all robots
// 		quadNW := 0
// 		quadNE := 0
// 		quadSW := 0
// 		quadSE := 0
// 		for _, robot := range robots {
// 			// is in northwest quadrant
// 			if robot.position.X < gridDims.X/2 && robot.position.Y < gridDims.Y/2 {
// 				quadNW++
// 				continue
// 			}
// 			// is in northeast quadrant
// 			if robot.position.X > gridDims.X/2 && robot.position.Y < gridDims.Y/2 {
// 				quadNE++
// 				continue
// 			}
// 			// is in southwest quadrant
// 			if robot.position.X < gridDims.X/2 && robot.position.Y > gridDims.Y/2 {
// 				quadSW++
// 				continue
// 			}
// 			// is in southeast quadrant
// 			if robot.position.X > gridDims.X/2 && robot.position.Y > gridDims.Y/2 {
// 				quadSE++
// 				continue
// 			}
// 		}

// 		// symmetrical across y axis?
// 		if quadNW == quadNE && quadSE == quadSW {
// 			return seconds
// 		}
// 	}

// 	return -1
// }

// this did not return after 10 minutes of running...
func part2TryLowestSafetyFactor(robots []robot, gridDims grid.Coordinate, debug bool) int {

	lowest := atomic.Uint32{}
	lowest.Store(math.MaxUint32)

	go func() {
		for seconds := 0; true; seconds++ {

			quadNW := 0
			quadNE := 0
			quadSW := 0
			quadSE := 0

			for _, robot := range robots {
				afterX := (robot.position.X + robot.velocity.X*seconds) % gridDims.X
				afterY := (robot.position.Y + robot.velocity.Y*seconds) % gridDims.Y

				// when negative, we simply add the grid dimensions to correct
				if afterX < 0 {
					afterX += gridDims.X
				}
				if afterY < 0 {
					afterY += gridDims.Y
				}

				// is in northwest quadrant
				if afterX < gridDims.X/2 && afterY < gridDims.Y/2 {
					quadNW++
					continue
				}
				// is in northeast quadrant
				if afterX > gridDims.X/2 && afterY < gridDims.Y/2 {
					quadNE++
					continue
				}
				// is in southwest quadrant
				if afterX < gridDims.X/2 && afterY > gridDims.Y/2 {
					quadSW++
					continue
				}
				// is in southeast quadrant
				if afterX > gridDims.X/2 && afterY > gridDims.Y/2 {
					quadSE++
					continue
				}
			}

			score := quadNW * quadNE * quadSW * quadSE
			if score < int(lowest.Load()) {
				lowest.Store(uint32(score))
				if debug {
					log.Printf("New lowest score after %d seconds: %d", seconds, score)
				}
			}
		}
	}()

	<-time.After(time.Second * 5)
	return int(lowest.Load())
}
