package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

func countDigits(n int) int {
	if n == 0 {
		return 1
	}
	count := 0
	for n != 0 {
		n /= 10
		count++
	}
	return count
}

func splitStone(n int) (int, int) {
	str := strconv.Itoa(n)
	mid := len(str) / 2
	left, _ := strconv.Atoi(str[:mid])
	right, _ := strconv.Atoi(str[mid:])
	return left, right
}

func parseInput(r io.Reader) ([]int, error) {
	var result []int
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}
			result = append(result, num)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func zeroRule(n int) []int {
	if n == 0 {
		return []int{1}
	}
	return []int{}
}

func evenDigitsRule(n int) []int {
	if countDigits(n)%2 == 0 {
		left, right := splitStone(n)
		return []int{left, right}
	}
	return []int{}
}

func defaultRule(n int) []int {
	return []int{n * 2024}
}

func hashCode(n, d int) uint64 {
	return uint64(n)<<32 | uint64(d)
}

func traverse(n, d, dMax int, cacheHit, cacheMiss *int, cache map[uint64]int) int {
	if v, ok := cache[hashCode(n, d)]; ok {
		(*cacheHit)++
		return v
	} else {
		(*cacheMiss)++
	}
	if d == dMax {
		return 1
	}
	acc := 0
	for _, rule := range []func(int) []int{zeroRule, evenDigitsRule, defaultRule} {
		newStones := rule(n)
		if len(newStones) > 0 {
			for _, stone := range newStones {
				val := traverse(stone, d+1, dMax, cacheHit, cacheMiss, cache)
				acc += val
				cache[hashCode(stone, d+1)] = val
			}
			break
		}

	}
	return acc
}

func run(part2 bool, input string) any {

	if part2 {
		return "not implemented"
	}

	stones, err := parseInput(strings.NewReader(input))
	if err != nil {
		panic(err)
	}
	blinks := 75

	acc := 0
	cacheHit := 0
	cacheMiss := 0

	cache := make(map[uint64]int)
	stones, err = parseInput(strings.NewReader(input))
	for _, stone := range stones {
		acc += traverse(stone, 0, blinks, &cacheHit, &cacheMiss, cache)
	}

	fmt.Println("cache hit rate:", float64(cacheHit)/float64(cacheHit+cacheMiss)*100)

	return acc
}
