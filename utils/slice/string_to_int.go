package slice

import (
	"fmt"
	"strconv"
)

func StringsToInts(strings []string) ([]int, error) {
	ints := make([]int, 0, len(strings))

	for i, elem := range strings {
		parsed, err := strconv.Atoi(elem)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value at index %d: %v", i, err)
		}
		ints = append(ints, parsed)
	}

	return ints, nil
}
