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
	parsed := parse(data)
	entries := stats(parsed)
	for _, e := range entries {
		fmt.Println(e.Position, e.Count, e.Average)
	}
}

func split(s string) [][]string {
	var result [][]string

	for _, line := range strings.Split(s, "\n") {
		result = append(result, strings.Split(line, ","))
	}

	return result
}

type PositionSalaries struct {
	Position string
	Salaries []int
}

func parse(data [][]string) []PositionSalaries {
	var (
		result []PositionSalaries
	)

	const positionIdx = 1
	const salaryIdx = 4

	for _, v := range data {
		if len(v) < 5 {
			continue
		}

		position := v[positionIdx]
		idx := -1
		for i, s := range result {
			if s.Position == position {
				idx = i
				break
			}
		}
		if idx == -1 {
			result = append(result, PositionSalaries{
				Position: position,
				Salaries: []int{},
			})
			idx = len(result) - 1
		}

		salary, err := strconv.Atoi(v[salaryIdx])
		check(err)

		result[idx].Salaries = append(result[idx].Salaries, salary)
	}

	return result
}

type Stats struct {
	Position string
	Average  int
	Count    int
}

func stats(positionSalaries []PositionSalaries) []Stats {
	var result []Stats
	for _, entry := range positionSalaries {
		result = append(result, Stats{
			Position: entry.Position,
			Average:  average(entry.Salaries),
			Count:    len(entry.Salaries),
		})
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
