package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adrianosela/adventofcode/utils/slice"
)

func solvePart1(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0
	for scanner.Scan() {
		ints, err := slice.StringsToInts(strings.Fields(scanner.Text()))
		if err != nil {
			return 0, fmt.Errorf("failed to convert line to integers: %v", err)
		}
		sum += getNext(ints)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed to scan file contents: %v", err)
	}

	return sum, nil
}

func getNext(ints []int) int {
	diffs := make([]int, len(ints)-1)

	refDiff := 0
	isAllSame := true

	for i := 0; i < len(ints)-1; i++ {
		diff := ints[i+1] - ints[i]
		diffs[i] = diff
		if i == 0 {
			refDiff = diff
			continue
		}
		if diff != refDiff {
			isAllSame = false
		}
	}

	if !isAllSame {
		diffs = append(diffs, getNext(diffs))
	}

	return ints[len(ints)-1] + diffs[len(diffs)-1]
}

func getPrevious(ints []int) int {
	diffs := make([]int, len(ints)-1)

	refDiff := 0
	isAllSame := true

	for i := len(ints) - 1; i >= 1; i-- {
		diff := ints[i] - ints[i-1]
		diffs[i-1] = diff
		if i == len(ints)-1 {
			refDiff = diff
			continue
		}
		if diff != refDiff {
			isAllSame = false
		}
	}

	if !isAllSame {
		diffs = append([]int{getPrevious(diffs)}, diffs...)
	}

	return ints[0] - diffs[0]
}

func solvePart2(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0
	for scanner.Scan() {
		ints, err := slice.StringsToInts(strings.Fields(scanner.Text()))
		if err != nil {
			return 0, fmt.Errorf("failed to convert line to integers: %v", err)
		}
		sum += getPrevious(ints)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed to scan file contents: %v", err)
	}

	return sum, nil
}

func main() {
	sampleSoln, err := solvePart1("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 1 for sample input: %v", err)
	}
	log.Printf("[Answer to Small Sample in Part 1] The sum is %d (should be 114)", sampleSoln)

	soln, err := solvePart1("input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 1 for input: %v", err)
	}
	log.Printf("[Answer to Part 1] The sum is %d", soln)

	sampleSoln2, err := solvePart2("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 2 for sample input: %v", err)
	}
	log.Printf("[Answer to Small Sample in Part 2] The number of cards is %d (should be 2)", sampleSoln2)

	soln2, err := solvePart2("input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 2 for input: %v", err)
	}
	log.Printf("[Answer to Part 2] The number of cards is %d", soln2)
}
