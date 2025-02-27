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
	hands []*hand
}

type hand struct {
	cards [5]byte
	bid   int

	bestHand int
}

var (
	handHighCard     = 0
	handOnePair      = 1
	handTwoPair      = 2
	handThreeOfAKind = 3
	handFullHouse    = 4
	handFourOfAKind  = 5
	handFiveOfAKind  = 6

	relativeHandStrength = map[byte]uint8{
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'T': 10,
		'J': 11,
		'Q': 12,
		'K': 13,
		'A': 14,
	}
)

func aBetterThanB(a, b *hand) bool {
	if a.bestHand < b.bestHand {
		return true
	}
	if a.bestHand > b.bestHand {
		return false
	}

	for i := 0; i < 5; i++ {
		rhsA := relativeHandStrength[a.cards[i]]
		rhsB := relativeHandStrength[b.cards[i]]
		if rhsA < rhsB {
			return true
		}
		if rhsA > rhsB {
			return false
		}
	}

	// they are exactly the same
	return true
}

func newHand(cards [5]byte, bid int, jRule bool /* TODO */) *hand {
	freq := make(map[byte]int)
	for _, card := range cards {
		freq[card]++
	}

	isFive := false
	isFour := false
	threes := 0
	twos := 0

	for _, fr := range freq {
		if fr > 4 {
			isFive = true
			break
		}
		if fr > 3 {
			isFour = true
			break
		}
		if fr > 2 {
			threes++
			continue
		}
		if fr > 1 {
			twos++
			continue
		}
	}

	// Find the best hand from possibilities
	bestHand := handHighCard
	if isFive {
		bestHand = handFiveOfAKind
	} else if isFour {
		bestHand = handFourOfAKind
	} else if threes > 0 && twos > 0 {
		bestHand = handFullHouse
	} else if threes > 0 {
		bestHand = handThreeOfAKind
	} else if twos > 1 {
		bestHand = handTwoPair
	} else if twos > 0 {
		bestHand = handOnePair
	} else {
		bestHand = handHighCard
	}

	return &hand{
		cards:    cards,
		bid:      bid,
		bestHand: bestHand,
	}
}

// Possibilities of hands with jokers
type possibility struct {
	fourKind  bool
	threeKind int
	pairs     int
}

// Helper function to calculate possibilities based on current distribution of jokers
func calculatePossibility(freq map[byte]int, jokers int) possibility {
	pairs := 0
	threes := 0
	fourKind := false

	for _, count := range freq {
		switch {
		case count == 4:
			fourKind = true
		case count == 3:
			threes++
		case count == 2:
			pairs++
		}
	}

	// Allocate remaining jokers to maximize hand value
	for jokers > 0 && !fourKind {
		if threes > 0 {
			fourKind = true
			threes--
			jokers--
		} else if pairs > 0 {
			threes++
			pairs--
			jokers--
		} else {
			break
		}
	}

	return possibility{
		fourKind:  fourKind,
		threeKind: threes,
		pairs:     pairs,
	}
}

func solvePart1(filename string) (int, error) {
	in, err := loadInput(filename, false)
	if err != nil {
		return 0, fmt.Errorf("failed to load input: %v", err)
	}

	sort.Slice(in.hands, func(a, b int) bool { return aBetterThanB(in.hands[a], in.hands[b]) })

	sum := 0
	for i, hand := range in.hands {
		sum += (hand.bid * (i + 1))
	}
	return sum, nil
}

func solvePart2(filename string) (int, error) {
	in, err := loadInput(filename, true)
	if err != nil {
		return 0, fmt.Errorf("failed to load input: %v", err)
	}

	sort.Slice(in.hands, func(a, b int) bool { return aBetterThanB(in.hands[a], in.hands[b]) })

	sum := 0
	for i, hand := range in.hands {
		sum += (hand.bid * (i + 1))
	}
	return sum, nil
}

func loadInput(filename string, jRule bool) (*input, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	input := &input{hands: []*hand{}}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			cardsStr, bidStr, ok := strings.Cut(line, " ")
			if !ok {
				return nil, fmt.Errorf("invalid line, no space: %s", line)
			}
			bid, err := strconv.Atoi(bidStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert bid to int: %v", err)
			}
			input.hands = append(input.hands, newHand([5]byte([]byte(cardsStr)), bid, jRule))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file contents: %v", err)
	}

	return input, nil
}

func main() {
	sampleSoln, err := solvePart1("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 1 for sample input: %v", err)
	}
	log.Printf("[Answer to Sample in Part 1] The total winnings are %d (should be 6440)", sampleSoln)

	soln, err := solvePart1("input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 1 for input: %v", err)
	}
	log.Printf("[Answer to Part 1] The total winnings are %d", soln)

	// TODO
	sampleSoln2, err := solvePart2("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 2 for sample input: %v", err)
	}
	log.Printf("[Answer to Sample in Part 2] The total winnings are %d (should be 5905)", sampleSoln2)

	soln2, err := solvePart2("input.txt")
	if err != nil {
		log.Fatalf("failed to solve part 2 for input: %v", err)
	}
	log.Printf("[Answer to Part 2] The total winnings are %d", soln2)
}
