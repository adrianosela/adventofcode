package main

import (
	"flag"
	"log"

	"github.com/adrianosela/adventofcode/utils/grid"
	"github.com/adrianosela/adventofcode/utils/set"
)

func main() {
	filename := flag.String("filename", "", "The path to the input file")
	debug := flag.Bool("debug", false, "Whether to print debug output or not")
	flag.Parse()

	g, err := grid.LoadByte(*filename)
	if err != nil {
		log.Fatalf("failed to load grid: %v", err)
	}

	log.Printf("[Answer to Part 1] The number of unique locations is: %d", bruteForceA(g, *debug))
	log.Printf("[Answer to Part 2] The number of unique locations is: %d", bruteForceB(g, *debug))

}

func bruteForceA(g grid.Grid[byte], debug bool) int {
	antinodes := set.New[string]()

	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			val := g[y][x]
			if val == '.' {
				continue
			}

			for yy := y; yy < len(g); yy++ {
				// note: we could start xx = x but only for row yy = y
				for xx := 0; xx < len(g[yy]); xx++ {
					// skip self
					if yy == y && xx == x {
						continue
					}
					// skip unequal character
					if val != g[yy][xx] {
						continue
					}

					// compute the distance between new point and original point
					yDiff, xDiff := yy-y, xx-x

					// mirrored a distance away from original character
					yRefA, xRefA := y-yDiff, x-xDiff
					if !(yRefA < 0 || yRefA >= len(g) || xRefA < 0 || xRefA >= len(g[yRefA])) {
						if debug {
							log.Printf("Found antinode of (y=%d,x=%d) and (y=%d,x=%d): (y=%d,x=%d)", y, x, yy, xx, yRefA, xRefA)
						}
						antinodes.Put((&grid.Coordinate{Y: yRefA, X: xRefA}).String())
					}

					// mirrored a distance away from second character
					yRefB, xRefB := yy+yDiff, xx+xDiff
					if !(yRefB < 0 || yRefB >= len(g) || xRefB < 0 || xRefB >= len(g[yRefB])) {
						if debug {
							log.Printf("Found antinode of (y=%d,x=%d) and (y=%d,x=%d): (y=%d,x=%d)", y, x, yy, xx, yRefB, xRefB)
						}
						antinodes.Put((&grid.Coordinate{Y: yRefB, X: xRefB}).String())
					}
				}
			}
		}
	}

	return antinodes.Size()
}

func bruteForceB(g grid.Grid[byte], debug bool) int {
	antinodes := set.New[string]()

	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			val := g[y][x]
			if val == '.' {
				continue
			}

			// add antenna itself as antinode
			if debug {
				log.Printf("Found antenna (also an antinode): (y=%d,x=%d)", y, x)
			}
			antinodes.Put((&grid.Coordinate{Y: y, X: x}).String())

			for yy := y; yy < len(g); yy++ {
				// note: we could start xx = x but only for row yy = y
				for xx := 0; xx < len(g[yy]); xx++ {
					// skip self
					if yy == y && xx == x {
						continue
					}
					// skip unequal character
					if val != g[yy][xx] {
						continue
					}

					// compute the distance between new point and original point
					yDiff, xDiff := yy-y, xx-x

					// mirrored a distance away from original character
					yRefA, xRefA := y-yDiff, x-xDiff
					for !(yRefA < 0 || yRefA >= len(g) || xRefA < 0 || xRefA >= len(g[yRefA])) {
						if debug {
							log.Printf("Found antinode of (y=%d,x=%d) and (y=%d,x=%d): (y=%d,x=%d)", y, x, yy, xx, yRefA, xRefA)
						}
						antinodes.Put((&grid.Coordinate{Y: yRefA, X: xRefA}).String())

						yRefA -= yDiff
						xRefA -= xDiff
					}

					// mirrored a distance away from second character
					yRefB, xRefB := yy+yDiff, xx+xDiff
					for !(yRefB < 0 || yRefB >= len(g) || xRefB < 0 || xRefB >= len(g[yRefB])) {
						if debug {
							log.Printf("Found antinode of (y=%d,x=%d) and (y=%d,x=%d): (y=%d,x=%d)", y, x, yy, xx, yRefB, xRefB)
						}
						antinodes.Put((&grid.Coordinate{Y: yRefB, X: xRefB}).String())

						yRefB += yDiff
						xRefB += xDiff
					}
				}
			}
		}
	}

	return antinodes.Size()
}
