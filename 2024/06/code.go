package main

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/jpillora/puzzler/harness/aoc"
	"github.com/pkg/profile"
)

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type Delta struct {
	dx, dy int
}

var DirectionDeltaMap = map[Direction]Delta{
	North: Delta{0, -1},
	East:  Delta{1, 0},
	South: Delta{0, 1},
	West:  Delta{-1, 0},
}

type Position struct {
	x, y int
}

func (p Position) HashCode() uint64 {
	return uint64(p.y)<<16 | uint64(p.x)
}

type SimulationStep struct {
	position  Position
	direction Direction
}

func (s SimulationStep) HashCode() uint64 {
	return uint64(s.position.x)<<32 | uint64(s.position.y)<<16 | uint64(s.direction)
}

func (d Direction) TurnRight() Direction {
	return (d + 1) % 4
}

func main() {
	aoc.Harness(run)
}

func parseGrid(r io.Reader) ([][]rune, Position, error) {
	scanner := bufio.NewScanner(r)
	var grid [][]rune
	var coord Position
	found := false

	for y := 0; scanner.Scan(); y++ {
		line := scanner.Text()
		row := []rune(line)
		grid = append(grid, row)
		if !found {
			for x, char := range row {
				if char == '^' {
					coord = Position{x, y}
					found = true
					break
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, Position{-1, -1}, err
	}

	return grid, coord, nil
}

func simulateMovementUntilEscape(grid [][]rune, startingPosition Position, direction Direction, extraBlockPosition Position) (map[uint64]Position, bool) {
	dx := DirectionDeltaMap[direction].dx
	dy := DirectionDeltaMap[direction].dy

	y := startingPosition.y
	x := startingPosition.x
	pathPositions := make(map[uint64]Position)
	simulationSteps := make(map[uint64]struct{})

	for {
		position := Position{x, y}
		pathPositions[position.HashCode()] = position
		simulationStep := SimulationStep{position, direction}
		if _, ok := simulationSteps[simulationStep.HashCode()]; ok {
			return pathPositions, false
		}

		simulationSteps[simulationStep.HashCode()] = struct{}{}
		nextX := x + dx
		nextY := y + dy

		if nextY < 0 || nextY >= len(grid) || nextX < 0 || nextX >= len(grid[nextY]) {
			return pathPositions, true
		}
		if grid[nextY][nextX] == '#' || (extraBlockPosition.x == nextX && extraBlockPosition.y == nextY) {
			direction = direction.TurnRight()
			dx = DirectionDeltaMap[direction].dx
			dy = DirectionDeltaMap[direction].dy
			nextX = x
			nextY = y
		}

		x = nextX
		y = nextY
	}
}

func run(part2 bool, input string) any {
	defer profile.Start(profile.ProfilePath(".")).Stop()
	if part2 {
		grid, startingPosition, err := parseGrid(strings.NewReader(input))
		if err != nil {
			panic(err)
		}
		pathCoords, _ := simulateMovementUntilEscape(grid, startingPosition, North, Position{-1, -1})
		sneakyGuardPathCount := 0
		results := make(chan int)
		var wg sync.WaitGroup
		for _, obstaclePosition := range pathCoords {
			wg.Add(1)
			go func(obstaclePosition Position) {
				defer wg.Done()
				if err != nil {
					panic(err)
				}
				_, escapes := simulateMovementUntilEscape(grid, startingPosition, North, obstaclePosition)
				if !escapes {
					results <- 1
				} else {
					results <- 0
				}
			}(obstaclePosition)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		for result := range results {
			sneakyGuardPathCount += result
		}

		return sneakyGuardPathCount
	}
	grid, startingPosition, err := parseGrid(strings.NewReader(input))
	if err != nil {
		panic(err)
	}
	pathCoords, _ := simulateMovementUntilEscape(grid, startingPosition, North, Position{-1, -1})
	return len(pathCoords)
}
