package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	instrDo   = "do()"
	instrDont = "don't()"
)

type input struct {
	pairs [][]int
}

func main() {
	inputPath := "input.txt"

	input, err := parseInputForPart1(inputPath)
	if err != nil {
		log.Fatalf("failed to parse input file: %v", err)
	}

	sum := 0
	for i := 0; i < len(input.pairs); i++ {
		mul := 1
		for j := 0; j < len(input.pairs[i]); j++ {
			mul *= input.pairs[i][j]
		}
		sum += mul
	}

	log.Printf("[Answer to Part 1] The sum is: %d", sum)

	input, err = parseInputForPart2(inputPath)
	if err != nil {
		log.Fatalf("failed to parse input file: %v", err)
	}

	sum = 0
	for i := 0; i < len(input.pairs); i++ {
		mul := 1
		for j := 0; j < len(input.pairs[i]); j++ {
			mul *= input.pairs[i][j]
		}
		sum += mul
	}

	log.Printf("[Answer to Part 2] The sum is: %d", sum)
}

func parseInputForPart1(path string) (*input, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file at path \"%s\": %v", path, err)

	}
	defer file.Close()

	pairs := [][]int{}

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		line := scanner.Text()
		lines++

		parts := strings.Split(line, "mul(")

		// valid instructions should be parts of the form
		// `${INT},${INT})` with some trailing garbage
		for _, part := range parts {
			maybeFirstInt, rest, ok := strings.Cut(part, ",")
			if !ok {
				// invalid, does not have a comma
				continue
			}

			firstInt, err := strconv.Atoi(maybeFirstInt)
			if err != nil {
				// invalid, not an integer
				continue
			}

			maybeSecondInt, _, ok := strings.Cut(rest, ")")
			if !ok {
				// invalid, does not have a closing parenthesis
				continue
			}

			secondInt, err := strconv.Atoi(maybeSecondInt)
			if err != nil {
				// invalid, not an integer
				continue
			}

			pairs = append(pairs, []int{firstInt, secondInt})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file at path \"%s\": %v", path, err)
	}

	return &input{pairs: pairs}, nil
}

func parseInputForPart2(path string) (*input, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file at path \"%s\": %v", path, err)

	}
	defer file.Close()

	re := regexp.MustCompile(`mul\(\d+,\d+\)|do\(\)|don't\(\)`)
	enabled := true

	scanner := bufio.NewScanner(file)
	lines := 0
	pairs := [][]int{}
	for scanner.Scan() {
		line := scanner.Bytes()
		lines++

		matches := re.FindAll(line, -1)
		for _, match := range matches {
			matchStr := string(match)

			switch matchStr {
			case instrDo:
				enabled = true
			case instrDont:
				enabled = false
			default:
				if !enabled {
					continue
				}

				_, part, ok := strings.Cut(matchStr, "mul(")
				if !ok {
					// invalid, does not start with expected prefix
					continue
				}

				maybeFirstInt, rest, ok := strings.Cut(part, ",")
				if !ok {
					// invalid, does not have a comma
					continue
				}

				firstInt, err := strconv.Atoi(maybeFirstInt)
				if err != nil {
					// invalid, not an integer
					continue
				}

				maybeSecondInt, _, ok := strings.Cut(rest, ")")
				if !ok {
					// invalid, does not have a closing parenthesis
					continue
				}

				secondInt, err := strconv.Atoi(maybeSecondInt)
				if err != nil {
					// invalid, not an integer
					continue
				}

				pairs = append(pairs, []int{firstInt, secondInt})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file at path \"%s\": %v", path, err)
	}

	return &input{pairs: pairs}, nil
}
