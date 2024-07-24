package main

import (
	"testing"
	"time"
)

func Test_calculateTimeSpent(t *testing.T) {
	today, _ := time.Parse(time.DateOnly, "2024-05-20")

	tests := []struct {
		name       string
		birthday   string
		wantYears  int
		wantMonths int
		wantDays   int
		wantError  bool
	}{
		{
			name:       "Today's year old",
			birthday:   today.Format(time.DateOnly),
			wantYears:  0,
			wantMonths: 0,
			wantDays:   0,
		},
		{
			name:       "Very young",
			birthday:   today.Add(-24 * time.Hour).Format(time.DateOnly),
			wantYears:  0,
			wantMonths: 0,
			wantDays:   1,
		},
		{
			name:       "Brian Kernighan",
			birthday:   "1942-01-30",
			wantYears:  82,
			wantMonths: 3,
			wantDays:   21,
		},
		{
			name:       "Linux Torvalds",
			birthday:   "1969-12-28",
			wantYears:  54,
			wantMonths: 4,
			wantDays:   23,
			wantError:  false,
		},
		{
			name:       "Peter's birthday",
			birthday:   "1981-08-17",
			wantYears:  42,
			wantMonths: 9,
			wantDays:   3,
			wantError:  false,
		},
		// Logo
		{
			name:       "Seymour Papert",
			birthday:   "1928-02-29",
			wantYears:  96,
			wantMonths: 2,
			wantDays:   20,
			wantError:  false,
		},
		{
			name:       "Future birthday",
			birthday:   "2050-08-17",
			wantYears:  0,
			wantMonths: 0,
			wantDays:   0,
			wantError:  true,
		},
		{
			name:       "Future birthday 2",
			birthday:   "2024-06-17",
			wantYears:  0,
			wantMonths: 0,
			wantDays:   0,
			wantError:  true,
		},
		{
			name:       "Future birthday 3",
			birthday:   "2024-05-21",
			wantYears:  0,
			wantMonths: 0,
			wantDays:   0,
			wantError:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			born, err := time.Parse(time.DateOnly, tt.birthday)
			if err != nil {
				t.Errorf("time.Parse() error = %v", err)
				return
			}

			gotYears, gotMonths, gotDays, err := calculateTimeSpent(today, born)
			if tt.wantError && err == nil || !tt.wantError && err != nil {
				t.Errorf("calculateTimeSpent() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if gotYears != tt.wantYears {
				t.Errorf("calculateTimeSpent() gotYears = %v, wantYears %v", gotYears, tt.wantYears)
			}
			if gotMonths != tt.wantMonths {
				t.Errorf("calculateTimeSpent() gotMonths = %v, wantMonths %v", gotMonths, tt.wantMonths)
			}
			if gotDays != tt.wantDays {
				t.Errorf("calculateTimeSpent() gotDays = %v, wantDays %v", gotDays, tt.wantDays)
			}
		})
	}
}
