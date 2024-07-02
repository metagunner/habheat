package database

import (
	"context"
	"fmt"
	"habheath"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TestYear = 2024
const TestMonth = time.January

var testDB *DB

func seedTestData(db *DB) {
	var b strings.Builder
	b.WriteString("INSERT INTO habit (title, day, is_completed, updated_at) VALUES")

	maxHabitsPerDay := 3
	var values []interface{}
	const monthDayCount = 30
	for i := 1; i <= monthDayCount; i++ {
		title, _ := habheath.CreateHabitTitle(fmt.Sprintf("Habit %d", i))
		day := time.Date(TestYear, TestMonth, i, 0, 0, 0, 0, time.UTC)
		habitsPerDay := rand.Intn(maxHabitsPerDay) + 1
		for j := 0; j < habitsPerDay; j++ {
			isCompleted := rand.Intn(2) == 1
			habit, _ := habheath.CreateHabit(title, day, isCompleted)

			if j > 0 || (j == 0 && i > 1) {
				b.WriteString(", ")
			}

			b.WriteString("(?, ?, ?, ?)")
			values = append(values, habit.Title, habit.Day.Format(time.RFC3339), habit.IsCompleted, habit.UpdatedAt.Format(time.RFC3339))
		}
	}

	bulkInsert := b.String()
	db.db.ExecContext(context.Background(), bulkInsert, values...)
}

func TestMain(m *testing.M) {
	db, _ := SetupTestDB()
	testDB = db
	seedTestData(db)
	m.Run()
	CloseTestDB(db)
}

func TestHabitService_Get(t *testing.T) {
	service := NewHabitService(testDB)

	day := time.Date(TestYear, TestMonth, 1, 0, 0, 0, 0, time.UTC)

	chain, err := service.GetAllByDay(context.Background(), day)

	assert.NoError(t, err)
	assert.NotEmpty(t, chain.Habits)
	assert.Equal(t, chain.Title, day.Format(time.DateOnly))
}

func TestHabitService_Create(t *testing.T) {
	service := NewHabitService(testDB)

	habitTitle, _ := habheath.CreateHabitTitle("Play chess")
	day := time.Date(TestYear, TestMonth, 1, 0, 0, 0, 0, time.UTC)
	habit, _ := habheath.CreateHabit(habitTitle, day, false)

	err := service.Create(context.Background(), habit)

	assert.NoError(t, err)
	assert.NotZero(t, habit.Id)
}

func TestHabitService_Delete(t *testing.T) {
	service := NewHabitService(testDB)
	t.Run("Given not exist habit id should fail", func(t *testing.T) {
		notExistsHabitId := habheath.HabitId(9999999)
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
