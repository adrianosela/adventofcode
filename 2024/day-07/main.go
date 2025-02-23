package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/adrianosela/adventofcode/utils/slice"
)

func main() {
	filename := flag.String("filename", "", "The path to the input file")
	flag.Parse()

	sum, err := part1(*filename)
	if err != nil {
		log.Fatalf("failed to solve part 1: %v", err)
	}
	log.Printf("[Answer to Part 1] The sum of valid totals is %d", sum)

	sum, err = part2(*filename)
	if err != nil {
		log.Fatalf("failed to solve part 2: %v", err)
	}
	log.Printf("[Answer to Part 2] The sum of valid totals is %d", sum)
}

func part1(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open input file: %v", err)
	}

	sum := 0
	scanner := bufio.NewScanner(file)
	for lineNo := 0; scanner.Scan(); lineNo++ {
		line := scanner.Text()

		totalStr, operandsStr, ok := strings.Cut(line, ": ")
		if !ok {
			return 0, fmt.Errorf("invalid input on line %d: does not contain colon", lineNo)
		}

		total, err := strconv.Atoi(totalStr)
		if err != nil {
			return 0, fmt.Errorf("invalid input on line %d: first part not an integer", lineNo)
		}

		operands, err := slice.StringsToInts(strings.Split(operandsStr, " "))
		if err != nil {
			return 0, fmt.Errorf("invalid input on line %d: %v", lineNo, err)
		}
		if len(operands) == 0 && (total != 0) {
			continue
		}

		if valid(total, operands[0], operands[1:]...) {
			sum += total
		}
	}

	return sum, nil
}

func valid(total int, val int, operands ...int) bool {
	if len(operands) == 0 {
		return total == val
	}
	return valid(total, val+operands[0], operands[1:]...) ||
		valid(total, val*operands[0], operands[1:]...)
}

func part2(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open input file: %v", err)
	}

	sum := 0
	scanner := bufio.NewScanner(file)
	for lineNo := 0; scanner.Scan(); lineNo++ {
		line := scanner.Text()

		totalStr, operandsStr, ok := strings.Cut(line, ": ")
		if !ok {
			return 0, fmt.Errorf("invalid input on line %d: does not contain colon", lineNo)
		}

		total, err := strconv.Atoi(totalStr)
		if err != nil {
			return 0, fmt.Errorf("invalid input on line %d: first part not an integer", lineNo)
		}

		operands, err := slice.StringsToInts(strings.Split(operandsStr, " "))
		if err != nil {
			return 0, fmt.Errorf("invalid input on line %d: %v", lineNo, err)
		}
		if len(operands) == 0 && (total != 0) {
			continue
		}

		if validWithConcat(total, operands[0], operands[1:]...) {
			sum += total
		}
	}

	return sum, nil
}

func validWithConcat(total int, val int, operands ...int) bool {
	if len(operands) == 0 {
		return total == val
	}
	return validWithConcat(total, val+operands[0], operands[1:]...) ||
		validWithConcat(total, val*operands[0], operands[1:]...) ||
		validWithConcat(total, concat(val, operands[0]), operands[1:]...)
}

func concat(a, b int) int {
	result, _ := strconv.Atoi(fmt.Sprintf("%s%s", strconv.Itoa(a), strconv.Itoa(b)))
	return result
}
