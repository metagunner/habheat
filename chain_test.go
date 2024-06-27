package habheath_test

import (
	"strconv"
	"testing"
	"time"

	"habheath"
)

var now = time.Now().UTC()

func TestCreateChain(t *testing.T) {
	chain, err := habheath.CreateChain(1, "", "Description", now)
	assertNotError(t, err)
	want := now.Format("02 01 2006")
	if chain.Title != want {
		t.Errorf("want %q, got %q", want, chain.Title)
	}
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
	t.Run("Given empty title add habit should fail", func(t *testing.T) {
		chain := createTestChain(1)
		_, err := habheath.CreateHabit(1, chain.Id, "", false)
		assertError(t, err, habheath.ErrHabitTitleIsRequired)
	})

	t.Run("Given valid data add habit should succeed", func(t *testing.T) {
		chain := createTestChain(1)
		habit, err := habheath.CreateHabit(1, chain.Id, "Habit 1", false)
		assertNotError(t, err)

		chain.AddHabit(habit)
		if len(chain.Habits) != 1 {
			t.Errorf("expected 1 habits, got %d", len(chain.Habits))
		}
	})
}

func TestChangeHabitTitle(t *testing.T) {
	t.Run("Given empty title change habit title should fail", func(t *testing.T) {
		chain := createTestChain(1)
		habit := createTestHabit(1, chain.Id)
		chain.AddHabit(habit)

		err := chain.ChangeHabitTitle(1, "")
		assertError(t, err, habheath.ErrInvalidHabitTitle)
	})

	t.Run("Given valid title change habit title should succeed", func(t *testing.T) {
		chain := createTestChain(1)
		habit := createTestHabit(1, chain.Id)
		chain.AddHabit(habit)

		err := chain.ChangeHabitTitle(1, "Updated Habit 1")
		assertNotError(t, err)
		if chain.Habits[0].Title != "Updated Habit 1" {
			t.Errorf("expected habit title to be %v , got %s", "Updated Habit 1", chain.Habits[0].Title)
		}
	})
}

func TestToggleHabitCompletion(t *testing.T) {
	chain := createTestChain(1)
	habit := createTestHabit(1, chain.Id)
	chain.AddHabit(habit)

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
	t.Run("Given not existing habit remove habit should fail", func(t *testing.T) {
		notExistHabitId := habheath.HabitId(99)
		chain := createTestChainWithHabit()
		err := chain.RemoveHabit(notExistHabitId)
		assertError(t, err, habheath.ErrHabitNotFound)
	})
	t.Run("Given chain with multiple habits remove habit should succeed", func(t *testing.T) {
		chain := createTestChain(1)
		chain.AddHabit(createTestHabit(1, chain.Id))
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
	chain := createTestChain(1)
	chain.AddHabit(createTestHabit(1, chain.Id))
	chain.AddHabit(createTestHabit(2, chain.Id))
	chain.AddHabit(createTestHabit(3, chain.Id))
	chain.ToggleHabitCompletion(1)
	chain.ToggleHabitCompletion(2)
	expectedProgress := 2 // Since there are 2 completed habits
	if chain.Progress != expectedProgress {
		t.Errorf("expected progress %d, got %d", expectedProgress, chain.Progress)
	}
}

func createTestChain(chainId habheath.ChainId) *habheath.Chain {
	chain, _ := habheath.CreateChain(chainId, "Test Chain", "Description", now)
	return chain
}

func createTestHabit(id habheath.HabitId, chainId habheath.ChainId) *habheath.Habit {
	title := "Habit " + strconv.Itoa(int(id))
	habit, _ := habheath.CreateHabit(id, chainId, title, false)
	return habit
}

func createTestChainWithHabit() *habheath.Chain {
	chain := createTestChain(1)
	habit := createTestHabit(1, chain.Id)
	chain.AddHabit(habit)
	return chain
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

	if got.Error() != want.Error() {
		t.Errorf("got %q want %q", got, want)
	}
}
