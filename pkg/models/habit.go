package models

import (
	"context"
	"time"

	"github.com/metagunner/habheath/pkg/app"
)

var ErrInvalidHabitTitle = app.Errorf(app.EINVALID, "Invalid habit title.")

type Habit struct {
	Id          HabitId    `json:"id"`
	Title       HabitTitle `json:"title"`
	Day         time.Time  `json:"day"`
	IsCompleted bool       `json:"is_completed"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type (
	HabitId    int
	HabitTitle string
)

func CreateHabitTitle(title string) (HabitTitle, error) {
	if title == "" || len(title) > 250 {
		return "", ErrInvalidHabitTitle
	}

	return HabitTitle(title), nil
}

func (ht HabitTitle) String() string {
	return string(ht)
}

func CreateHabit(title HabitTitle, day time.Time, isCompleted bool) (*Habit, error) {
	day = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	habit := &Habit{Title: title, Day: day, UpdatedAt: time.Now().UTC(), IsCompleted: isCompleted}
	return habit, nil
}

func (h *Habit) ToggleCompletion() {
	h.IsCompleted = !h.IsCompleted
	h.UpdatedAt = time.Now().UTC()
}

func (h *Habit) ChangeTitle(title string) error {
	if h.Title.String() == title {
		return nil
	}

	habitTitle, err := CreateHabitTitle(title)
	if err != nil {
		return err
	}

	h.Title = habitTitle
	h.UpdatedAt = time.Now().UTC()
	return nil
}

type HabitService interface {
	// All the habits for heath map
	HeatMap(ctx context.Context, from time.Time, to time.Time) (map[time.Time]*HeathMap, int, error)
	// Get all the habits for the given day
	GetAllByDay(ctx context.Context, day time.Time) (*Chain, error)
	Create(ctx context.Context, habit *Habit) error
	// delete a habit
	Delete(ctx context.Context, id HabitId) error
	Update(ctx context.Context, habit *Habit) error
}

type Chain struct {
	Title  string
	Habits []*Habit
}

type HeathMap struct {
	TotalNumberOfHabits int
	CompletedHabits     int
	Day                 int
	Month               int
	Year                int
}
