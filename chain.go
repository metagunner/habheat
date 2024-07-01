package habheath

//
//import (
//	"context"
//	"database/sql"
//	"time"
//)
//
//var (
//	ErrHabitNotFound     = Errorf(ENOTFOUND, "Habit not found.")
//	ErrInvalidChainTitle = Errorf(EINVALID, "Invalid chain title.")
//)
//
////func (c *Chain) Progress(habit *Habit) int {
////	total := len(c.Habits)
////	completed := 0
////	for _, habit := range c.Habits {
////		if habit.IsCompleted {
////			completed++
////		}
////	}
////
////	return total - completed
////}
//
//func (c *Chain) AddHabit(habit *Habit) {
//	c.Habits = append(c.Habits, habit)
//	c.UpdatedAt = time.Now().UTC()
//}
//
//func (h *Chain) ChangeHabitTitle(habitId HabitId, title string) error {
//	for _, habit := range h.Habits {
//		if habit.Id == habitId {
//			if err := habit.ChangeTitle(title); err != nil {
//				return err
//			}
//			return nil
//		}
//	}
//
//	return ErrHabitNotFound
//}
//
//func (c *Chain) ToggleHabitCompletion(habitId HabitId) error {
//	for _, habit := range c.Habits {
//		if habit.Id == habitId {
//			habit.ToggleCompletion()
//			if habit.IsCompleted {
//				c.Progress++
//			} else {
//				c.Progress--
//			}
//			return nil
//		}
//	}
//
//	return ErrHabitNotFound
//}
//
//func (c *Chain) RemoveHabit(habitId HabitId) error {
//	//	for i, h := range c.Habits {
//	//		if h.Id == habitId {
//	//			c.Habits[i] = c.Habits[len(c.Habits)-1]
//	//			c.Habits[len(c.Habits)-1] = nil
//	//			c.Habits = c.Habits[:len(c.Habits)-1]
//	//			c.UpdatedAt = &now
//	//			return nil
//	//		}
//	//	}
//
//	index := -1
//	for i, habit := range c.Habits {
//		if habit.Id == habitId {
//			index = i
//		}
//	}
//
//	if index < 0 {
//		return ErrHabitNotFound
//	}
//
//	habit := c.Habits[index]
//	c.Habits = append(c.Habits[:index], c.Habits[index+1:]...)
//	if habit.IsCompleted {
//		c.Progress--
//	}
//	c.UpdatedAt = time.Now().UTC()
//
//	return nil
//}
