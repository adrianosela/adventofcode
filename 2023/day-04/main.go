package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/adrianosela/adventofcode/utils/set"
	"github.com/adrianosela/adventofcode/utils/slice"
)

type scratchCard struct {
	winning set.Set[int]
	numbers []int
}

func newScratchCard(line string) (*scratchCard, error) {
	_, after, ok := strings.Cut(line, ":")
	if !ok {
		return nil, fmt.Errorf("no ':' in line \"%s\"", line)
	}

	winningNumbersPart, cardNumbersPart, ok := strings.Cut(after, "|")
	if !ok {
		return nil, fmt.Errorf("no '|' in line \"%s\"", line)
	}

	winningNumbersSet := set.New[int]()
	for _, winner := range strings.Fields(winningNumbersPart) {
		winnerInt, err := strconv.Atoi(winner)
		if err != nil {
			return nil, fmt.Errorf("failed to parse winning number: %v", err)
		}
		winningNumbersSet.Put(winnerInt)
	}

	cardNumbersStrings := strings.Fields(cardNumbersPart)
	cardNumbers := make([]int, 0, len(cardNumbersStrings))
	for _, num := range cardNumbersStrings {
		cardNum, err := strconv.Atoi(num)
		if err != nil {
			return nil, fmt.Errorf("failed to parse card number: %v", err)
		}
		cardNumbers = append(cardNumbers, cardNum)
	}

	return &scratchCard{winning: winningNumbersSet, numbers: cardNumbers}, nil
}

func (s *scratchCard) score() int {
	matches := s.matches()
	if matches == 0 {
		return 0
	}
	return 1 << (matches - 1)
}

func (s *scratchCard) matches() int {
	matches := 0
	for _, n := range s.numbers {
		if s.winning.Has(n) {
			matches++
		}
	}
	return matches
}

func solvePart1(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0
	for scanner.Scan() {
		sc, err := newScratchCard(scanner.Text())
		if err != nil {
			return 0, fmt.Errorf("failed to parse scratch card: %v", err)
		}
		sum += sc.score()
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed to scan file contents: %v", err)
	}

	return sum, nil
}

func solvePart2(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	cards := []*scratchCard{}
	for scanner.Scan() {
		sc, err := newScratchCard(scanner.Text())
		if err != nil {
			return 0, fmt.Errorf("failed to parse scratch card: %v", err)
		}
		cards = append(cards, sc)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed to scan file contents: %v", err)
	}

	// init all card counts to one (the original card)
	cardCounts := slice.Of(len(cards), 1)
	totalCards := 0
	for i, card := range cards {
		totalCards += cardCounts[i]

		for j := 1; j <= card.matches() && i+j < len(cards); j++ {
			cardCounts[i+j] += cardCounts[i]
		}
	}
	return totalCards, nil
}

func main() {
	sampleSoln, err := solvePart1("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 1 for sample input: %v", err)
	}
	log.Printf("[Answer to Small Sample in Part 1] The sum is %d (should be 13)", sampleSoln)

	soln, err := solvePart1("input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 1 for input: %v", err)
	}
	log.Printf("[Answer to Part 1] The sum is %d", soln)

	sampleSoln2, err := solvePart2("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 2 for sample input: %v", err)
	}
	log.Printf("[Answer to Small Sample in Part 2] The number of cards is %d (should be 30)", sampleSoln2)

	soln2, err := solvePart2("input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 2 for input: %v", err)
	}
	log.Printf("[Answer to Part 2] The number of cards is %d", soln2)
}
