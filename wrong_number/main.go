package main

import "fmt"

var allNumber = [][]float64{
	{40610, 41721, 42832, 43943, 45054, 46156, 47276},
	{91350, 90500, 89650, 88800, 87000, 87100, 86250},
	{12345, 17350, 22355, 26360, 32365, 37370, 42375},
	{10, 20, 30, 44, 50, 60, 70},
	{100, 200, 300, 400, 505, 600, 700},
	{0, 2, 3, 4, 5, 6, 7},
	{1, 2, 3, 4, 5, 6, 8},
}

type wrongNumberResult struct {
	defaultDiff float64
	wrongNrVal  float64
	wrongNrPos  int
	correctNr   float64
}

func (wn wrongNumberResult) String() string {
	return fmt.Sprintf("default difference: %v, wrong number: %v, wrong number position: %v, correct number: %v",
		wn.defaultDiff, wn.wrongNrVal, wn.wrongNrPos, wn.correctNr)
}

func diff(numbers []float64) []float64 {
	diffs := make([]float64, len(numbers)-1)
	for i := 0; i < len(numbers)-1; i++ {
		diffs[i] = numbers[i+1] - numbers[i]
	}

	return diffs
}

func mostFrequentDiff(diffs []float64) float64 {
	diffMap := make(map[float64]int)
	for _, diff := range diffs {
		diffMap[diff]++
	}

	freqCnt := 0
	freqVal := 0.0
	for v, c := range diffMap {
		if freqCnt < c {
			freqCnt = c
			freqVal = v
		}
	}

	return freqVal
}

func findWrongNumber(numbers []float64) wrongNumberResult {
	diffs := diff(numbers)

	// y = mx + n
	m := mostFrequentDiff(diffs)
	n := numbers[0]
	if numbers[1]-numbers[0] != m { // first or second index are the wrong number
		n = numbers[2] - 2*m // cannot use second index, it will interfere with first index
	}

	badIdx := 0
	for x := range numbers {
		if numbers[x] != m*float64(x)+n {
			badIdx = x
			break
		}
	}

	return wrongNumberResult{
		defaultDiff: m,
		wrongNrVal:  numbers[badIdx],
		wrongNrPos:  badIdx + 1,
		correctNr:   m*float64(badIdx) + n,
	}
}

func main() {
	for _, numbers := range allNumber {
		fmt.Println("numbers are:", numbers)
		res := findWrongNumber(numbers)
		fmt.Println(res)
		fmt.Println()
	}
}
