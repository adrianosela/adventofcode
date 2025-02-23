package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	sampleSoln, err := solvePt1("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 1 for sample input: %v", err)
	}
	log.Printf("[Answer to Small Sample in Part 1] The sum is %d (should be 142)", sampleSoln)

	soln, err := solvePt1("input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 1 for input: %v", err)
	}
	log.Printf("[Answer to Part 1] The sum is %d", soln)

	sampleSoln2, err := solvePt2("sample-input-pt-2.txt")
	if err != nil {
		log.Fatalf("failed to solve part 2 for sample input: %v", err)
	}
	log.Printf("[Answer to Small Sample in Part 2] The sum is %d (should be 281)", sampleSoln2)

	soln2, err := solvePt2("input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 2 for input: %v", err)
	}
	log.Printf("[Answer to Part 2] The sum is %d", soln2)
}

func solvePt1(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	sum := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		firstSeen := false
		first := byte(0x48)
		last := byte(0x48)
		for _, elem := range scanner.Bytes() {
			if elem >= 48 && elem <= 57 {
				last = elem
				if !firstSeen {
					firstSeen = true
					first = elem
				}
			}
		}
		num, _ := strconv.Atoi(string([]byte{first, last}))
		sum += num
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed to scan file contents: %v", err)
	}

	return sum, nil
}

func solvePt2(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	stringValues := map[string]byte{
		"one":   '1',
		"two":   '2',
		"three": '3',
		"four":  '4',
		"five":  '5',
		"six":   '6',
		"seven": '7',
		"eight": '8',
		"nine":  '9',
	}

	sum := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()

		firstSeen := false
		first := byte(0x48)
		last := byte(0x48)

		for i := 0; i < len(str); i++ {
			for k, v := range stringValues {
				if strings.HasPrefix(str[i:], k) {
					last = v
					if !firstSeen {
						firstSeen = true
						first = v
					}
				}
			}

			if byte(str[i]) >= 48 && byte(str[i]) <= 57 {
				last = byte(str[i])
				if !firstSeen {
					firstSeen = true
					first = byte(str[i])
				}
			}
		}

		num, _ := strconv.Atoi(string([]byte{first, last}))
		sum += num
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed to scan file contents: %v", err)
	}

	return sum, nil
}
