package gui

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/metagunner/habheat/pkg/config"
	"github.com/metagunner/habheat/pkg/models"
	"github.com/metagunner/habheat/pkg/utils"
	"github.com/samber/lo"
)

type ChainPanelContext struct {
	viewModel    *ChainPanelViewModel
	view         *gocui.View
	habitService models.HabitService
	gui          *Gui
}

type ChainPanelViewModel struct {
	list        *SelectList
	selectedDay time.Time
}

func NewChainPanelContext(v *gocui.View, gui *Gui, habitService models.HabitService) *ChainPanelContext {
	viewModel := &ChainPanelViewModel{}
	getDisplayStrings := func() []SelectItem {
		date := viewModel.selectedDay
		date = utils.CreateDate(date.Year(), date.Month(), date.Day())
		habitChain, err := habitService.GetAllByDay(context.Background(), date)
		if err != nil {
			return []SelectItem{}
		}

		result := []SelectItem{}
		if len(habitChain.Habits) > 0 {
			for i, habit := range habitChain.Habits {
				status := " "
				if habit.IsCompleted {
					status = "X"
				}
				result = append(result, SelectItem{id: int(habit.Id), option: fmt.Sprintf("%d. [%s] %s", i+1, status, habit.Title)})
			}
		}
		return result
	}
	viewModel.list = NewSelectList(gui, v, getDisplayStrings)
	viewModel.list.SetEmptyMessage("No habits for this day. Create one by presing the " + gui.Config.Keybinding.Heatmap.CreateHabit)

	chainPanelContext := &ChainPanelContext{
		viewModel:    viewModel,
		view:         v,
		habitService: habitService,
		gui:          gui,
	}

	heatmapKeys := gui.Config.Keybinding.Heatmap
	gui.g.SetKeybinding(v.Name(), config.GetKey(heatmapKeys.ToggleHabit), gocui.ModNone, gui.wrappedHandler(chainPanelContext.ToggleHabitCompletion))
	gui.g.SetKeybinding(v.Name(), config.GetKey(heatmapKeys.EditHabit), gocui.ModNone, gui.wrappedHandler(chainPanelContext.UpdateHabit))
	gui.g.SetKeybinding(v.Name(), config.GetKey(heatmapKeys.CreateHabit), gocui.ModNone, gui.wrappedHandler(chainPanelContext.AddHabit))
	gui.g.SetKeybinding(v.Name(), config.GetKey(heatmapKeys.DeleteHabit), gocui.ModNone, gui.wrappedHandler(chainPanelContext.RemoveHabit))
	gui.g.SetKeybinding(v.Name(), config.GetKey(gui.Config.Keybinding.Universal.Close), gocui.ModNone, gui.wrappedHandler(chainPanelContext.CloseChainPanel))

	return chainPanelContext
}

func (self *ChainPanelContext) OpenChainPanel() error {
	selectedDate := self.gui.GetDateFromHeatmapCursor()
	if selectedDate.IsZero() {
		return nil
	}
	self.view.Subtitle = selectedDate.Format(time.DateOnly)
	self.view.Clear()
	self.view.Visible = true

	viewName := self.view.Name()
	self.viewModel.selectedDay = selectedDate
	self.viewModel.list.RefreshOptions()
	self.viewModel.list.Render()

	if _, err := self.gui.g.SetViewOnTop(viewName); err != nil {
		return err
	}
	if _, err := self.gui.g.SetCurrentView(viewName); err != nil {
		return err
	}

	return nil
}

func (self *ChainPanelContext) CloseChainPanel() error {
	self.view.Clear()
	self.view.Visible = false
	if _, err := self.gui.g.SetCurrentView(self.gui.ViewHeatmap.Name()); err != nil {
		return err
	}
	self.gui.YearsSelectList.RefreshOptions()
	self.gui.YearsSelectList.Render()
	selected := self.gui.YearsSelectList.GetSelected().option
	if err := self.gui.reInitGrid(selected); err != nil {
		return err
	}
	if err := self.gui.renderHeatmap(); err != nil {
		return err
	}

	return nil
}

func (self *ChainPanelContext) RemoveHabit() error {
	selected := self.viewModel.list.GetSelected()
	if selected == (SelectItem{}) {
		return nil
	}

	if err := self.habitService.Delete(context.Background(), models.HabitId(selected.id)); err != nil {
		return err
	}
	self.view.Clear()
	self.viewModel.list.RefreshOptions()
	self.viewModel.list.Render()
	return nil
}

func (self *ChainPanelContext) ToggleHabitCompletion() error {
	selected := self.viewModel.list.GetSelected()
	if selected == (SelectItem{}) {
		return nil
	}

	chain, err := self.habitService.GetAllByDay(context.Background(), self.viewModel.selectedDay)
	if err != nil {
		return err
	}
	habit, finded := lo.Find(chain.Habits, func(x *models.Habit) bool { return x.Id == models.HabitId(selected.id) })
	if !finded {
		return errors.New("not found")
	}
	habit.ToggleCompletion()
	if err := self.habitService.Update(context.Background(), habit); err != nil {
		return err
	}
	self.view.Clear()
	self.viewModel.list.RefreshOptions()
	self.viewModel.list.Render()
	return nil
}

func (self *ChainPanelContext) UpdateHabit() error {
	selected := self.viewModel.list.GetSelected()
	if selected == (SelectItem{}) {
		return nil
	}

	chain, err := self.habitService.GetAllByDay(context.Background(), self.viewModel.selectedDay)
	if err != nil {
		return err
	}
	habit, finded := lo.Find(chain.Habits, func(x *models.Habit) bool { return x.Id == models.HabitId(selected.id) })
	if !finded {
		return errors.New("not found")
	}

	onConfirm := func(newtitle string) error {
		if err := habit.ChangeTitle(newtitle); err != nil {
			return err
		}
		if err := self.habitService.Update(context.Background(), habit); err != nil {
			return err
		}
		self.gui.HabitsPanel.CloseHabitPanel()

		return nil
	}
	self.gui.HabitsPanel.SetPanelState(int(habit.Id), habit.Title.String(), fmt.Sprintf("Habit %d", habit.Id), onConfirm)
	viewName := self.gui.HabitsPanel.view.Name()
	if _, err := self.gui.g.SetViewOnTop(viewName); err != nil {
		return err
	}
	if _, err := self.gui.g.SetCurrentView(viewName); err != nil {
		return err
	}

	return nil
}

func (self *ChainPanelContext) AddHabit() error {

	onConfirm := func(newtitle string) error {
		title, err := models.CreateHabitTitle(newtitle)
		if err != nil {
			return err
		}
		habit, err := models.CreateHabit(title, self.viewModel.selectedDay, false)
		if err != nil {
			return err
		}
		if err := self.habitService.Create(context.Background(), habit); err != nil {
			return err
		}
		self.gui.HabitsPanel.CloseHabitPanel()
		return nil
	}
	self.gui.HabitsPanel.SetPanelState(0, "", "New Habit", onConfirm)
	viewName := self.gui.HabitsPanel.view.Name()
	if _, err := self.gui.g.SetViewOnTop(viewName); err != nil {
		return err
	}
	if _, err := self.gui.g.SetCurrentView(viewName); err != nil {
		return err
	}

	return nil
}
