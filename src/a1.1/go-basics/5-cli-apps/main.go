package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	now := time.Now()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("What is your name?")
		scanner.Scan()
		name := strings.Trim(scanner.Text(), " ")
		if shouldExit(name) {
			break
		}

		fmt.Println("What is your birthday? (YYYY-MM-DD)")
		scanner.Scan()
		birthday := strings.Trim(scanner.Text(), " ")
		if shouldExit(birthday) {
			break
		}

		born, err := parseBirthday(birthday)
		if err != nil {
			fmt.Println(err)
			return
		}
		year, month, day, err := calculateTimeSpent(now, born)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s was %d years, %d months, %d days old on %s.\n", name, year, month, day, birthday)
	}

	fmt.Println("Bye!")
}

func shouldExit(name string) bool {
	switch name {
	case "exit", "e", "quit", "q", "":
		return true
	}

	return false
}

func parseBirthday(birthday string) (time.Time, error) {
	bDay, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return time.Time{}, err
	}

	return bDay, nil
}

func calculateTimeSpent(now, born time.Time) (int, int, int, error) {
	if born.After(now) {
		return 0, 0, 0, errors.New("born in the future")
	}

	currentYear, cMonth, _ := now.Date()
	bornYear, bMonth, _ := born.Date()
	days := now.Day() - born.Day()

	bornMonth := int(bMonth)
	currentMonth := int(cMonth)

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
