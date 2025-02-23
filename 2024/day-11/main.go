package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/adrianosela/adventofcode/utils/slice"
)

func main() {
	input := flag.String("input", "814 1183689 0 1 766231 4091 93836 46", "The raw input")
	logState := flag.Bool("log-state", false, "Whether to log the stone's state after each blink")
	logDurations := flag.Bool("log-durations", false, "Whether to log the blinking computation's durations")
	flag.Parse()

	ints, err := slice.StringsToInts(strings.Split(*input, " "))
	if err != nil {
		log.Fatalf("failed to convert string to integer slice: %v", err)
	}

	log.Printf("[Answer to Part 1] The number of stones is %d", solve(ints, 25, *logState, *logDurations))
	log.Printf("[Answer to Part 2] The number of stones is %d", solveRecursive(ints, 75, *logState, *logDurations))
}

func solve(input []int, blinks int, logState bool, logDurations bool) int {
	slice := input
	if logState {
		log.Printf("After 0 blinks: %v", slice)
	}
	for blink := 1; blink <= blinks; blink++ {
		newSlice := []int{}
		start := time.Now()
		for i := 0; i < len(slice); i++ {
			val := slice[i]

			// rule 1: if zero, new rock is 1
			if val == 0 {
				newSlice = append(newSlice, 1)
				continue
			}

			// rule 2: if even number of digits, split in half
			str := strconv.Itoa(val)
			if len(str)%2 == 0 {
				firstHalf, _ := strconv.Atoi(str[:len(str)/2])
				secondHalf, _ := strconv.Atoi(str[len(str)/2:])
				newSlice = append(newSlice, firstHalf, secondHalf)

				continue
			}

			// rule 3: anything else, multiply by 2024
			newSlice = append(newSlice, val*2024)
		}
		if logState {
			log.Printf("After %d blinks: %v", blink, newSlice)
		}
		if logDurations {
			log.Printf("Blink %d took %s to compute", blink, time.Since(start))
		}
		slice = newSlice
	}
	return len(slice)
}

// part 2 takes waaaaay too long to solve iteratively, so we'll basically split
// every single stone 75 times individually, and count the leaf nodes of the tree.
func solveRecursive(input []int, blinks int, logState bool, logDurations bool) int {
	sum := 0
	for i, stone := range input {
		start := time.Now()

		stones := branch(make(map[string]int), stone, 0, blinks)
		sum += stones

		if logState {
			log.Printf("Processed stones after blinks for stone %d/%d (%d stones after 75 blinks of %d)", i+1, len(input), stones, stone)
		}
		if logDurations {
			log.Printf("Processed stones after blinks for stone %d/%d, took %s", i+1, len(input), time.Since(start))
		}
	}
	return sum
}

func branch(memo map[string]int, stone int, blink int, blinks int) int {
	if blink == blinks {
		return 1
	}

	key := fmt.Sprintf("%d-%d", stone, blink)
	if visited, ok := memo[key]; ok {
		return visited
	}

	if stone == 0 {
		result := branch(memo, 1, blink+1, blinks)
		memo[key] = result
		return result
	}

	str := strconv.Itoa(stone)
	if len(str)%2 == 0 {
		firstHalf, _ := strconv.Atoi(str[:len(str)/2])
		secondHalf, _ := strconv.Atoi(str[len(str)/2:])
		result := branch(memo, firstHalf, blink+1, blinks) + branch(memo, secondHalf, blink+1, blinks)
		memo[key] = result
		return result
	}

	result := branch(memo, stone*2024, blink+1, blinks)
	memo[key] = result
	return result
}
