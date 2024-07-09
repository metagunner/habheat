package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/metagunner/habheath/pkg/database"
	"github.com/metagunner/habheath/pkg/gui"
)

func main() {
	// SuperGrid()
	db := database.NewDB("./test.db")
	if err := db.Open(); err != nil {
		panic(err)
	}
	database.SeedTestData(context.Background(), db, 2023, 7)

	gui := gui.NewGui(db)
	err := gui.Run()
	if err != nil {
		if !errors.Is(err, gocui.ErrQuit) {
			panic(err)
		}
	}
}

func SuperGrid() {
	// Get today's date
	today := time.Now().UTC()

	if today.Weekday() == time.Sunday {
		today = today.AddDate(0, 0, 7)
	}

	// Find the next upcoming Sunday after today
	for today.Weekday() != time.Sunday {
		today = today.AddDate(0, 0, 1)
	}

	// Calculate the start date: 53 weeks before the next upcoming Sunday
	startDate := today.AddDate(0, 0, -53*7)

	// Find the first Sunday on or after the start date
	firstSunday := startDate
	for firstSunday.Weekday() != time.Sunday {
		firstSunday = firstSunday.AddDate(0, 0, 1)
	}
	startDate = firstSunday

	// Generate random shades of green (30-37 ANSI codes)
	shades := []int{30, 32, 34, 36} // ANSI codes for different shades of green

	// Create the grid
	grid := make([][]string, 7)
	for i := range grid {
		grid[i] = make([]string, 53)
	}

	// Fill the grid with colored boxes
	today = time.Now().UTC()
	currentDate := startDate
	for col := 0; col < 53; col++ {
		for row := 0; row < 7; row++ {
			if currentDate.After(today) {
				grid[row][col] = fmt.Sprintf("\x1b[48;5;52m%d \x1b[0m", currentDate.Day())
			} else {
				// Select a random shade of green
				colorIndex := rand.Intn(len(shades))
				// ANSI escape sequence for coloring (green shades)
				colorCode := shades[colorIndex]
				grid[row][col] = fmt.Sprintf("\x1b[48;5;%dm%d \x1b[0m", colorCode, currentDate.Day())
			}
			currentDate = currentDate.AddDate(0, 0, 1) // move to the next day
		}
	}

	// Print the grid
	for _, row := range grid {
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println()
	}
}
