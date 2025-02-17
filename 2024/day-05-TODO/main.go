package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/adrianosela/adventofcode/utils/set"
	"github.com/adrianosela/adventofcode/utils/sliceconv"
)

func main() {
	filename := flag.String("filename", "sample-input.txt", "The path to the input file")
	debug := flag.Bool("debug", false, "Whether to print debug output or not")
	flag.Parse()

	fmt.Println(solve(*filename, *debug))
}

func solve(filename string, debug bool) int {
	rules, updates, err := loadInput(filename)
	if err != nil {
		log.Fatalf("failed to load inputs from file: %v", err)
	}

	if debug {
		log.Printf("Got Rules: %v", rules)
		log.Printf("Got Updates: %v", updates)
	}

	// TODO
	return -1
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
