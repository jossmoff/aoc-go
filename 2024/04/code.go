package main

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/jossmoff/aoc-go/utils"
	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

type Direction int

const (
	North Direction = iota
	South
	East
	West
	NorthEast
	NorthWest
	SouthEast
	SouthWest
)

type Delta struct {
	DeltaX int
	DeltaY int
}

var DirectionDeltaMap = map[Direction]Delta{
	North:     {DeltaX: 0, DeltaY: -1},
	South:     {DeltaX: 0, DeltaY: 1},
	East:      {DeltaX: 1, DeltaY: 0},
	West:      {DeltaX: -1, DeltaY: 0},
	NorthEast: {DeltaX: 1, DeltaY: -1},
	NorthWest: {DeltaX: -1, DeltaY: -1},
	SouthEast: {DeltaX: 1, DeltaY: 1},
	SouthWest: {DeltaX: -1, DeltaY: 1},
}

var DiagonalDirectionDeltaMap = map[Direction]Delta{
	NorthEast: {DeltaX: 1, DeltaY: -1},
	NorthWest: {DeltaX: -1, DeltaY: -1},
	SouthEast: {DeltaX: 1, DeltaY: 1},
	SouthWest: {DeltaX: -1, DeltaY: 1},
}

func readInput(r io.Reader) [][]byte {
	ls := bufio.NewScanner(r)

	var lines [][]byte
	for ls.Scan() {
		runSlice := []byte(ls.Text())
		lines = append(lines, runSlice)
	}
	return lines
}

func generateValidDirections(gridSizeX, gridSizeY, currentX, currentY, vectorLength int, directionMap map[Direction]Delta) []Direction {
	validDirections := []Direction{}

	for direction, delta := range directionMap {
		// Calculate the endpoint after moving vectorLength steps in this direction
		newX := currentX + delta.DeltaX*(vectorLength-1)
		newY := currentY + delta.DeltaY*(vectorLength-1)

		// Check if the endpoint is within bounds
		if newX >= 0 && newX < gridSizeY && newY >= 0 && newY < gridSizeX {
			validDirections = append(validDirections, direction)
		}
	}

	return validDirections
}

func countWordOccurencesInGrid(grid [][]byte, searchWord string) int {
	acc := 0
	gridSizeY := len(grid)
	gridSizeX := len(grid[0])
	for y, row := range grid {
		for x, _ := range row {
			// Little optimsiation to skip the search if the first letter doesn't match
			if grid[y][x] != searchWord[0] {
				continue
			}
			for _, direction := range generateValidDirections(gridSizeX, gridSizeY, x, y, len(searchWord), DirectionDeltaMap) {
				if searchWordPresent(grid, searchWord, x, y, direction) {
					acc++
				}
			}
		}
	}
	return acc
}

func countCrossesInGrid(grid [][]byte, searchWord string) int {
	acc := 0
	gridSizeY := len(grid)
	gridSizeX := len(grid[0])
	var foundMidPoints []int
	for y, row := range grid {
		for x, _ := range row {
			if grid[y][x] != searchWord[0] {
				continue
			}
			for _, direction := range generateValidDirections(gridSizeX, gridSizeY, x, y, len(searchWord), DiagonalDirectionDeltaMap) {
				if searchWordPresent(grid, searchWord, x, y, direction) {
					deltaX := DiagonalDirectionDeltaMap[direction].DeltaX
					deltaY := DiagonalDirectionDeltaMap[direction].DeltaY
					midPointNormalisedCoordinate := (y+deltaY)*gridSizeX + (x + deltaX)
					foundMidPoints = append(foundMidPoints, midPointNormalisedCoordinate)
				}
			}
		}
	}
	tally := utils.Tally(foundMidPoints)
	for _, v := range tally {
		if v >= 2 {
			acc++
		}
	}
	return acc
}

func countCrossesInGridFast(grid [][]byte, searchWord string) int {
	acc := 0
	gridSizeY := len(grid)
	gridSizeX := len(grid[0])
	var foundMidPoints []int
	var wg sync.WaitGroup

	for y, row := range grid {
		for x := range row {
			if grid[y][x] != searchWord[0] {
				continue
			}
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()
				for _, direction := range generateValidDirections(gridSizeX, gridSizeY, x, y, len(searchWord), DiagonalDirectionDeltaMap) {
					if searchWordPresent(grid, searchWord, x, y, direction) {
						deltaX := DiagonalDirectionDeltaMap[direction].DeltaX
						deltaY := DiagonalDirectionDeltaMap[direction].DeltaY
						midPointNormalisedCoordinate := (y+deltaY)*gridSizeX + (x + deltaX)
						foundMidPoints = append(foundMidPoints, midPointNormalisedCoordinate)
					}
				}
			}(x, y)
		}
	}

	wg.Wait()

	tally := utils.Tally(foundMidPoints)
	for _, v := range tally {
		if v >= 2 {
			acc++
		}
	}
	return acc
}

func searchWordPresent(grid [][]byte, searchWord string, x, y int, direction Direction) bool {
	directionDelta := DirectionDeltaMap[direction]

	for i := 0; i < len(searchWord); i++ {
		newX := x + directionDelta.DeltaX*i
		newY := y + directionDelta.DeltaY*i

		if grid[newY][newX] != searchWord[i] {
			return false
		}
	}

	return true

}

func run(part2 bool, input string) any {
	var grid [][]byte = readInput(strings.NewReader(input))
	if part2 {
		return countCrossesInGridFast(grid, "MAS")
	}
	return countWordOccurencesInGrid(grid, "XMAS")
}
