package habheath_test

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"habheath"
)

func TestCreateChain(t *testing.T) {
	t.Run("Given valid data creating Chain should succeed", func(t *testing.T) {
		chainId := 1
		day := time.Now()
		chain, err := habheath.CreateChain(chainId, "Test Chain", "Description", day, []*habheath.Habit{createTestHabit(1, chainId)})
		assertNotError(t, err)

		if chain.Id != 1 {
			t.Errorf("expected id 1, got %d", chain.Id)
		}
		if chain.Title != "Test Chain" {
			t.Errorf("expected title 'Test Chain', got %s", chain.Title)
		}
		if chain.Progress != 0 {
			t.Errorf("expected progress 0, got %d", chain.Progress)
		}
		if len(chain.Habits) != 1 {
			t.Errorf("expected 1 habit, got %d", len(chain.Habits))
		}
	})
	t.Run("Given chain without habit should fail", func(t *testing.T) {
		day := time.Now()
		_, err := habheath.CreateChain(1, "Test Chain", "Description", day, nil)
		assertError(t, err, habheath.ErrAtleastOneHabitIsRequired)
	})
}

func TestAddHabit(t *testing.T) {
	//	addHabitTests := []struct {
	//		name string
	//	}{
	//		{name: ""},
	//		{name: ""},
	//	}
	//
	//	for _, tt := range addHabitTests {
	//		got := tt.habit
	//		t.Run(tt.name, func(t *testing.T) {
	//			if got != tt.habit {
	//				t.Errorf("%#v got %g want %g", tt.shape, got, tt.hasArea)
	//			}
	//		})
	//	}
	t.Run("Given habit without title should fail", func(t *testing.T) {
		chain := createTestChainWithSingleHabit(1)
		habitId := len(chain.Habits) + 1
		_, err := habheath.CreateHabit(habitId, chain.Id, "", false)
		assertError(t, err, habheath.ErrHabitTitleIsRequired)
	})

	t.Run("Given valid habit should succeed", func(t *testing.T) {
		chain := createTestChainWithSingleHabit(1)
		habitId := len(chain.Habits) + 1
		habitTitle := "Habit " + strconv.Itoa(habitId)
		habit, err := habheath.CreateHabit(habitId, chain.Id, habitTitle, false)
		assertNotError(t, err)
		chain.AddHabit(habit)
		if len(chain.Habits) != 2 {
			t.Errorf("expected 2 habits, got %d", len(chain.Habits))
		}
	})
}

func TestChangeHabitTitle(t *testing.T) {
	t.Run("Given invalid habit title should fail", func(t *testing.T) {
		chain := createTestChainWithSingleHabit(1)
		err := chain.ChangeHabitTitle(1, "")
		assertError(t, err, habheath.ErrInvalidHabitTitle)
	})

	t.Run("Given valid habit title should succeed", func(t *testing.T) {
		chain := createTestChainWithSingleHabit(1)
		err := chain.ChangeHabitTitle(1, "Updated Habit 1")
		assertNotError(t, err)
		if chain.Habits[0].Title != "Updated Habit 1" {
			t.Errorf("expected 'Updated Habit 1', got %s", chain.Habits[0].Title)
		}
	})
}

func TestToggleHabitCompletion(t *testing.T) {
	chain := createTestChainWithSingleHabit(1)
	err := chain.ToggleHabitCompletion(1)
	assertNotError(t, err)

	if !chain.Habits[0].IsCompleted {
		t.Errorf("expected habit to be completed, got %v", chain.Habits[0].IsCompleted)
	}
	if chain.Progress != 1 {
		t.Errorf("expected progress 1, got %d", chain.Progress)
	}

	err = chain.ToggleHabitCompletion(1)
	assertNotError(t, err)

	if chain.Habits[0].IsCompleted {
		t.Errorf("expected habit to be not completed, got %v", chain.Habits[0].IsCompleted)
	}
	if chain.Progress != 0 {
		t.Errorf("expected progress 0, got %d", chain.Progress)
	}
}

func TestRemoveHabit(t *testing.T) {
	t.Run("Given chain with single habit removing habit should fail", func(t *testing.T) {
		chain := createTestChainWithSingleHabit(1)
		err := chain.RemoveHabit(1)
		assertError(t, err, habheath.ErrAtleastOneHabitIsRequired)
	})
	t.Run("Given not existing habit removing habit should fail", func(t *testing.T) {
		notExistHabitId := 99
		chain := createTestChainWithSingleHabit(1)
		chain.AddHabit(createTestHabit(2, chain.Id))
		err := chain.RemoveHabit(notExistHabitId)
		assertError(t, err, habheath.ErrHabitNotFound(notExistHabitId))
	})
	t.Run("Given chain with multiple habits removing habit should succeed", func(t *testing.T) {
		chain := createTestChainWithSingleHabit(1)
		chain.AddHabit(createTestHabit(2, chain.Id))

		err := chain.RemoveHabit(1)
		assertNotError(t, err)
		if len(chain.Habits) != 1 {
			t.Errorf("expected 1 habit, got %d", len(chain.Habits))
		}
		if chain.Habits[0].Id != 2 {
			t.Errorf("expected habit with id 2, got %d", chain.Habits[0].Id)
		}
	})
}

func TestChainProgress(t *testing.T) {
	chain := createTestChainWithSingleHabit(1)
	chain.AddHabit(createTestHabit(2, chain.Id))
	chain.AddHabit(createTestHabit(3, chain.Id))
	chain.AddHabit(createTestHabit(4, chain.Id))
	chain.ToggleHabitCompletion(3)
	chain.ToggleHabitCompletion(4)
	expectedProgress := 2 // Since there are 2 completed habits
	if chain.Progress != expectedProgress {
		t.Errorf("expected progress %d, got %d", expectedProgress, chain.Progress)
	}
}

func createTestChainWithSingleHabit(chainId int) *habheath.Chain {
	habitId := 1
	chain, _ := habheath.CreateChain(chainId, "Test Chain", "Description", time.Now().UTC(), []*habheath.Habit{createTestHabit(habitId, chainId)})
	return chain
}

func createTestHabit(id int, chainId int) *habheath.Habit {
	title := "Habit " + strconv.Itoa(id)
	habit, _ := habheath.CreateHabit(1, chainId, title, false)
	return habit
}

// Should not produce error
func assertNotError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but did not want one, got %v", got)
	}
}

// Should produce error
func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("wanted an error but did not get one")
	}

	if !errors.Is(got, want) {
		t.Errorf("got %q want %q", got, want)
	}
}
