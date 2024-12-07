package main

import (
	"bufio"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/jossmoff/aoc-go/utils"
	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

func parseInput(r io.Reader) ([][]int, error) {

	var result [][]int
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		target, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, err
		}

		numsStr := strings.Fields(parts[1])
		nums := make([]int, len(numsStr)+1)
		nums[0] = target
		for i, numStr := range numsStr {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, err
			}
			nums[i+1] = num
		}

		result = append(result, nums)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

type Calibration struct {
	target         int
	lastValueIndex int
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

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func isValidCalibration(originalTarget int, calibrationValues []int) int {
	stack := arraystack.New()

	last := len(calibrationValues) - 1

	calibration := Calibration{originalTarget, last}
	stack.Push(calibration)

	for stack.Size() > 0 {
		calibrationToCheck, _ := stack.Pop()
		lastValueIndex := calibrationToCheck.(Calibration).lastValueIndex
		if lastValueIndex < 0 {
			continue
		}
		lastValue := calibrationValues[lastValueIndex]
		target := calibrationToCheck.(Calibration).target

		// Check minus
		if target-lastValue == 0 {
			return originalTarget
		}
		if target-lastValue > 0 {
			stack.Push(Calibration{target - lastValue, lastValueIndex - 1})
		}

		// Check divide
		if target%lastValue == 0 {
			if target/lastValue == 1 {
				return originalTarget
			} else {
				stack.Push(Calibration{target / lastValue, lastValueIndex - 1})
			}
		}

		// Check concat
		numDigitsOfLastValue := countDigits(lastValue)
		if (target-lastValue)%powInt(10, numDigitsOfLastValue) == 0 {
			nextNum := (target - lastValue) / powInt(10, numDigitsOfLastValue)
			if lastValueIndex-1 == 0 && nextNum == calibrationValues[lastValueIndex-1] {
				return originalTarget
			}
			stack.Push(Calibration{nextNum, lastValueIndex - 1})
		}
	}

	return -1
}

func run(part2 bool, input string) any {
	calibrations, _ := parseInput(strings.NewReader(input))
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		return "not implemented"
	}
	// solve part 1 here
	return utils.Reduce(func(acc int, calibration []int) int {
		target := calibration[0]
		calibrationValues := calibration[1:]
		value := isValidCalibration(target, calibrationValues)
		if value != -1 {
			return acc + value
		} else {
			return acc
		}
	}, 0, slices.Values(calibrations))
}
