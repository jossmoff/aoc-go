package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/jossmoff/aoc-go/utils"
	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

type Position struct {
	Row int
	Col int
}

func (p Position) LessThan(o Position) bool {
	return p.Row < o.Row || (p.Row == o.Row && p.Col < o.Col)
}

type Vector struct {
	Dy int
	Dx int
}

func parseInput(r io.Reader) ([][]rune, map[rune][]Position, error) {

	var grid [][]rune
	antennaMap := make(map[rune][]Position)
	scanner := bufio.NewScanner(r)

	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		gridRow := []rune(line)
		grid = append(grid, gridRow)
		for col, char := range gridRow {
			if char != '.' {
				antennaMap[char] = append(antennaMap[char], Position{Row: row, Col: col})
			}
		}
		row++
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return grid, antennaMap, nil
}

func isValidAntinode(gridX int, gridY int, antiNode Position) bool {
	return antiNode.Row >= 0 && antiNode.Row < gridX && antiNode.Col >= 0 && antiNode.Col < gridY
}

func computeAntiNodesForAntenna(grid [][]rune, antennaPositions []Position) []Position {
	antiNodes := make([]Position, 0)
	antennaPositionPairs := utils.CombinationsWithOrdering(antennaPositions, antennaPositions, func(p1 Position, p2 Position) bool {
		return p1.LessThan(p2)
	})
	gridY := len(grid)
	gridX := len(grid[0])
	for _, pair := range antennaPositionPairs {
		antenna1 := pair.V1
		antenna2 := pair.V2
		vector := Vector{Dy: antenna2.Row - antenna1.Row, Dx: antenna2.Col - antenna1.Col}
		for i := 0; isValidAntinode(gridX, gridY, Position{Row: antenna1.Row - i*vector.Dy, Col: antenna1.Col - i*vector.Dx}); i++ {
			antiNodes = append(antiNodes, Position{Row: antenna1.Row - i*vector.Dy, Col: antenna1.Col - i*vector.Dx})
		}
		for i := 0; isValidAntinode(len(grid[0]), len(grid), Position{Row: antenna2.Row + i*vector.Dy, Col: antenna2.Col + i*vector.Dx}); i++ {
			antiNodes = append(antiNodes, Position{Row: antenna2.Row + i*vector.Dy, Col: antenna2.Col + i*vector.Dx})
		}
	}

	return antiNodes
}

func run(part2 bool, input string) any {
	grid, antennaPositionMap, err := parseInput(strings.NewReader(input))
	if err != nil {
		panic(err)
	}
	antiNodes := hashset.New()
	for _, positions := range antennaPositionMap {
		for _, antiNode := range computeAntiNodesForAntenna(grid, positions) {
			antiNodes.Add(antiNode)
		}
	}
	return antiNodes.Size()
}
