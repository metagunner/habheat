package gui

import (
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) createAllViews() error {
	maxX, maxY := gui.g.Size()
	roundedFrameRunes := []rune{'─', '│', '╭', '╮', '╰', '╯'}

	yearsV, err := gui.g.SetView("years", 0, 0, 10, maxY-1, 0)
	if err != nil && !gocui.IsUnknownView(err) {
		return err
	}
	yearsV.Title = "Years"
	yearsV.FrameRunes = roundedFrameRunes
	yearsV.TitlePrefix = "1"

	gui.ViewYears = yearsV

	heathmapV, err := gui.g.SetView("heathmap", 11, 0, maxX-1, maxY-1, 0)
	if err != nil && !gocui.IsUnknownView(err) {
		return err
	}
	heathmapV.Title = "Habheath"
	heathmapV.Wrap = true
	heathmapV.FrameColor = gocui.ColorWhite
	heathmapV.SelBgColor = gocui.ColorBlue
	heathmapV.FrameRunes = roundedFrameRunes
	heathmapV.TitlePrefix = "2"

	gui.ViewHeathmap = heathmapV

	return nil
}
