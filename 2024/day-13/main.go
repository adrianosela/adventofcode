package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	prefixButtonA = "Button A: "
	prefixButtonB = "Button B: "
	prefixPrize   = "Prize: "
)

type input struct {
	machines []machine
}

type machine struct {
	ax     int
	ay     int
	bx     int
	by     int
	prizex int
	prizey int
}

type result struct {
	tokens   int
	solvable bool
}

func (in *input) solvePart1(debug bool) int {
	sum := 0
	for i := 0; i < len(in.machines); i++ {

		if result := in.machines[i].solveRecursively(make(map[string]result), 0, 0, 0, 0); result.solvable {
			if debug {
				fmt.Printf("Machine %d is solvable with %d tokens\n", i+1, result.tokens)
			}
			sum += result.tokens
			continue
		}
		if debug {
			fmt.Printf("Machine %d is NOT solvable\n", i+1)
		}
	}
	return sum
}

func (in *input) solvePart2(debug bool) int {
	offset := 10000000000000
	sum := 0
	for i := 0; i < len(in.machines); i++ {
		machine := machine{
			ax:     in.machines[i].ax,
			ay:     in.machines[i].ay,
			bx:     in.machines[i].bx,
			by:     in.machines[i].by,
			prizex: in.machines[i].prizex + offset,
			prizey: in.machines[i].prizey + offset,
		}
		if result := machine.solveAnalytically(); result.solvable {
			if debug {
				fmt.Printf("Machine %d is solvable with %d tokens\n", i+1, result.tokens)
			}
			sum += result.tokens
			continue
		}
		if debug {
			fmt.Printf("Machine %d is NOT solvable\n", i+1)
		}
	}
	return sum
}

func (m *machine) solveRecursively(
	memo map[string]result,
	depthA int,
	depthB int,
	currentX int,
	currentY int,
) result {
	// it is memoized
	memoKey := fmt.Sprintf("%d-%d", currentX, currentY)
	if res, ok := memo[memoKey]; ok {
		return res
	}
	// it is solved
	if currentX == m.prizex && currentY == m.prizey {
		res := result{solvable: true}
		memo[memoKey] = res
		return res
	}
	// it is not solvable, we moved past
	if currentX > m.prizex || currentY > m.prizey {
		res := result{solvable: false}
		memo[memoKey] = res
		return res
	}
	// it is not solvable, we went too deep
	if depthA > 100 || depthB > 100 {
		res := result{solvable: false}
		memo[memoKey] = res
		return res
	}

	resultWithA := m.solveRecursively(memo, depthA+1, depthB, currentX+m.ax, currentY+m.ay)
	resultWithB := m.solveRecursively(memo, depthA, depthB+1, currentX+m.bx, currentY+m.by)

	// solvable either button press, return smallest number of moves
	if resultWithA.solvable && resultWithB.solvable {
		if resultWithA.tokens > resultWithB.tokens {
			res := result{tokens: 3 + resultWithA.tokens, solvable: true}
			memo[memoKey] = res
			return res
		} else {
			res := result{tokens: 1 + resultWithB.tokens, solvable: true}
			memo[memoKey] = res
			return res
		}
	}
	// solvable only with pressing A
	if resultWithA.solvable {
		res := result{tokens: 3 + resultWithA.tokens, solvable: true}
		memo[memoKey] = res
		return res
	}
	// solvable only with pressing B
	if resultWithB.solvable {
		res := result{tokens: 1 + resultWithB.tokens, solvable: true}
		memo[memoKey] = res
		return res
	}
	// not solvable
	res := result{solvable: false}
	memo[memoKey] = res
	return res
}

