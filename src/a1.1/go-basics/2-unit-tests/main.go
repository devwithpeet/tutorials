package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()

	days := calc(now, "1981-08-17")

	fmt.Printf("%d\n", days)
}

func calc(now time.Time, input string) int64 {
	born, _ := time.Parse(time.DateOnly, input)

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	diff := today.Unix() - born.Unix()

	return diff / 24 / 60 / 60
}
