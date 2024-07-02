package habheath_test

import (
	"habheath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var now = time.Now().UTC()

func TestCreateHabitTitle(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		hasTitle string
		err      error
	}{
		{"Given valid title should succeed", "Exercise", "Exercise", nil},
		{"Given empty title should fail", "", "", habheath.ErrInvalidHabitTitle},
		{"Given long title should fail", string(make([]byte, 251)), "", habheath.ErrInvalidHabitTitle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := habheath.CreateHabitTitle(tt.title)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, habheath.HabitTitle(tt.hasTitle), got)
		})
	}
}

func TestCreateHabit(t *testing.T) {
	title := habheath.HabitTitle("Exercise")
	isCompleted := false

	habit, err := habheath.CreateHabit(title, now, isCompleted)

	assert.NoError(t, err)
	assert.Equal(t, title, habit.Title)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), habit.Day)
	assert.Equal(t, isCompleted, habit.IsCompleted)
	assert.WithinDuration(t, time.Now().UTC(), habit.UpdatedAt, time.Second)
}

func TestToggleCompletion(t *testing.T) {
	title := habheath.HabitTitle("Exercise")
	day := time.Now().UTC()
	isCompleted := false

	habit, _ := habheath.CreateHabit(title, day, isCompleted)

	t.Run("toggle to complete", func(t *testing.T) {
		habit.ToggleCompletion()
		assert.True(t, habit.IsCompleted)
	})

	t.Run("toggle to incomplete", func(t *testing.T) {
		habit.ToggleCompletion()
		assert.False(t, habit.IsCompleted)
	})
}

func TestChangeTitle(t *testing.T) {
	tests := []struct {
		name         string
		initialTitle habheath.HabitTitle
		newTitle     string
		expectedErr  error
	}{
		{"change to valid title should succeed", "Exercise", "Read", nil},
		{"change to same title should succeed", "Exercise", "Exercise", nil},
		{"change to long title should fail", "Exercise", string(make([]byte, 251)), habheath.ErrInvalidHabitTitle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			habit, _ := habheath.CreateHabit(tt.initialTitle, time.Now().UTC(), false)
			err := habit.ChangeTitle(tt.newTitle)

			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, habheath.HabitTitle(tt.newTitle), habit.Title)
				assert.WithinDuration(t, time.Now().UTC(), habit.UpdatedAt, time.Second)
			}
		})
	}
}
