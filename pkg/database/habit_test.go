package database

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/metagunner/habheath/pkg/models"
	"github.com/metagunner/habheath/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const TestYear = 2024
const TestMonth = time.January

var testDB *DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	db, _ := SetupTestDB()
	testDB = db
	SeedTestData(ctx, testDB, TestYear, TestMonth)
	code := m.Run()
	if err := testDB.Close(); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func TestHabitService_GetAllByDay(t *testing.T) {
	service := NewHabitService(testDB)

	day := time.Date(TestYear, TestMonth, 1, 0, 0, 0, 0, time.UTC)

	chain, err := service.GetAllByDay(context.Background(), day)

	assert.NoError(t, err)
	assert.NotEmpty(t, chain.Habits)
	assert.Equal(t, chain.Title, day.Format(time.DateOnly))
}

func TestHabitService_Create(t *testing.T) {
	service := NewHabitService(testDB)

	habitTitle, _ := models.CreateHabitTitle("Play chess")
	day := time.Date(TestYear, TestMonth, 1, 0, 0, 0, 0, time.UTC)
	habit, _ := models.CreateHabit(habitTitle, day, false)

	err := service.Create(context.Background(), habit)

	assert.NoError(t, err)
	assert.NotZero(t, habit.Id)
}

func TestHabitService_Delete(t *testing.T) {
	service := NewHabitService(testDB)
	t.Run("Given not exist habit id should fail", func(t *testing.T) {
		notExistsHabitId := models.HabitId(9999999)
		err := service.Delete(context.Background(), notExistsHabitId)

		assert.ErrorIs(t, err, ErrHabitNotFound)
	})
	t.Run("Given valid habit id should succeed", func(t *testing.T) {
		day := time.Date(TestYear, TestMonth, 1, 0, 0, 0, 0, time.UTC)
		chain, _ := service.GetAllByDay(context.Background(), day)
		err := service.Delete(context.Background(), chain.Habits[0].Id)

		assert.NoError(t, err)
	})
}

func TestHabitService_Update(t *testing.T) {
	service := NewHabitService(testDB)

	day := time.Date(TestYear, TestMonth, 2, 0, 0, 0, 0, time.UTC)
	chain, _ := service.GetAllByDay(context.Background(), day)

	habit := chain.Habits[0]
	isCompleted := habit.IsCompleted
	habit.ChangeTitle("Updated Habit")
	habit.ToggleCompletion()

	err := service.Update(context.Background(), habit)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Habit", habit.Title.String())
	assert.Equal(t, !isCompleted, habit.IsCompleted)
}

func TestHabitService_HeathMap(t *testing.T) {
	service := NewHabitService(testDB)
	ctx := context.Background()

	// just for test
	testYear := 1990

	title, _ := models.CreateHabitTitle("Habit 1")
	habit, _ := models.CreateHabit(title, utils.CreateDate(testYear, 1, 1), true)
	service.Create(ctx, habit)

	title2, _ := models.CreateHabitTitle("Habit 2")
	habit2, _ := models.CreateHabit(title2, utils.CreateDate(testYear, 1, 2), false)
	service.Create(ctx, habit2)

	from := utils.CreateDate(testYear, 1, 1)
	to := utils.CreateDate(testYear, 12, 31)

	heatMap, count, err := service.HeatMap(ctx, from, to)
	assert.NoError(t, err)
	assert.NotNil(t, heatMap)

	expected := map[time.Time]*models.HeathMap{
		utils.CreateDate(testYear, 1, 1): {TotalNumberOfHabits: 1, CompletedHabits: 1, Day: 1, Month: 1, Year: testYear},
		utils.CreateDate(testYear, 1, 2): {TotalNumberOfHabits: 1, CompletedHabits: 0, Day: 2, Month: 1, Year: testYear},
	}

	assert.Equal(t, len(expected), count)

	// Verify the heat map results
	for date, expectedHeathMap := range expected {
		actualHeathMap, exists := heatMap[date]
		assert.True(t, exists, "Date %v not found in heat map", date)
		assert.Equal(t, expectedHeathMap.TotalNumberOfHabits, actualHeathMap.TotalNumberOfHabits)
		assert.Equal(t, expectedHeathMap.CompletedHabits, actualHeathMap.CompletedHabits)
		assert.Equal(t, expectedHeathMap.Day, actualHeathMap.Day)
		assert.Equal(t, expectedHeathMap.Month, actualHeathMap.Month)
		assert.Equal(t, expectedHeathMap.Year, actualHeathMap.Year)
	}
}
