package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/adrianosela/adventofcode/utils/sliceconv"
)

const (
	directionUnknown    = "UNKNOWN"
	directionIncreasing = "INCREASING"
	directionDecreasing = "DECREASING"
)

func main() {
	inputPath := "input.txt"

	part1(inputPath)
	part2(inputPath)
}

func part1(inputPath string) {
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("failed to open input file at path \"%s\": %v", inputPath, err)
	}
	defer file.Close()

	safeReports := 0
	totalReports := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		totalReports++

		levels, err := sliceconv.StringsToInts(strings.Split(line, " "))
		if err != nil {
			log.Fatalf("failed to parse values at line %d \"%s\": %v", totalReports, line, err)
		}

		if isSafeReport(levels) {
			safeReports++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("failed to scan file at path \"%s\": %v", inputPath, err)
	}

	log.Printf("[Answer to Part 1] The number of safe reports is: %d/%d", safeReports, totalReports)
}

func part2(inputPath string) {
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("failed to open input file at path \"%s\": %v", inputPath, err)
	}
	defer file.Close()

	safeReports := 0
	totalReports := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		totalReports++

		levels, err := sliceconv.StringsToInts(strings.Split(line, " "))
		if err != nil {
			log.Fatalf("failed to parse values at line %d \"%s\": %v", totalReports, line, err)
		}

		if isSafeReportWithDampener(levels) {
			safeReports++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("failed to scan file at path \"%s\": %v", inputPath, err)
	}

	log.Printf("[Answer to Part 2] The number of safe reports is: %d/%d", safeReports, totalReports)
}

func isSafeReport(levels []int) bool {
	direction := directionUnknown
	for i := 1; i < len(levels); i++ {
		diff := levels[i] - levels[i-1]

		abs := diff
		if abs < 0 {
			abs *= -1
		}

		// valid only if difference is at least one and at most three.
		if abs < 1 || abs > 3 {
			return false
		}

		switch direction {
		case directionUnknown:
			// difference is positive (sequence is increasing)
			if diff > 0 {
				direction = directionIncreasing
				continue
			}
			// difference is negative (sequence is decreasing)
			direction = directionDecreasing
			continue
		case directionDecreasing:
			// difference is positive (sequence is increasing)
			if diff > 0 {
				return false
			}
			continue
		case directionIncreasing:
			// difference is negative (sequence is decreasing)
			if diff < 0 {
				return false
			}
			continue
		}
	}

	return true
}

func isSafeReportWithDampener(levels []int) bool {
	if isSafeReport(levels) {
		return true
	}
	for i := 0; i < len(levels); i++ {
		if isSafeReport(copyWithIndexRemoved(levels, i)) {
			return true
		}
	}
	return false
}

func copyWithIndexRemoved(original []int, index int) []int {
	clone := make([]int, 0, len(original)-1)
	clone = append(clone, original[:index]...)
	clone = append(clone, original[index+1:]...)
	return clone
}
