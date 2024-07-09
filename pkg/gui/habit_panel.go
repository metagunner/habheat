package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
)

type HabitPanelContext struct {
	chainPanelContext *ChainPanelContext
	view              *gocui.View
	viewModel         *HabitPanelViewModel
	gui               *Gui
}

type HabitPanelViewModel struct {
	id        int
	title     string
	onConfirm func(string) error
}

func NewHabitPanelContext(v *gocui.View, gui *Gui) *HabitPanelContext {
	viewModel := &HabitPanelViewModel{}
	habitPanelContext := &HabitPanelContext{
		view:      v,
		viewModel: viewModel,
		gui:       gui,
	}
	gui.g.SetKeybinding(v.Name(), gocui.KeyEnter, gocui.ModNone, gui.wrappedHandler(habitPanelContext.OnConfirm))
	gui.g.SetKeybinding(v.Name(), gocui.KeyEsc, gocui.ModNone, gui.wrappedHandler(habitPanelContext.CloseHabitPanel))
	return habitPanelContext
}

func (self *HabitPanelContext) OnConfirm() error {
	title := self.GetHabitTitle()
	return self.viewModel.onConfirm(title)
}

func (self *HabitPanelContext) SetPanelState(
	id int,
	title string,
	habitInputTitle string,
	onConfirm func(string) error,
) {
	self.viewModel.id = id
	self.viewModel.title = title
	self.viewModel.onConfirm = onConfirm

	self.gui.g.Cursor = true
	self.view.Title = habitInputTitle
	self.view.Visible = true
	self.view.ClearTextArea()
	self.view.TextArea.TypeString(title)
	self.view.RenderTextArea()
}

func (self *HabitPanelContext) GetHabitTitle() string {
	return strings.TrimSpace(self.view.TextArea.GetContent())
}

func (self *HabitPanelContext) CloseHabitPanel() error {
	self.view.Clear()
	self.view.Visible = false
	self.gui.g.Cursor = false
	if _, err := self.gui.g.SetCurrentView(self.gui.ChainPanel.view.Name()); err != nil {
		return err
	}
	self.gui.ChainPanel.viewModel.list.RefreshOptions()
	self.gui.ChainPanel.viewModel.list.Render()

	return nil
}
