package habheath

import (
	"context"
	"time"
)

var ErrAtleastOneHabitIsRequired = Errorf(EINVALID, "At least one habit is required.")

func ErrHabitNotFound(habitId int) *Error {
	return Errorf(ENOTFOUND, "%q %d %q", "Habit with the given id:", habitId, " not found.")
}

type Chain struct {
	Id int `json:"id"`

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

func CreateChain(id int, title, description string, day time.Time, habits []*Habit) (*Chain, error) {
	if len(habits) == 0 {
		return nil, ErrAtleastOneHabitIsRequired
	}

	now := time.Now().UTC()
	if title == "" {
		title = now.Format("02 01 2006")
	}

	var progress int
	for i := 0; i < len(habits); i++ {
		if habits[i].IsCompleted {
			progress++
		}
	}

	return &Chain{id, title, description, day, progress, habits, now, nil}, nil
}

func (c *Chain) TotalNumberOfHabits() int {
	return len(c.Habits)
}

func (c *Chain) AddHabit(habit *Habit) {
	c.Habits = append(c.Habits, habit)
	now := time.Now().UTC()
	c.UpdatedAt = &now
}

func (h *Chain) ChangeHabitTitle(habitId int, title string) error {
	for _, habit := range h.Habits {
		if habit.Id == habitId {
			if err := habit.ChangeTitle(title); err != nil {
				return err
			}
			return nil
		}
	}

	return ErrHabitNotFound(habitId)
}

func (c *Chain) ToggleHabitCompletion(habitId int) error {
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

	return ErrHabitNotFound(habitId)
}

func (c *Chain) RemoveHabit(habitId int) error {
	if len(c.Habits) == 1 {
		return ErrAtleastOneHabitIsRequired
	}

	index := -1
	for i, habit := range c.Habits {
		if habit.Id == habitId {
			index = i
		}
	}

	if index < 0 {
		return ErrHabitNotFound(habitId)
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
