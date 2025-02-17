package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/adrianosela/adventofcode/utils/grid"
	"github.com/adrianosela/adventofcode/utils/set"
)

func main() {
	filename := flag.String("filename", "", "The path to the input file")
	debug := flag.Bool("debug", false, "Whether to print debug output or not")
	trailStart := flag.Int("trail-start", 0, "Value indicating start of the trail")
	trailEnd := flag.Int("trail-end", 9, "Value indicating end of the trail")
	flag.Parse()

	g, err := grid.LoadInt(*filename, "")
	if err != nil {
		log.Fatalf("failed to load input grid: %v", err)
	}

	if *debug {
		fmt.Printf("Grid:\n-----------------\n%s-----------------\n", g.String())
	}

	fmt.Println(countUniqueTrails(g, *trailStart, *trailEnd, *debug)) // answer to part one
	fmt.Println(countPaths(g, *trailStart, *trailEnd, *debug))        // answer to part two
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
