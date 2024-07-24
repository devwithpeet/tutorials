package main

import (
	"testing"
	"time"
)

func Test_calc(t *testing.T) {
	today, _ := time.Parse(time.DateOnly, "2024-05-20")

	t.Run("Today's year old", func(t *testing.T) {
		input := today.Format(time.DateOnly)
		want := int64(0)
		if got := calc(today, input); got != want {
			t.Errorf("calc() = %v, want %v", got, want)
		}
	})

	t.Run("1 day old", func(t *testing.T) {
		input := today.Add(-24 * time.Hour).Format(time.DateOnly)
		want := int64(1)
		if got := calc(today, input); got != want {
			t.Errorf("calc() = %v, want %v", got, want)
		}
	})

	t.Run("1 year old", func(t *testing.T) {
		input := today.Add(-365 * 24 * time.Hour).Format(time.DateOnly)
		want := int64(365)
		if got := calc(today, input); got != want {
			t.Errorf("calc() = %v, want %v", got, want)
		}
	})

	t.Run("Peter's birthday", func(t *testing.T) {
		input := "1981-08-17"
		want := int64(15617)
		if got := calc(today, input); got != want {
			t.Errorf("calc() = %v, want %v", got, want)
		}
	})
}
