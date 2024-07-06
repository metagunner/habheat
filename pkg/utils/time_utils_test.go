package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetYearsBetween(t *testing.T) {
	from := 2022
	to := 2024
	years := GetYearsBetween(from, to)

	assert.Equal(t, 3, len(years))
	assert.Equal(t, 2024, years[0])
	assert.Equal(t, 2023, years[1])
	assert.Equal(t, 2022, years[2])
}

func TestGetMonths(t *testing.T) {
	t.Run("Given start of the year should give all months", func(t *testing.T) {
		date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		months := GetMonths(date)

		assert.Equal(t, 12, len(months))

		expected := date.AddDate(0, -11, 0)
		for _, month := range months {
			assert.Equal(t, expected.Format("Jan"), month.Format("Jan"))
			expected = expected.AddDate(0, 1, 0)
		}
	})

	t.Run("Given middle of the year should include months from previous year", func(t *testing.T) {
		date := time.Date(2024, 3, 7, 0, 0, 0, 0, time.UTC)
		months := GetMonths(date)

		assert.Equal(t, 12, len(months))

		expected := date.AddDate(0, -11, 0)
		for _, month := range months {
			assert.Equal(t, expected.Format("Jan"), month.Format("Jan"))
			expected = expected.AddDate(0, 1, 0)
		}
	})
}

func TestGetDaysInMonth(t *testing.T) {
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	monthDayCount := GetDaysInMonth(date)

	assert.Equal(t, 31, monthDayCount)
}

func TestGetOrdinalSuffix(t *testing.T) {
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "1st", GetOrdinalSuffix(date.Day()))
	date = time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "22nd", GetOrdinalSuffix(date.Day()))
}
