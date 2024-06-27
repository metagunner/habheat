package habheath

import (
	"context"
	"time"
)

var ErrHabitNotFound = Errorf(ENOTFOUND, "Habit not found.")

type ChainId int

type Chain struct {
	Id ChainId `json:"id"`

	// Default to current day
	Title       string `json:"title"`
	Description string `json:"description"`

	// Which day this chain is created
	Day time.Time `json:"day"`

	// Total number of habits - completed habits
	Progress int `json:"progress"`

	Habits []*Habit `json:"habits"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

func CreateChain(id ChainId, title, description string, day time.Time) (*Chain, error) {
	if title == "" {
		title = day.Format("02 01 2006")
	}

	now := time.Now().UTC()
	return &Chain{Id: id, Title: title, Description: description, Day: day, Progress: 0, Habits: nil, CreatedAt: now, UpdatedAt: nil}, nil
}

func (c *Chain) TotalNumberOfHabits() int {
	return len(c.Habits)
}

func (c *Chain) AddHabit(habit *Habit) {
	c.Habits = append(c.Habits, habit)
	now := time.Now().UTC()
	c.UpdatedAt = &now
}

func (h *Chain) ChangeHabitTitle(habitId HabitId, title string) error {
	for _, habit := range h.Habits {
		if habit.Id == habitId {
			if err := habit.ChangeTitle(title); err != nil {
				return err
			}
			return nil
		}
	}

	return ErrHabitNotFound
}

func (c *Chain) ToggleHabitCompletion(habitId HabitId) error {
	for _, habit := range c.Habits {
		if habit.Id == habitId {
			habit.ToggleCompletion()
			if habit.IsCompleted {
				c.Progress++
			} else {
				c.Progress--
			}
			return nil
		}
	}

	return ErrHabitNotFound
}

func (c *Chain) RemoveHabit(habitId HabitId) error {
	index := -1
	for i, habit := range c.Habits {
		if habit.Id == habitId {
			index = i
		}
	}

	if index < 0 {
		return ErrHabitNotFound
	}

	habit := c.Habits[index]
	c.Habits = append(c.Habits[:index], c.Habits[index+1:]...)
	if habit.IsCompleted {
		c.Progress--
	}
	now := time.Now().UTC()
	c.UpdatedAt = &now

	//for i, habit := range c.Habits {
	//	if habit.Id == habitId {
	//		c.Habits = append(c.Habits[:i], c.Habits[i+1:]...)
	//		if habit.IsCompleted {
	//			c.Progress--
	//		}
	//		now := time.Now().UTC()
	//		c.UpdatedAt = &now
	//		return nil
	//	}
	//}

	return nil
}

type ChainService interface {
	// All habits for heath map.
	HeatMap(ctx context.Context, year int) ([]*HeathMap, int, error)
	// Don't break the chain list
	MontlyChain(ctx context.Context, year, month int) ([]*DontBreakTheChain, int, error)
	Get(ctx context.Context, id int) (*Chain, error)
	GetByDay(ctx context.Context, day time.Time) (*Chain, error)
	Create(ctx context.Context, chain *Chain) (int, error)
	Delete(ctx context.Context, id int) error
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
