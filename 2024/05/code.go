package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	sets "github.com/jossmoff/aoc-go/utils/ds"
	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

func parseInput(r io.Reader) (map[int]*sets.Set[int], [][]int, error) {
	scanner := bufio.NewScanner(r)
	dependencies := make(map[int]*sets.Set[int])
	var runs [][]int
	isRunSection := false

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			isRunSection = true
			continue
		}

		if !isRunSection {
			parts := strings.Split(line, "|")
			if len(parts) != 2 {
				return nil, nil, fmt.Errorf("invalid dependency format: %s", line)
			}
			key, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, nil, err
			}
			values := strings.Split(parts[1], ",")
			for _, v := range values {
				val, err := strconv.Atoi(v)
				if err != nil {
					return nil, nil, err
				}
				if _, exists := dependencies[key]; !exists {
					dependencies[key] = sets.NewSet[int]()
				}
				sets.Add(dependencies[key], val)
			}
		} else {
			values := strings.Split(line, ",")
			var run []int
			for _, v := range values {
				val, err := strconv.Atoi(v)
				if err != nil {
					return nil, nil, err
				}
				run = append(run, val)
			}
			runs = append(runs, run)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return dependencies, runs, nil
}

func run(part2 bool, input string) any {
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return "not implemented"
	}
	dependencySetMap, updates, err := parseInput(strings.NewReader(input))
	if err != nil {
		return err
	}
	println(updates)
	for _, update := range updates {
		seenPageDependencies := sets.NewSet[int]()
		for _, page := range update {
			pageDependencies := dependencySetMap[page]
			var seenAllDependencies = true
			for _, dep := range sets.List(pageDependencies) {
				seenAllDependencies = seenAllDependencies && sets.Contains(seenPageDependencies, dep)
			}

		}
	}
	return 42
}
