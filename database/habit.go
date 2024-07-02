package database

import (
	"context"
	"database/sql"
	"habheath"
	"time"
)

var (
	ErrHabitNotFound = habheath.Errorf(habheath.ENOTFOUND, "Habit not found.")
)

type HabitServiceImpl struct {
	db *DB
}

func NewHabitService(db *DB) habheath.HabitService {
	return &HabitServiceImpl{db: db}
}

// Compile-time check to ensure HabitServiceImpl implements ChainService
var _ habheath.HabitService = (*HabitServiceImpl)(nil)

func (s *HabitServiceImpl) HeatMap(ctx context.Context, year int) ([]*habheath.HeathMap, int, error) {
	// Implementation
	return nil, 0, nil
}

func (s *HabitServiceImpl) MontlyChain(ctx context.Context, year int) ([]*habheath.DontBreakTheChain, int, error) {
	// Implementation
	return nil, 0, nil
}

func (s *HabitServiceImpl) GetAllByDay(ctx context.Context, day time.Time) (*habheath.Chain, error) {
	const getHabitsQuery = `
		SELECT 
		    id,
		    title,
		    day,
			is_completed,
		    updated_at
		FROM habit
		ORDER BY id ASC
	`
	rows, err := s.db.db.QueryContext(ctx, getHabitsQuery)
	if err != nil {
		return nil, err
	}

	habits := make([]*habheath.Habit, 0)
	for rows.Next() {
		var h habheath.Habit
		//(*NullTime)(chain.UpdatedAt)
		var dayStr string
		var updatedAtStr string
		if err := rows.Scan(&h.Id, &h.Title, &dayStr, &h.IsCompleted, &updatedAtStr); err != nil {
			return nil, err
		}
		h.Day, _ = time.Parse(time.RFC3339, dayStr)
		h.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
		habits = append(habits, &h)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := &habheath.Chain{Title: day.Format(time.DateOnly), Habits: habits}
	return result, nil
}

func (s *HabitServiceImpl) Create(ctx context.Context, habit *habheath.Habit) error {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const createHabitQuery = `INSERT INTO habit (title, day, is_completed, updated_at)
			VALUES (?, ?, ?, ?)`

	ads := habit.Day.Format(time.RFC3339)
	result, err := tx.ExecContext(ctx, createHabitQuery, habit.Title, ads, habit.IsCompleted, habit.UpdatedAt)

	if err != nil {
		return FormatError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	habit.Id = habheath.HabitId(id)

	return tx.Commit()
}

func (s *HabitServiceImpl) Delete(ctx context.Context, id habheath.HabitId) error {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = checkHabitExists(ctx, tx, id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM habit WHERE id = ?`, id); err != nil {
		return FormatError(err)
	}

	return tx.Commit()
}

func (s *HabitServiceImpl) Update(ctx context.Context, habit *habheath.Habit) error {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = checkHabitExists(ctx, tx, habit.Id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE habit
		SET title = ?,
			is_completed = ?,
			updated_at = ?
		WHERE id = ?
	`,
		habit.Title,
		habit.IsCompleted,
		habit.UpdatedAt,
		habit.Id); err != nil {
		return FormatError(err)
	}

	return tx.Commit()
}

func checkHabitExists(ctx context.Context, tx *sql.Tx, id habheath.HabitId) error {
	var n int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM habit WHERE id = ?`, id).Scan(&n); err != nil {
		return FormatError(err)
	} else if n == 0 {
		return ErrHabitNotFound
	}

	return nil
}
