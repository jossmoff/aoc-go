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
		numstr := strings.Split(ls.Text(), "   ")
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

func distanceBetweenList(l []int, r []int) int {

	zipped := utils.Zip(slices.Values(l), slices.Values(r))
	diffs := utils.Map(func(v utils.Zipped[int, int]) int {
		return utils.Abs(v.V1 - v.V2)
	}, zipped)
	return utils.Sum(diffs)
}

func similarityScore(l []int, r []int) int {
	var rightCountMap map[int]int = make(map[int]int)
	for _, v := range r {
		if _, ok := rightCountMap[v]; !ok {
			rightCountMap[v] = 0
		}
		rightCountMap[v]++
	}

	return utils.Reduce(func(acc int, v int) int {
		return acc + v*rightCountMap[v]
	}, 0, slices.Values(l))
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
	leftList := utils.Map(func(v []int) int {
		return v[0]
	}, slices.Values(lines))

	rightList := utils.Map(func(v []int) int {
		return v[1]
	}, slices.Values(lines))

	if part2 {
		return similarityScore(slices.Collect(leftList), slices.Collect(rightList))
	}
	sortedLeftList := slices.Sorted(leftList)
	sortedRightList := slices.Sorted(rightList)

	return distanceBetweenList(sortedLeftList, sortedRightList)
}
