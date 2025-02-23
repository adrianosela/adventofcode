package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/adrianosela/adventofcode/utils/slice"
)

type input struct {
	seeds  []int
	layers []layer
}

type layer struct {
	conds []cond
}

type cond struct {
	dst  int
	src  int
	size int
}

func solvePart1(filename string) (int, error) {
	in, err := loadInput(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to load input: %v", err)
	}

	lowest := int(math.MaxInt)
	for _, seed := range in.seeds {
		loc := location(seed, in.layers)
		if loc < lowest {
			lowest = loc
		}
	}
	return lowest, nil
}

// there are better ways to solve this (e.g. recursive DFS checking both
// overlapping and non-overlapping parts)... but this works (just takes
// ~5 mins on my 16GB laptop).
func solvePart2(filename string) (int, error) {
	in, err := loadInput(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to load input: %v", err)
	}

	lowest := int(math.MaxInt)

	// in part 2 seeds come in pairs where the first part
	// is the start and the second is the range
	for s := 0; s < len(in.seeds)-1; s += 2 {
		for x := 0; x < in.seeds[s+1]; x++ {
			loc := location(in.seeds[s]+x, in.layers)
			if loc < lowest {
				lowest = loc
			}
		}
	}
	return lowest, nil
}

func location(cur int, layers []layer) int {
	for l := 0; l < len(layers); l++ {
		for c := 0; c < len(layers[l].conds); c++ {
			if cur >= layers[l].conds[c].src && cur < layers[l].conds[c].src+layers[l].conds[c].size {
				diff := cur - layers[l].conds[c].src
				cur = layers[l].conds[c].dst + diff
				break
			}
		}
	}
	return cur
}

func loadInput(filename string) (*input, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	input := &input{}
	parsedSeeds := false
	activeLayer := (*layer)(nil)
	for scanner.Scan() {
		line := scanner.Text()

		if !parsedSeeds {
			_, after, ok := strings.Cut(line, ":")
			if !ok {
				return nil, fmt.Errorf("unexpected seed line (doesnt have ':'): %s", line)
			}
			seeds, err := slice.StringsToInts(strings.Fields(after))
			if err != nil {
				return nil, fmt.Errorf("failed to convert seed part to integers: %v", err)
			}
			input.seeds = seeds
			parsedSeeds = true
			continue
		}

		if len(line) == 0 {
			if activeLayer != nil {
				input.layers = append(input.layers, *activeLayer)
			}
			activeLayer = nil
			continue
		}

		if strings.Contains(line, "map") {
			activeLayer = &layer{}
			continue
		}

		ints, err := slice.StringsToInts(strings.Fields(line))
		if err != nil {
			return nil, fmt.Errorf("failed to convert cond line to integers: %v", err)
		}
		if len(ints) != 3 {
			return nil, fmt.Errorf("cond line has more than 3 parts: %s", line)
		}

		activeLayer.conds = append(activeLayer.conds, cond{
			dst:  ints[0],
			src:  ints[1],
			size: ints[2],
		})
	}
	if activeLayer != nil {
		input.layers = append(input.layers, *activeLayer)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file contents: %v", err)
	}

	return input, nil
}

func main() {
	sampleSoln, err := solvePart1("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 1 for sample input: %v", err)
	}
	log.Printf("[Answer to Small Sample in Part 1] The lowest location is %d (should be 35)", sampleSoln)

	soln, err := solvePart1("input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 1 for input: %v", err)
	}
	log.Printf("[Answer to Part 1] The lowest location is %d", soln)

	sampleSoln2, err := solvePart2("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 2 for sample input: %v", err)
	}
	log.Printf("[Answer to Small Sample in Part 2] The lowest location is %d (should be 46)", sampleSoln2)

	soln2, err := solvePart2("input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 2 for input: %v", err)
	}
	log.Printf("[Answer to Part 2] The lowest location is %d", soln2)
}
