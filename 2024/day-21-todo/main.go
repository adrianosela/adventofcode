package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/adrianosela/adventofcode/utils/grid"
)

const (
	buttonActivate = byte('A')
	buttonUp       = byte('^')
	buttonDown     = byte('v')
	buttonRight    = byte('>')
	buttonLeft     = byte('<')
)

var (
	// +---+---+---+
	// | 7 | 8 | 9 |
	// +---+---+---+
	// | 4 | 5 | 6 |
	// +---+---+---+
	// | 1 | 2 | 3 |
	// +---+---+---+
	// |N/A| 0 | A |
	// +---+---+---+
	numPadCoords = map[byte]*grid.Coordinate{
		'7': {X: 0, Y: 0},
		'8': {X: 1, Y: 0},
		'9': {X: 2, Y: 0},
		'4': {X: 0, Y: 1},
		'5': {X: 1, Y: 1},
		'6': {X: 2, Y: 1},
		'1': {X: 0, Y: 2},
		'2': {X: 1, Y: 2},
		'3': {X: 2, Y: 2},
		'0': {X: 1, Y: 3},
		'A': {X: 2, Y: 3},
	}
	numPadGapCoords = &grid.Coordinate{X: 0, Y: 3}

	// +---+---+---+
	// |N/A| ^ | A |
	// +---+---+---+
	// | < | v | > |
	// +---+---+---+
	dirPadCoords = map[byte]*grid.Coordinate{
		buttonUp:       {X: 1, Y: 0},
		buttonActivate: {X: 2, Y: 0},
		buttonLeft:     {X: 0, Y: 1},
		buttonDown:     {X: 1, Y: 1},
		buttonRight:    {X: 2, Y: 1},
	}
	dirPadGapCoords = &grid.Coordinate{X: 0, Y: 0}
)

func main() {
	sampleInput := []string{"029A", "980A", "179A", "456A", "379A"}
	sampleInputSolution, err := solvePart1(sampleInput, true)
	if err != nil {
		log.Fatalf("failed to solve part 1 for (sample) input %v: %v", sampleInputSolution, err)
	}
	log.Printf("[Answer to Sample in Part 1] The result is: %d (should be 126384)", sampleInputSolution)

	input := []string{"169A", "279A", "540A", "869A", "789A"}
	inputSolution, err := solvePart1(input, true)
	if err != nil {
		log.Fatalf("failed to solve part 1 for input %v: %v", input, err)
	}
	log.Printf("[Answer to Part 1] The result is: %d", inputSolution)

	part2Soln, err := solvePart2(input, true)
	if err != nil {
		log.Fatalf("failed to solve part 2 for input %v: %v", input, err)
	}
	log.Printf("[Answer to Part 2] The result is: %d", part2Soln)
}

func solvePart1(codes []string, debug bool) (int, error) {
	return solveWithBots(codes, 2, debug)
}

func solvePart2(codes []string, debug bool) (int, error) {
	return solveWithBots(codes, 25, debug)
}

// NOTE: in the problem statement, the fact that there
// were multiple robots inbetween hinted at needing the
// ability to handle a variable amount so had that since
// part1.
func solveWithBots(codes []string, bots int, debug bool) (int, error) {
	sum := 0
	for i := 0; i < len(codes); i++ {
		complexity, err := getComplexity(codes[i], bots, debug)
		if err != nil {
			return 0, fmt.Errorf("failed to get complexity for code at index %d (%s): %v", i, codes[i], err)
		}
		sum += complexity
	}
	return sum, nil
}

func getComplexity(code string, dirPadBots int, debug bool) (int, error) {
	seq := []byte{}
	currentNumPadPos := numPadCoords[buttonActivate]
	for _, c := range []byte(code) {
		if nextNumPadPos, ok := numPadCoords[c]; ok {
			seq = append(seq, numPadShortestSeq(currentNumPadPos, nextNumPadPos)...)
			currentNumPadPos = nextNumPadPos
		}
	}

	currentDirPadPos := dirPadCoords[buttonActivate]
	for b := 0; b < dirPadBots; b++ {
		botStart := time.Now()

		nextBotSeq := []byte{}
		for _, c := range seq {
			if nextDirPadPos, ok := dirPadCoords[c]; ok {
				nextBotSeq = append(nextBotSeq, dirPadShortestSeq(currentDirPadPos, nextDirPadPos)...)
				currentDirPadPos = nextDirPadPos
			}
		}
		seq = nextBotSeq

		if debug {
			log.Printf("done with bot %d after %s", b, time.Since(botStart))
		}
	}

	numericPartOfCode, err := strconv.Atoi(strings.TrimSuffix(code, "A"))
	if err != nil {
		return 0, fmt.Errorf("failed to get numeric part of code %s: %v", code, err)
	}

	return numericPartOfCode * len(seq), nil
}

func numPadShortestSeq(src, dst *grid.Coordinate) []byte {
	gapPossible := (src.Y == numPadGapCoords.Y && dst.X == numPadGapCoords.X) ||
		(src.X == numPadGapCoords.X && dst.Y == numPadGapCoords.Y)

	movesY := moveY(src, dst)
	movesX, dx := moveX(src, dst)

	var path []byte
	if dx < 0 != gapPossible {
		path = append(movesX, movesY...)
		return append(path, buttonActivate)
	}

	path = append(movesY, movesX...)
	return append(path, buttonActivate)
}

func dirPadShortestSeq(src, dst *grid.Coordinate) []byte {
	gapPossible := (src.Y == dirPadGapCoords.Y && dst.X == dirPadGapCoords.X) ||
		(src.X == dirPadGapCoords.X && dst.Y == dirPadGapCoords.Y)

	movesY := moveY(src, dst)
	movesX, dx := moveX(src, dst)

	var path []byte
	if dx < 0 != gapPossible {
		path = append(movesX, movesY...)
		return append(path, buttonActivate)
	}

	path = append(movesY, movesX...)
	return append(path, buttonActivate)
}

func moveY(start, end *grid.Coordinate) []byte {
	dy := end.Y - start.Y

	distance := dy
	char := buttonDown
	if distance < 0 {
		char = buttonUp
		distance *= -1
	}

	path := make([]byte, 0, distance)
	for range distance {
		path = append(path, char)
	}
	return path
}

func moveX(start, end *grid.Coordinate) ([]byte, int) {
	dx := end.X - start.X

	distance := dx
	char := buttonRight
	if distance < 0 {
		char = buttonLeft
		distance *= -1
	}

	path := make([]byte, 0, distance)
	for range distance {
		path = append(path, char)
	}
	return path, dx
}
