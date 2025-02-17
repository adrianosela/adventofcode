package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adrianosela/adventofcode-2024/utils/grid"
	"github.com/adrianosela/adventofcode-2024/utils/set"
	"github.com/adrianosela/adventofcode-2024/utils/sliceconv"
)

func main() {
	filename := flag.String("filename", "", "The path to the input file")
	debug := flag.Bool("debug", false, "Whether to print debug output or not")
	trailStart := flag.Int("trail-start", 0, "Value indicating start of the trail")
	trailEnd := flag.Int("trail-end", 9, "Value indicating end of the trail")
	flag.Parse()

	g, err := loadGrid(*filename)
	if err != nil {
		log.Fatalf("failed to load input grid: %v", err)
	}

	if *debug {
		fmt.Printf("Grid:\n-----------------\n%s-----------------\n", g.String())
	}

	fmt.Println(countUniqueTrails(g, *trailStart, *trailEnd, *debug)) // answer to part one
	fmt.Println(countPaths(g, *trailStart, *trailEnd, *debug))        // answer to part two
}

func loadGrid(path string) (grid.Grid[int], error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	grid := make(grid.Grid[int], 0)

	scanner := bufio.NewScanner(file)
	for lineNo := 0; scanner.Scan(); lineNo++ {
		line := scanner.Text()

		row, err := sliceconv.StringsToInts(strings.Split(line, ""))
		if err != nil {
			return nil, fmt.Errorf("failed to convert line %d to a slice of integers: %v", lineNo, err)
		}

		grid = append(grid, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan input file: %v", err)
	}

	return grid, nil
}

func countPaths(
	g grid.Grid[int],
	trailStart int,
	trailEnd int,
	debug bool,
) int {
	sum := 0
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			if g[y][x] == trailStart {
				sum += countAllPathsStartingAt(g, y, x, trailStart, trailEnd, y, x, []string{}, debug)
			}
		}
	}
	return sum
}

func countAllPathsStartingAt(
	g grid.Grid[int],
	y int,
	x int,
	currentValue int,
	targetValue int,

	// used for debugging output
	startY int,
	startX int,
	path []string,
	debug bool,
) int {
	// out of bounds
	if y < 0 || y >= len(g) || x < 0 || x >= len(g[y]) {
		return 0
	}
	// not expected value
	if g[y][x] != currentValue {
		return 0
	}
	// got to the end
	if currentValue == targetValue {
		if debug {
			log.Printf("Found path path starting at (y=%d,x=%d)\t%v", startY, startX, path)
		}
		return 1
	}
	return countAllPathsStartingAt(g, y+1, x, currentValue+1, targetValue, startY, startX, append(path, "⬇️"), debug) + // down
		countAllPathsStartingAt(g, y-1, x, currentValue+1, targetValue, startY, startX, append(path, "⬆️"), debug) + // up
		countAllPathsStartingAt(g, y, x+1, currentValue+1, targetValue, startY, startX, append(path, "➡️"), debug) + // right
		countAllPathsStartingAt(g, y, x-1, currentValue+1, targetValue, startY, startX, append(path, "⬅️"), debug) // left
}

func countUniqueTrails(
	g grid.Grid[int],
	trailStart int,
	trailEnd int,
	debug bool,
) int {
	dedup := set.New[string]()

	sum := 0
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			if g[y][x] == trailStart {
				sum += countTrailsStartingAt(g, dedup, y, x, trailStart, trailEnd, y, x, []string{}, debug)
			}
		}
	}
	return sum
}

func countTrailsStartingAt(
	g grid.Grid[int],
	dedup set.Set[string],
	y int,
	x int,
	currentValue int,
	targetValue int,

	// used for debugging output
	startY int,
	startX int,
	path []string,
	debug bool,
) int {
	// out of bounds
	if y < 0 || y >= len(g) || x < 0 || x >= len(g[y]) {
		return 0
	}
	// not expected value
	if g[y][x] != currentValue {
		return 0
	}
	// got to the end
	if currentValue == targetValue {
		pathKey := fmt.Sprintf("(%d,%d)-(%d,%d)", startY, startX, y, x)
		if dedup.Has(pathKey) {
			return 0
		}
		dedup.Put(pathKey)
		if debug {
			log.Printf("Found unique path starting at (y=%d,x=%d)\t%v", startY, startX, path)
		}
		return 1
	}
	return countTrailsStartingAt(g, dedup, y+1, x, currentValue+1, targetValue, startY, startX, append(path, "⬇️"), debug) + // down
		countTrailsStartingAt(g, dedup, y-1, x, currentValue+1, targetValue, startY, startX, append(path, "⬆️"), debug) + // up
		countTrailsStartingAt(g, dedup, y, x+1, currentValue+1, targetValue, startY, startX, append(path, "➡️"), debug) + // right
		countTrailsStartingAt(g, dedup, y, x-1, currentValue+1, targetValue, startY, startX, append(path, "⬅️"), debug) // left
}
