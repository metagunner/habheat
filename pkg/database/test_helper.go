package database

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/metagunner/habheath/pkg/models"
	"github.com/metagunner/habheath/pkg/utils"
)

func SetupTestDB() (*DB, error) {
	dsn := "file::memory:?cache=shared"
	db := NewDB(dsn)
	if err := db.Open(); err != nil {
		return nil, err
	}

	return db, nil
}

// Seed the database with random dated habits for 12 month. Given year month are guaranteed to be seeded
func SeedTestData(ctx context.Context, db *DB, year int, month time.Month) {
	var habitCount int
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM habit`).Scan(&habitCount); err != nil {
		panic(err)
	}
	if habitCount == 0 {
		log.Println("Seeding the database")
		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			panic(err)
		}
		defer tx.Rollback()

		var b strings.Builder
		b.WriteString("INSERT INTO habit (title, day, is_completed, updated_at) VALUES")

		maxHabitsPerDay := 6
		var values []interface{}
		date := utils.CreateDate(year, month, 1)
		months := utils.GetMonths(date)

		for _, m := range months {
			for i := 1; i <= utils.GetDaysInMonth(m); i++ {
				title, _ := models.CreateHabitTitle(fmt.Sprintf("Habit %s %d", m.Format(time.DateOnly), i))
				day := utils.CreateDate(m.Year(), m.Month(), i)
				habitsPerDay := rand.Intn(maxHabitsPerDay)
				if m.Year() == year && m.Month() == month {
					habitsPerDay += 1
				}
				for j := 0; j < habitsPerDay; j++ {
					isCompleted := rand.Intn(2) == 1
					habit, _ := models.CreateHabit(title, day, isCompleted)

					if len(values) > 0 {
						b.WriteString(", ")
					}

					b.WriteString("(?, ?, ?, ?)")
					values = append(values, habit.Title, habit.Day.Format(time.RFC3339), habit.IsCompleted, habit.UpdatedAt.Format(time.RFC3339))
				}
			}
		}

		bulkInsert := b.String() + ";"
		result, err := db.db.ExecContext(ctx, bulkInsert, values...)
		if err != nil {
			panic(err)
		}
		rows, _ := result.RowsAffected()
		log.Printf("inserted into db: %d", rows)

		if err := tx.Commit(); err != nil {
			panic(err)
		}
	}
}
