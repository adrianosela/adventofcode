package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	schematicTypeUnknown = "UNKNOWN"
	schematicTypeLock    = "LOCK"
	schematicTypeKey     = "KEY"
)

type input struct {
	keys  [][]int
	locks [][]int
}

func main() {
	filename := flag.String("filename", "", "The path to the input file")
	schematicHeight := flag.Int("schematic-height", 6, "The height (in characters) of each schematic")
	debug := flag.Bool("debug", false, "Whether to print debug output or not")
	flag.Parse()

	input, err := loadInput(*filename, *schematicHeight)
	if err != nil {
		log.Fatalf("failed to load input: %v", err)
	}

	if *debug {
		log.Printf("Got keys:  %v", input.keys)
		log.Printf("Got locks: %v", input.locks)
	}

	log.Printf("[Answer to Part 1] The number of fit combinations is: %d", fitCombinations(input.keys, input.locks, *schematicHeight))
}

func fitCombinations(keys, locks [][]int, schematicHeight int) int {
	sum := 0
	for lockIndex := 0; lockIndex < len(locks); lockIndex++ {
		for keyIndex := 0; keyIndex < len(keys); keyIndex++ {
			fits := true
			for pinIndex := 0; pinIndex < len(locks[lockIndex]); pinIndex++ {
				if locks[lockIndex][pinIndex]+keys[keyIndex][pinIndex] > schematicHeight-1 {
					fits = false
					break
				}
			}
			if fits {
				sum++
			}
		}
	}
	return sum
}

func loadInput(path string, schematicHeight int) (*input, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	input := &input{
		keys:  make([][]int, 0),
		locks: make([][]int, 0),
	}

	reset := true
	currentSchematic := schematicTypeUnknown

	currentKey := 0
	currentLock := 0

	scanner := bufio.NewScanner(file)
	for lineNo := 0; scanner.Scan(); lineNo++ {
		line := scanner.Text()

		if len(line) == 0 {
			if currentSchematic == schematicTypeKey {
				currentKey++
			}
			if currentSchematic == schematicTypeLock {
				currentLock++
			}
			currentSchematic = schematicTypeUnknown
			reset = true
			continue
		}

		if reset {
			switch {
			case !strings.Contains(line, "."):
				currentSchematic = schematicTypeKey
				input.keys = append(input.keys, make([]int, len(line)))
			case !strings.Contains(line, "#"):
				currentSchematic = schematicTypeLock
				fill := make([]int, len(line))
				for i := 0; i < len(fill); i++ {
					fill[i] = schematicHeight - 1
				}
				input.locks = append(input.locks, fill)
			default:
				return nil, fmt.Errorf("invalid input at line %d: is beginning of schematic but not a key nor lock (%s)", lineNo, line)
			}
			reset = false
			continue
		}

		switch currentSchematic {
		case schematicTypeKey:
			for i, char := range line {
				if char == '#' {
					input.keys[currentKey][i]++
				}
			}
			continue
		case schematicTypeLock:
			for i, char := range line {
				if char == '.' {
					input.locks[currentLock][i]--
				}
			}
			continue
		default:
			return nil, fmt.Errorf("do not have a valid schematic type at line %d", lineNo)
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan input file: %v", err)
	}

	return input, nil
}
