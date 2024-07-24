package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getFileName() string {
	if len(os.Args) < 2 {
		return "a1.1/6-generate-random-salary-data/salaries.csv"
	}

	return os.Args[1]
}

func main() {
	dat, err := os.ReadFile(getFileName())
	check(err)

	data := split(string(dat))
	result := stats(data)
	for _, v := range result {
		fmt.Println(v[0], v[1])
	}
}

func split(s string) [][]string {
	var result [][]string

	for _, line := range strings.Split(s, "\n") {
		result = append(result, strings.Split(line, ","))
	}

	return result
}

func stats(data [][]string) [][2]string {
	var (
		positions []string
		numbers   [][]int
	)

	const positionIdx = 1
	const salaryIdx = 4

	for _, v := range data {
		if len(v) < 5 {
			continue
		}

		position := v[positionIdx]
		idx := -1
		for i, p := range positions {
			if p == position {
				idx = i
				break
			}
		}
		if idx == -1 {
			positions = append(positions, position)
			numbers = append(numbers, []int{})
			idx = len(positions) - 1
		}

		salary, err := strconv.Atoi(v[salaryIdx])
		check(err)

		numbers[idx] = append(numbers[idx], salary)
	}

	var result [][2]string
	for i, p := range positions {
		avg := average(numbers[i])
		result = append(result, [2]string{p, strconv.Itoa(avg)})
	}

	return result
}

func average(numbers []int) int {
	sum := 0
	for _, v := range numbers {
		sum += v
	}

	return sum / len(numbers)
}