// Basically solves as a system of 2 equations and 2 unknowns:
// ⚫ Equation 1: (ax * P) + (bx * Q) = prizex
// ⚫ Equation 2: (ay * P) + (by * Q) = prizey
//
// First we isolate Q:
// ⚫ Q = ((pricey * ax) = (pricex * ay)) / ((by * ax) - (bx * ay))
// ⚫ Valid solutions will have integer values of Q
//
// Next we solve for P:
// ⚫ P = (prizex - (bx * Q)) / ax
//
// There would be a special case if ax is zero, but looking at inputs
// there are no buttons with No-Op values for the X coordinate.
func (m *machine) solveAnalytically() result {
	// calculate Q using the cross-multiplication of positions and movements
	qNumerator := m.prizey*m.ax - m.prizex*m.ay
	qDenominator := m.by*m.ax - m.bx*m.ay

	// check if Q can be calculated without any remainder (i.e., Q is an integer)
	if qDenominator != 0 && qNumerator%qDenominator == 0 {
		q := qNumerator / qDenominator
		// calculate P using the previously derived value of Q
		pNumerator := m.prizex - q*m.bx
		// ensure that P is an integer and calculate tokens if it is solvable
		if m.ax != 0 && pNumerator%m.ax == 0 {
			return result{tokens: 3*(pNumerator/m.ax) + q, solvable: true}
		}
	}

	// NOTE: there would be a special case where ax is zero, but looking at the
	// inputs, none of the machines have an A button with X zero.

	// not solvable
	return result{solvable: false}
}

func loadInput(filename string) (*input, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	machines := []machine{}
	currentMachine := machine{}

	for lineNo := 0; scanner.Scan(); lineNo++ {
		line := scanner.Text()
		if line == "" {
			machines = append(machines, currentMachine)
			currentMachine = machine{}
			continue
		}

		if strings.HasPrefix(line, prefixButtonA) {
			buttonAParts := strings.Split(strings.TrimPrefix(line, prefixButtonA), ", ")
			if len(buttonAParts) != 2 {
				return nil, fmt.Errorf("line %d does not have 2 parts (\"%s\")", lineNo, line)
			}

			axStr := strings.TrimPrefix(buttonAParts[0], "X+")
			ax, err := strconv.Atoi(axStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert X offset to integer: %v", err)
			}

			ayStr := strings.TrimPrefix(buttonAParts[1], "Y+")
			ay, err := strconv.Atoi(ayStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert Y offset to integer: %v", err)
			}

			currentMachine.ax = ax
			currentMachine.ay = ay
			continue
		}

		if strings.HasPrefix(line, prefixButtonB) {
			buttonBParts := strings.Split(strings.TrimPrefix(line, prefixButtonB), ", ")
			if len(buttonBParts) != 2 {
				return nil, fmt.Errorf("line %d does not have 2 parts (\"%s\")", lineNo, line)
			}

			bxStr := strings.TrimPrefix(buttonBParts[0], "X+")
			bx, err := strconv.Atoi(bxStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert X offset to integer: %v", err)
			}

			byStr := strings.TrimPrefix(buttonBParts[1], "Y+")
			by, err := strconv.Atoi(byStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert Y offset to integer: %v", err)
			}

			currentMachine.bx = bx
			currentMachine.by = by
			continue
		}

		if strings.HasPrefix(line, prefixPrize) {
			prizeParts := strings.Split(strings.TrimPrefix(line, prefixPrize), ", ")
			if len(prizeParts) != 2 {
				return nil, fmt.Errorf("line %d does not have 2 parts (\"%s\")", lineNo, line)
			}

			prizeXStr := strings.TrimPrefix(prizeParts[0], "X=")
			prizex, err := strconv.Atoi(prizeXStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert X offset to integer: %v", err)
			}

			prizeyStr := strings.TrimPrefix(prizeParts[1], "Y=")
			prizey, err := strconv.Atoi(prizeyStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert Y offset to integer: %v", err)
			}

			currentMachine.prizex = prizex
			currentMachine.prizey = prizey
			continue
		}
	}
	machines = append(machines, currentMachine)

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan over file: %v", err)
	}

	return &input{machines: machines}, nil
}

func main() {
	sampleIn, err := loadInput("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to load sample input: %v", err)
	}
	in, err := loadInput("input.txt")
	if err != nil {
		log.Fatalf("failed to load input: %v", err)
	}
	log.Printf("Answer to sample in part 1: %d", sampleIn.solvePart1(false))
	log.Printf("Answer to part 1: %d", in.solvePart1(false))
	log.Printf("Answer to sample in part 2: %d", sampleIn.solvePart2(false))
	log.Printf("Answer to part 2: %d", in.solvePart2(false))
}
