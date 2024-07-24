package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	now := time.Now()

	const birthday = "1981-08-17"

	year, month, day, err := calc(now, birthday)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Peter was %d years, %d months, %d days old on %s.\n", year, month, day, birthday)
}

func calc(now time.Time, birthday string) (int, int, int, error) {
	split := strings.Split(birthday, "-")

	bornYear, _ := strconv.Atoi(split[0])
	bornMonth, _ := strconv.Atoi(split[1])
	bornDay, _ := strconv.Atoi(split[2])

	currentYear := now.Year()
	currentMonth := int(now.Month())
	currentDay := now.Day()

	if bornYear > currentYear {
		return 0, 0, 0, fmt.Errorf("Peter was born in the future.")
	}

	if bornYear == currentYear && bornMonth > currentMonth {
		return 0, 0, 0, fmt.Errorf("Peter was born in the future.")
	}

	if bornYear == currentYear && bornMonth == currentMonth && bornDay > currentDay {
		return 0, 0, 0, fmt.Errorf("Peter was born in the future.")
	}

	days := currentDay - bornDay
	if days < 0 {
		days = getMonthLength(currentYear, bornMonth) + days
		if bornMonth == 12 {
			bornMonth = 1
			bornYear += 1
		} else {
			bornMonth += 1
		}
	}

	months := currentMonth - bornMonth
	if months < 0 {
		months = 12 + months
		bornYear += 1
	}

	return currentYear - bornYear, months, days, nil
}

func getMonthLength(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 2:
		if isLeapYear(year) {
			return 29
		}

		return 28
	}

	return 30
}

// Leap year: February is 29 days long if the year is divisible by 4 but not divisible by 100, or also if the year is divisible by 400.
// Examples of leap years: 1244, 1600, 2000
// Examples of non-leap years: 2011, 1700, 1800,
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
