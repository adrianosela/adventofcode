package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

type offset struct {
	x int
	y int
}

const (
	charM = 'M'
	charA = 'A'
	charS = 'S'

	directionN  = "N"
	directionNE = "NE"
	directionE  = "E"
	directionSE = "SE"
	directionS  = "S"
	directionSW = "SW"
	directionW  = "W"
	directionNW = "NW"
)

var (
	directions = map[string]*offset{
		directionN:  {x: 0, y: -1},
		directionNE: {x: 1, y: -1},
		directionE:  {x: 1, y: 0},
		directionSE: {x: 1, y: 1},
		directionS:  {x: 0, y: 1},
		directionSW: {x: -1, y: 1},
		directionW:  {x: -1, y: 0},
		directionNW: {x: -1, y: -1},
	}
)

func main() {
	filename := flag.String("filename", "", "The path to the input file")
	debug := flag.Bool("debug", false, "Whether to print debug output or not")
	flag.Parse()

	// should be 18 for XMAS in sample-input.txt
	// should be 2642 for XMAS in input.txt
	find := "XMAS"
	log.Printf(
		"[Answer to Part 1] The number of %s occurrences is: %d",
		find,
		findString(*filename, find, *debug),
	)

	// should be 9 for XMAS in sample-input.txt
	// should be ? for XMAS in input.txt
	log.Printf(
		"[Answer to Part 2] The number of %s occurrences is: %d",
		"X-MAS",
		findCrossedMASes(*filename, *debug),
	)

}

func loadGrid(path string, debug bool) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	grid := [][]byte{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			grid = append(grid, []byte(line))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan input file: %v", err)
	}

	if debug && len(grid) > 0 {
		log.Printf("loaded %d by %d grid", len(grid), len(grid[0]))
	}

	return grid, nil
}

func findCrossedMASes(inputPath string, debug bool) int {
	grid, err := loadGrid(inputPath, debug)
	if err != nil {
		log.Fatalf("failed to load grid from input file: %v", err)
	}

	occurrences := 0
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == charA {
				cNW := getCharInDirection(grid, y, x, directionNW)
				cNE := getCharInDirection(grid, y, x, directionNE)
				cSW := getCharInDirection(grid, y, x, directionSW)
				cSE := getCharInDirection(grid, y, x, directionSE)
				// m northwest, s northeast, m southwest, s southeast
				if cNW == charM && cNE == charS && cSW == charM && cSE == charS {
					occurrences++
					continue
				}
				// m northwest, m northeast, s southwest, s southeast
				if cNW == charM && cNE == charM && cSW == charS && cSE == charS {
					occurrences++
					continue
				}
				// s northwest, m northeast, s southwest, m southeast
				if cNW == charS && cNE == charM && cSW == charS && cSE == charM {
					occurrences++
					continue
				}
				// s northwest, s northeast, m southwest, m southeast
				if cNW == charS && cNE == charS && cSW == charM && cSE == charM {
					occurrences++
					continue
				}
			}
		}
	}
	return occurrences
}

func findString(inputPath string, str string, debug bool) int {
	grid, err := loadGrid(inputPath, debug)
	if err != nil {
		log.Fatalf("failed to load grid from input file: %v", err)
	}

	if len(str) == 0 {
		return 0
	}

	occurrences := 0
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == str[0] {
				for direction := range directions {
					if match(grid, y, x, str, 1, direction, y, x, debug) {
						occurrences++
					}
				}
			}
		}
	}
	return occurrences
}

func getCharInDirection(
	grid [][]byte,
	posY int,
	posX int,
	direction string,
) byte {
	// calculate movement based on direction
	offset, ok := directions[direction]
	if !ok {
		panic(fmt.Sprintf("no offsets for direction %s", direction))
	}
	newX := posX + offset.x
	newY := posY + offset.y

	// check new positions are valid (stil in grid)
	if newY >= len(grid) || newY < 0 || newX >= len(grid[0]) || newX < 0 {
		return '.' // represents nil/invalid
	}

	return grid[newY][newX]
}

func match(
	grid [][]byte,
	posY int,
	posX int,
	find string,
	lookForPos int,
	direction string,
	startY int, // used for debugging
	startX int, // used for debugging
	debug bool, // used for debugging
) bool {
	if lookForPos == len(find) {
		if debug {
			log.Printf(
				"found %s from (y=%d,x=%d) to (y=%d,x=%d) (%s direction)",
				find,
				startY,
				startX,
				posY,
				posX,
				direction,
			)
		}
		return true
	}

	// calculate movement based on direction
	offset, ok := directions[direction]
	if !ok {
		panic(fmt.Sprintf("no offsets for direction %s", direction))
	}
	newX := posX + offset.x
	newY := posY + offset.y

	// check new positions are valid (stil in grid)
	if newY >= len(grid) || newY < 0 || newX >= len(grid[0]) || newX < 0 {
		return false
	}

	// if the current position does not match
	// the expected character, return early.
	if grid[newY][newX] != find[lookForPos] {
		return false
	}

	// if it does match, keep matching.
	return match(grid, newY, newX, find, lookForPos+1, direction, startY, startX, debug)
}
