package habheath

import "time"

var ErrInvalidHabitTitle = Errorf(EINVALID, "Invalid habit title.")
var ErrHabitTitleIsRequired = Errorf(EINVALID, "Title is required.")

type HabitId int

type Habit struct {
	Id          HabitId   `json:"id"`
	ChainId     ChainId   `json:"chainid"`
	Title       string    `json:"title"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsCompleted bool      `json:"isCompleted"`
	// TODO: Order       int       `json:"order"`
}

func CreateHabit(id HabitId, chainId ChainId, title string, isCompleted bool) (*Habit, error) {
	if title == "" {
		return nil, ErrHabitTitleIsRequired
	}

	now := time.Now().UTC()
	updatedAt := now
	habit := &Habit{id, chainId, "", updatedAt, false}
	habit.ChangeTitle(title)
	return habit, nil
}

func (h *Habit) ToggleCompletion() {
	h.IsCompleted = !h.IsCompleted
	h.UpdatedAt = time.Now().UTC()
}

func (h *Habit) ChangeTitle(title string) error {
	if title == "" || len(title) > 250 {
		return ErrInvalidHabitTitle
	}

	if h.Title == title {
		return nil
	}

	h.Title = title
	h.UpdatedAt = time.Now().UTC()
	return nil
}
