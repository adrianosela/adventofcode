package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/adrianosela/adventofcode/utils/set"
)

func main() {
	sampleInputMini, err := loadInput("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to load sample input data: %v", err)
	}
	log.Printf("[Answer to Sample in Part 1] The result is: %d (should be 37327623)", part1(sampleInputMini, 2000, false))

	input, err := loadInput("input.txt")
	if err != nil {
		log.Fatalf("failed to load input data: %v", err)
	}
	log.Printf("[Answer to Part 1] The result is: %d", part1(input, 2000, false))
	log.Printf("[Answer to Part 2] The result is: %d", part2(input, 4, 2000, false))
}

func loadInput(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	buyers := []int{}
	scanner := bufio.NewScanner(file)
	for lineNo := 0; scanner.Scan(); lineNo++ {
		line := scanner.Text()

		buyer, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line %d as an integer: %v", lineNo, err)
		}

		buyers = append(buyers, buyer)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan input file: %v", err)
	}

	return buyers, nil
}

func mix(a, b int) int {
	return a ^ b
}

func prune(a int) int {
	return a % 16777216
}

func do(a int) int {
	val := prune(mix(a*64, a))
	val = prune(mix(int(math.Floor(float64(val)/32.0)), val))
	val = prune(mix(val*2048, val))
	return val
}

func part1(buyers []int, cycles int, debug bool) int {
	sum := 0
	for b := 0; b < len(buyers); b++ {
		val := buyers[b]
		for c := 0; c < cycles; c++ {
			val = do(val)
		}
		sum += val

		if debug {
			log.Printf("After %d cycles, buyer %d with initial value %d became %d", cycles, b, buyers[b], val)
		}
	}
	return sum
}

func part2(buyers []int, sequenceLength int, cycles int, debug bool) int {
	diffSequenceToSum := make(map[string]int)
	maxSum := 0

	for b := 0; b < len(buyers); b++ {
		value := buyers[b]
		seenSequences := set.New[string]()
		diffSequence := []int{}

		for range cycles {
			nextValue := do(value)
			diff := (nextValue % 10) - (value % 10)
			diffSequence = append(diffSequence, diff)
			value = nextValue

			if len(diffSequence) < sequenceLength {
				// we don't have a full sequence of diffs
				// to test yet, move on to the next cycle.
				continue
			}

			sequenceKey := fmt.Sprintf("%v", diffSequence)

			// if this sequence had already been seen, the
			// monkey would have sold to this buyer, so we
			// only process sequences not already seen.
			if !seenSequences.Has(sequenceKey) {
				seenSequences.Put(sequenceKey)

				sequenceSumSoFar, ok := diffSequenceToSum[sequenceKey]
				if !ok {
					sequenceSumSoFar = 0
				}
				sequenceSumSoFar += (nextValue % 10)
				diffSequenceToSum[sequenceKey] = sequenceSumSoFar

				// update max sum in the loop to avoid having
				// to iterate over sequence sums map later
				if sequenceSumSoFar > maxSum {
					if debug {
						log.Printf("High score exceeded by sequence %v: %d", diffSequence, sequenceSumSoFar)
					}
					maxSum = sequenceSumSoFar
				}
			}

			diffSequence = diffSequence[1:]
		}
	}

	return maxSum
}
