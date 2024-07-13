package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/metagunner/habheat/pkg/config"
)

type SelectList struct {
	gui               *Gui
	view              *gocui.View
	items             []SelectItem
	selectedIndex     int
	cursorPos         int
	getDisplayStrings func() []SelectItem
	emptyMessage      string
	isRendered        bool
}

type SelectItem struct {
	id     int
	option string
}

func NewSelectList(g *Gui, view *gocui.View, getDisplayStrings func() []SelectItem) *SelectList {
	s := &SelectList{gui: g, view: view, getDisplayStrings: getDisplayStrings}

	// handlers
	if err := s.gui.g.SetKeybinding(s.view.Name(), config.GetKey(g.Config.Keybinding.Universal.NextItem), gocui.ModNone, g.wrappedHandler(s.HandleNextLine)); err != nil {
		panic(err)
	}
	if err := s.gui.g.SetKeybinding(s.view.Name(), config.GetKey(g.Config.Keybinding.Universal.PrevItem), gocui.ModNone, g.wrappedHandler(s.HandlePrevLine)); err != nil {
		panic(err)
	}
	if err := s.gui.g.SetKeybinding(s.view.Name(), config.GetKey(g.Config.Keybinding.Universal.NextItemAlt), gocui.ModNone, g.wrappedHandler(s.HandleNextLine)); err != nil {
		panic(err)
	}
	if err := s.gui.g.SetKeybinding(s.view.Name(), config.GetKey(g.Config.Keybinding.Universal.PrevItemAlt), gocui.ModNone, g.wrappedHandler(s.HandlePrevLine)); err != nil {
		panic(err)
	}

	return s
}

func (self *SelectList) HandlePrevLine() error {
	maxOpt := len(self.items)
	if maxOpt == 0 {
		return nil
	}

	next := self.selectedIndex - 1
	if next < 0 {
		next = 0
	}

	self.selectedIndex = next
	self.ScrollUp()
	self.Render()

	return nil
}

func (self *SelectList) HandleNextLine() error {
	maxOpt := len(self.items)
	if maxOpt == 0 {
		return nil
	}

	next := self.selectedIndex + 1
	if next >= maxOpt {
		next = self.selectedIndex
	}

	self.selectedIndex = next
	self.ScrollDown()
	self.Render()

	return nil
}

func (self *SelectList) GetSelected() SelectItem {
	if len(self.items) == 0 {
		return SelectItem{}
	}
	return self.items[self.selectedIndex]
}

func (self *SelectList) SetEmptyMessage(message string) {
	self.emptyMessage = message
}

func (self *SelectList) RefreshOptions() {
	items := self.getDisplayStrings()
	self.items = items
}

func (self *SelectList) Render() {
	if !self.isRendered {
		self.RefreshOptions()
		self.isRendered = true
	}
	if len(self.items) == 0 && self.emptyMessage != "" {
		fmt.Fprintln(self.view, self.emptyMessage)
		return
	}

	self.view.Clear()
	for _, item := range self.items {
		self.view.WriteString(item.option + "\n")
	}
	self.view.SetCursor(0, self.cursorPos)
}

func (self *SelectList) ScrollUp() {
	viewPortStart, viewPortHeight := self.ViewPortYBounds()

	linesToScroll := 0
	if viewPortStart > 0 && self.selectedIndex <= viewPortStart+viewPortHeight {
		marginEnd := viewPortStart
		if self.selectedIndex < marginEnd {
			linesToScroll = 1
		}
	}

	if linesToScroll != 0 {
		self.view.ScrollUp(linesToScroll)

		cp := self.selectedIndex - viewPortStart
		if self.selectedIndex < viewPortStart {
			cp = 0
		}
		self.cursorPos = cp
	} else {
		// set cursor to upper item
		if self.cursorPos != 0 {
			self.cursorPos = self.view.CursorY() - 1
		}
	}
}

func (self *SelectList) ScrollDown() {
	viewPortStart, viewPortHeight := self.ViewPortYBounds()

	linesToScroll := 0
	if len(self.items) > viewPortStart+viewPortHeight {
		marginStart := viewPortStart + viewPortHeight
		if self.selectedIndex > marginStart {
			linesToScroll = 1
		}
	}

	if linesToScroll != 0 {
		self.view.ScrollDown(linesToScroll)
		self.cursorPos = viewPortHeight
	} else {
		if len(self.items)-1 != self.cursorPos {
			self.cursorPos = self.view.CursorY() + 1
		}
	}
}

// tells us the start of line indexes shown in the view currently as well as the capacity of lines shown in the viewport.
func (self *SelectList) ViewPortYBounds() (int, int) {
	_, start := self.view.Origin()
	length := self.view.InnerHeight()
	return start, length
}
