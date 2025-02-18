package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/adrianosela/adventofcode/utils/set"
	"github.com/adrianosela/adventofcode/utils/sliceconv"
)

func main() {
	log.Printf("[Answer to Sample in Part 1] The result is: %d (should be 143)", part1("sample-input.txt"))
	log.Printf("[Answer to Part 1] The result is: %d", part1("input.txt"))

	log.Printf("[Answer to Sample in Part 2] The result is: %d (should be 123)", part2("sample-input.txt"))
	log.Printf("[Answer to Part 2] The result is: %d", part2("input.txt"))
}

func part2(filename string) int {
	rules, updates, err := loadInput(filename)
	if err != nil {
		log.Fatalf("failed to load inputs from file: %v", err)
	}
	sum := 0
	for u := 0; u < len(updates); u++ {
		if !isCorrectOrder(updates[u], rules) {
			correctOrderInPlace(updates[u], rules)
			sum += updates[u][len(updates[u])/2]
		}
	}
	return sum
}

func part1(filename string) int {
	rules, updates, err := loadInput(filename)
	if err != nil {
		log.Fatalf("failed to load inputs from file: %v", err)
	}
	sum := 0
	for u := 0; u < len(updates); u++ {
		if isCorrectOrder(updates[u], rules) {
			sum += updates[u][len(updates[u])/2]
		}
	}
	return sum
}

func isCorrectOrder(update []int, rules map[int]set.Set[int]) bool {
	for i := 0; i < len(update); i++ {
		for j := i + 1; j < len(update); j++ {
			requirements, hasRequirements := rules[update[i]]
			if hasRequirements {
				if requirements.Has(update[j]) {
					return false
				}
			}
		}
	}
	return true
}

func correctOrderInPlace(update []int, rules map[int]set.Set[int]) {
	for i := 0; i < len(update)-1; i++ {
		for j := i + 1; j < len(update); j++ {
			requirements, hasRequirements := rules[update[i]]
			if hasRequirements && requirements.Has(update[j]) {
				update[i], update[j] = update[j], update[i]
				i = -1 // to ensure order
				break
			}
		}
	}
}

func loadInput(filename string) (map[int]set.Set[int], [][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open input file: %v", err)
	}

	scanner := bufio.NewScanner(file)

	// start by loading rules until we reach
	// an empty line, then we start loading updates.
	loadingRules := true

	rules := make(map[int]set.Set[int])
	updates := make([][]int, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if loadingRules {
				loadingRules = false
			}
			continue
		}

		if loadingRules {
			beforeStr, afterStr, ok := strings.Cut(line, "|")
			if !ok {
				return nil, nil, fmt.Errorf("invalid input while loading rules (does not contain delimeter '|'): \"%s\"", line)
			}

			before, err := strconv.Atoi(beforeStr)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid input while loading rules (first part not an integer): \"%s\"", line)
			}

			after, err := strconv.Atoi(afterStr)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid input while loading rules (second part not an integer): \"%s\"", line)
			}

			if _, ok := rules[after]; ok {
				rules[after].Put(before)
			} else {
				rules[after] = set.New[int](before)
			}

			continue
		}

		ints, err := sliceconv.StringsToInts(strings.Split(line, ","))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert update line to a slice of integers: %v", err)
		}
		updates = append(updates, ints)
	}

	return rules, updates, nil
}
