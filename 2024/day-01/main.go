package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type input struct {
	size  int
	listA []int
	listB []int
}

func main() {
	// For part one, we'll parse the inputs and sort the lists. Then
	// simply compute the difference between values on the same index
	// and add all those differences up.

	input, err := parseInput("input.txt", "   ")
	if err != nil {
		log.Fatalf("failed to parse input file: %v", err)
	}

	sort.Slice(input.listA, func(i, j int) bool { return input.listA[i] < input.listA[j] })
	sort.Slice(input.listB, func(i, j int) bool { return input.listB[i] < input.listB[j] })

	sum := 0
	for i := 0; i < input.size; i++ {
		diff := input.listA[i] - input.listB[i]
		if diff < 0 {
			diff *= -1
		}
		sum += diff
	}

	log.Printf("[Answer to Part 1] The sum is: %d", sum)

	// For part two, we'll keep using the sorted lists such that once we know
	// the index of the first occurrence of a given number, we can simply find
	// the amount of times that number appears consecutively and stop when we
	// see a different (higher) number.

	similarity := 0
	position := 0

	// NOTE: Can probably skip the map. Kust keep track of the last value checked.
	memo := make(map[int]int)
	for i := 0; i < input.size; i++ {
		lookFor := input.listA[i]

		if count, ok := memo[lookFor]; ok {
			similarity += (lookFor * count)
			continue
		}

		count, newPosition := countOccurrences(input.listB, position, lookFor)
		memo[lookFor] = count
		position = newPosition

		similarity += (lookFor * count)
	}

	log.Printf("[Answer to Part 2] The similarity score is: %d", similarity)
}

func parseInput(path string, separator string) (*input, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file at path \"%s\": %v", path, err)

	}
	defer file.Close()

	listA := []int{}
	listB := []int{}

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		line := scanner.Text()
		lines++

		parts := strings.Split(line, separator)
		if len(parts) != 2 {
			return nil, fmt.Errorf("unexpected number of parts at line %d (\"%s\")", lines, line)
		}

		locationIDA, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("first value not an integer at line %d (\"%s\")", lines, line)
		}
		listA = append(listA, locationIDA)

		locationIDB, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("second value not an integer at line %d (\"%s\")", lines, line)
		}
		listB = append(listB, locationIDB)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file at path \"%s\": %v", path, err)
	}

	return &input{size: lines, listA: listA, listB: listB}, nil
}

func countOccurrences(slice []int, startPos int, val int) (int, int) {
	occurrences := 0

	for i := startPos; i < len(slice); i++ {
		if slice[i] == val {
			occurrences++
			continue
		}

		if slice[i] > val {
			return occurrences, i
		}
	}

	return occurrences, len(slice)
}
