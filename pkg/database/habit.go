package database

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/metagunner/habheath/pkg/app"
	"github.com/metagunner/habheath/pkg/models"
	"github.com/metagunner/habheath/pkg/utils"
)

var ErrHabitNotFound = app.Errorf(app.ENOTFOUND, "Habit not found.")

type HabitServiceImpl struct {
	db *DB
}

func NewHabitService(db *DB) models.HabitService {
	return &HabitServiceImpl{db: db}
}

// Compile-time check to ensure HabitServiceImpl implements ChainService
var _ models.HabitService = (*HabitServiceImpl)(nil)

func (s *HabitServiceImpl) HeatMap(ctx context.Context, from time.Time, to time.Time) (map[time.Time]*models.HeathMap, int, error) {
	const getHeathMapQuery = `
		SELECT 
			COUNT(*) AS total_number_of_habits,
			SUM(CASE WHEN is_completed = 1 THEN 1 ELSE 0 END) AS completed_habits,
		    strftime('%d', day) AS day,
		    strftime('%m', day) AS month,
		    strftime('%Y', day) AS year
		FROM habit
		WHERE day >= ? 
			AND day <= ?
		GROUP BY day, month, year
		ORDER BY year, month, day
	`

	fromQuery := from.UTC().Format(time.RFC3339)
	toQuery := to.UTC().Format(time.RFC3339)
	rows, err := s.db.db.QueryContext(ctx, getHeathMapQuery, fromQuery, toQuery)
	if err != nil {
		return nil, 0, err
	}

	result := make(map[time.Time]*models.HeathMap, 0)
	for rows.Next() {
		var h models.HeathMap
		var day string
		var month string
		var year string
		if err := rows.Scan(&h.TotalNumberOfHabits, &h.CompletedHabits, &day, &month, &year); err != nil {
			return nil, 0, err
		}
		h.Day, _ = strconv.Atoi(day)
		h.Month, _ = strconv.Atoi(month)
		h.Year, _ = strconv.Atoi(year)
		time := utils.CreateDate(h.Year, time.Month(h.Month), h.Day)
		result[time] = &h
	}

	if err := rows.Close(); err != nil {
		return nil, 0, err
	}

	return result, len(result), nil
}

func (s *HabitServiceImpl) GetAllByDay(ctx context.Context, day time.Time) (*models.Chain, error) {
	const getHabitsQuery = `
		SELECT 
		    id,
		    title,
		    day,
			is_completed,
		    updated_at
		FROM habit
		WHERE day = ?
		ORDER BY id ASC
	`

	tsQuery := day.UTC().Format(time.RFC3339)
	rows, err := s.db.db.QueryContext(ctx, getHabitsQuery, tsQuery)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	habits := make([]*models.Habit, 0)
	for rows.Next() {
		var h models.Habit
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

	result := &models.Chain{Title: day.Format(time.DateOnly), Habits: habits}
	return result, nil
}

func (s *HabitServiceImpl) Create(ctx context.Context, habit *models.Habit) error {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const createHabitQuery = `INSERT INTO habit (title, day, is_completed, updated_at) VALUES (?, ?, ?, ?)`

	habitDay := habit.Day.Format(time.RFC3339)
	result, err := tx.ExecContext(ctx, createHabitQuery, habit.Title, habitDay, habit.IsCompleted, habit.UpdatedAt)
	if err != nil {
		return FormatError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	habit.Id = models.HabitId(id)

	return tx.Commit()
}

func (s *HabitServiceImpl) Delete(ctx context.Context, id models.HabitId) error {
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

func (s *HabitServiceImpl) Update(ctx context.Context, habit *models.Habit) error {
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

func checkHabitExists(ctx context.Context, tx *sql.Tx, id models.HabitId) error {
	var n int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM habit WHERE id = ?`, id).Scan(&n); err != nil {
		return FormatError(err)
	} else if n == 0 {
		return ErrHabitNotFound
	}

	return nil
}
