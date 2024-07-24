package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()

	born, _ := time.Parse(time.DateOnly, "1981-08-17")

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	secondsSinceBorn := today.Unix() - born.Unix()

	secondsInADay := 24 * 60 * 60

	daysSinceBorn := secondsSinceBorn / int64(secondsInADay)

	fmt.Println(daysSinceBorn)
}
