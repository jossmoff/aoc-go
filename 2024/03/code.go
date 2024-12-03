package main

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jossmoff/aoc-go/utils"
	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		r := regexp.MustCompile(`do\(\)|don't\(\)|mul\((\d{1,3},\d{1,3})\)`)
		matches := r.FindAllStringSubmatch(input, -1)
		for _, match := range matches {
			fmt.Println(match)
		}
	}
	r := regexp.MustCompile(`mul\((\d{1,3},\d{1,3})\)`)
	matches := r.FindAllStringSubmatch(input, -1)
	multipliedPairs := utils.Map(func(s string) int {
		parts := strings.Split(s, ",")
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		return a * b
	}, utils.Map(func(match []string) string {
		return match[1]
	}, slices.Values(matches)))
	return utils.Reduce(func(acc int, v int) int {
		return acc + v
	}, 0, multipliedPairs)
}
