package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

type input struct {
	patterns [][]byte
	designs  [][]byte
}

func main() {
	debug := false

	sampleInput, err := loadInput("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to load sample input data: %v", err)
	}
	log.Printf("[Answer to Sample in Part 1] The result is: %d, should be 6", part1(sampleInput, debug))
	log.Printf("[Answer to Sample in Part 2] The result is: %d, should be 16", part2(sampleInput, debug))

	input, err := loadInput("input.txt")
	if err != nil {
		log.Fatalf("failed to load input data: %v", err)
	}
	log.Printf("[Answer to Part 1] The result is: %d", part1(input, debug))
	log.Printf("[Answer to Part 2] The result is: %d", part2(input, debug))
}

func loadInput(filename string) (*input, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	patterns := [][]byte{}
	designs := [][]byte{}
	patternsDone := false
	scanner := bufio.NewScanner(file)
	for lineNo := 0; scanner.Scan(); lineNo++ {
		line := scanner.Text()

		if len(line) == 0 {
			patternsDone = true
			continue
		}

		if !patternsDone {
			patterns = append(patterns, bytes.Split([]byte(line), []byte{',', ' '})...)
			continue
		}

		designs = append(designs, []byte(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan input file: %v", err)
	}

	return &input{
		patterns: patterns,
		designs:  designs,
	}, nil
}

func part1(in *input, debug bool) int {
	memo := make(map[string]bool)
	possible := 0
	for d := 0; d < len(in.designs); d++ {
		if isPossible(memo, in.designs[d], in.patterns) {
			possible++
		}
		if debug {
			log.Printf("Processed %d/%d designs", d+1, len(in.designs))
		}
	}
	return possible
}

func isPossible(memo map[string]bool, design []byte, patterns [][]byte) bool {
	if len(design) == 0 {
		return true
	}

	key := string(design)
	if possible, ok := memo[key]; ok {
		return possible
	}

	for p := 0; p < len(patterns); p++ {
		if bytes.HasPrefix(design, patterns[p]) {
			if isPossible(memo, bytes.TrimPrefix(design, patterns[p]), patterns) {
				memo[key] = true
				return true
			}
			memo[key] = false
		}
	}

	return false
}

func part2(in *input, debug bool) int {
	memo := make(map[string]int)
	possible := 0
	for d := 0; d < len(in.designs); d++ {
		possible += countPossible(memo, in.designs[d], in.patterns)
		if debug {
			log.Printf("Processed %d/%d designs", d+1, len(in.designs))
		}
	}
	return possible
}

func countPossible(memo map[string]int, design []byte, patterns [][]byte) int {
	if len(design) == 0 {
		return 1
	}

	if possible, ok := memo[string(design)]; ok {
		return possible
	}

	possible := 0
	for p := 0; p < len(patterns); p++ {
		if bytes.HasPrefix(design, patterns[p]) {
			trimmed := bytes.TrimPrefix(design, patterns[p])
			trimmedPossible := countPossible(memo, trimmed, patterns)
			memo[string(trimmed)] = trimmedPossible
			possible += trimmedPossible
		}
	}
	return possible
}
