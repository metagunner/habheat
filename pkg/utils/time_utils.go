package utils

import (
	"fmt"
	"slices"
	"time"
)

func GetYearsBetween(from int, to int) []int {
	years := make([]int, 0, to-from+1)
	for i := from; i <= to; i++ {
		years = append(years, i)
	}
	slices.Reverse(years)
	return years
}

func GetMonths(first time.Time) []time.Time {
	months := make([]time.Time, 0)
	nextMonth := first
	for i := 0; len(months) != 12; i++ {
		months = append(months, nextMonth)
		nextMonth = nextMonth.AddDate(0, 1, 0)
	}

	return months
}

// Returns the number of days in the given month of the specified year
func GetDaysInMonth(t time.Time) int {
	nextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	// Subtract one day to get the last day of the current month
	lastDayOfMonth := nextMonth.AddDate(0, 0, -1)
	return lastDayOfMonth.Day()
}

func CreateDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// Returns the ordinal suffix for a given day
func GetOrdinalSuffix(day int) string {
	if day >= 11 && day <= 13 {
		return fmt.Sprintf("%dth", day)
	}

	switch day % 10 {
	case 1:
		return fmt.Sprintf("%dst", day)
	case 2:
		return fmt.Sprintf("%dnd", day)
	case 3:
		return fmt.Sprintf("%drd", day)
	default:
		return fmt.Sprintf("%dth", day)
	}
}
