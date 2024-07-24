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

	year, month, day := calculateTimeSpent(now, birthday)

	fmt.Printf("Peter was %d years, %d months, %d days old on %s.\n", year, month, day, birthday)
}

func calculateTimeSpent(now time.Time, birthday string) (int, int, int) {
	split := strings.Split(birthday, "-")

	bornYear, _ := strconv.Atoi(split[0])
	bornMonth, _ := strconv.Atoi(split[1])
	bornDay, _ := strconv.Atoi(split[2])

	currentYear := now.Year()
	currentMonth := int(now.Month())
	currentDay := now.Day()

	if bornYear > currentYear {
		return -1, -1, -1
	}

	if bornYear == currentYear && bornMonth > currentMonth {
		return -1, -1, -1
	}

	if bornYear == currentYear && bornMonth == currentMonth && bornDay > currentDay {
		return -1, -1, -1
	}

	days := currentDay - bornDay
	if days < 0 {
		bornMonthLength := 30
		// January, March, May, July, August, October, December are 31 days long
		if bornMonth == 1 || bornMonth == 3 || bornMonth == 5 || bornMonth == 7 || bornMonth == 8 || bornMonth == 10 || bornMonth == 12 {
			bornMonthLength = 31
		} else if bornMonth == 2 {
			bornMonthLength = 28
			// Leap year: February is 29 days long if the year is divisible by 4 but not divisible by 100, or also if the year is divisible by 400.
			// Examples of leap years: 1244, 1600, 2000
			// Examples of non-leap years: 2011, 1700, 1800,
			if currentYear%4 == 0 || currentYear%100 == 0 && currentYear%400 != 0 {
				bornMonthLength = 29
			}
		}

		days = bornMonthLength + days
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

	if bornYear > currentYear {
		return -1, -1, -1
	}

	return currentYear - bornYear, months, days
}
