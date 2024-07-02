package habheath

import (
	"context"
	"time"
)

var (
	ErrInvalidHabitTitle = Errorf(EINVALID, "Invalid habit title.")
)

type Habit struct {
	Id          HabitId    `json:"id"`
	Title       HabitTitle `json:"title"`
	Day         time.Time  `json:"day"`
	IsCompleted bool       `json:"is_completed"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type HabitId int
type HabitTitle string

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
	HeatMap(ctx context.Context, year int) ([]*HeathMap, int, error)
	// Don't break the chain list
	MontlyChain(ctx context.Context, year int) ([]*DontBreakTheChain, int, error)
	// Get all the habits for the given day
	GetAllByDay(ctx context.Context, day time.Time) (*Chain, error)
	Create(ctx context.Context, habit *Habit) error
	// delete a habit
	Delete(ctx context.Context, id HabitId) error
	Update(ctx context.Context, habit *Habit) error
}

type Chain struct {
	//HabitTitle(time.Now().Format(time.DateOnly)), nil
	// 12/2 habits on June 22nd
	// No habits on August 26th
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

type DontBreakTheChain struct {
	Day     int
	Month   int
	Year    int
	IsBreak bool
}
