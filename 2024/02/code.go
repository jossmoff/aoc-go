package main

import (
	"bufio"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/jossmoff/aoc-go/utils"
	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

func readInput(r io.Reader) [][]int {
	ls := bufio.NewScanner(r)

	var lines [][]int
	for ls.Scan() {
		numstr := strings.Split(ls.Text(), " ")
		nums := utils.Map(func(v string) int {
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalln(err)
			}
			return n
		}, slices.Values(numstr))
		lines = append(lines, slices.Collect(nums))
	}

	return lines
}

func isLevelSafe(l []int) bool {
	if len(l) < 2 {
		return true
	}
	isIncreasing := l[0] <= l[1]
	for i := 1; i < len(l); i++ {
		diff := l[i] - l[i-1]
		if isIncreasing && diff <= 0 || !isIncreasing && diff >= 0 {
			return false
		}
		absDiff := utils.Abs(diff)
		if absDiff > 3 {
			return false
		}
	}
	return true
}

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	// when you're ready to do part 2, remove this "not implemented" block
	var lines [][]int = readInput(strings.NewReader(input))
	if part2 {
		valid := utils.Filter(func(line []int) bool {
			var variants [][]int
			for i := range len(line) {
				variants = append(variants, slices.Delete(slices.Clone(line), i, i+1))
			}
			return utils.Any(isLevelSafe, slices.Values(variants))
		}, slices.Values(lines))
		return len(slices.Collect(valid))
	}

	valid := utils.Filter(func(line []int) bool {
		return isLevelSafe(line)
	}, slices.Values(lines))

	return len(slices.Collect(valid))
}
