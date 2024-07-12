package gui

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/metagunner/habheath/pkg/utils"
	"github.com/samber/lo"
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
	yearsV.FgColor = gocui.ColorDefault
	yearsV.SelBgColor = gocui.ColorBlue
	yearsV.InactiveViewSelBgColor = gocui.ColorDefault | gocui.AttrBold

	getDisplayStrings := func() []SelectItem {
		currentYear := time.Now().Year()
		var years []int
		var dayStr string
		var ts time.Time
		err := gui.db.QueryRow(context.Background(), `SELECT day FROM habit ORDER BY day LIMIT 1`).Scan(&dayStr)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return []SelectItem{{id: 0, option: strconv.Itoa(currentYear)}}
			}
		}

		ts, _ = time.Parse(time.RFC3339, dayStr)
		years = utils.GetYearsBetween(ts.Year(), currentYear)

		return append([]SelectItem{{id: 0, option: "Default"}}, lo.Map(years, func(year int, _ int) SelectItem {
			return SelectItem{id: 0, option: strconv.Itoa(year)}
		})...)
	}
	gui.YearsSelectList = NewSelectList(gui, yearsV, getDisplayStrings)
	gui.YearsSelectList.view.Highlight = true

	heathmapV, err := gui.g.SetView("heathmap", 11, 0, maxX-1, maxY-1, 0)
	if err != nil && !gocui.IsUnknownView(err) {
		return err
	}
	heathmapV.Title = "Habheath"
	heathmapV.FrameRunes = roundedFrameRunes
	heathmapV.TitlePrefix = "2"

	gui.ViewHeathmap = heathmapV

	colorsV, err := gui.g.SetView("colors", maxX-23, 0, maxX-2, 2, 0)
	if err != nil && !gocui.IsUnknownView(err) {
		return err
	}
	colorsV.Title = "Colors"
	colorsV.FrameRunes = roundedFrameRunes
	defaultTheme := gui.Config.Gui.Theme.Selected
	color := gui.Config.Gui.Theme.ColorSchemes[defaultTheme]
	fmt.Fprint(colorsV, "Less ")
	for i := 1; i <= len(color.StatusValues); i++ {
		fmt.Fprint(colorsV, color.StatusValues[i])
	}
	fmt.Fprint(colorsV, " More")

	habitPanel, err := gui.g.SetView("habitpanel", maxX/2-30, maxY/2-2, maxX/2+30, maxY/2, 0)
	if err != nil && !gocui.IsUnknownView(err) {
		return err
	}
	gui.HabitsPanel = NewHabitPanelContext(habitPanel, gui)
	habitPanel.Visible = false
	habitPanel.Editable = true
	habitPanel.Highlight = true

	chainPanel, err := gui.g.SetView("chainpanel", maxX/4, maxY/4, 3*maxX/4, 3*maxY/4, 0)
	if err != nil && !gocui.IsUnknownView(err) {
		return err
	}
	chainPanel.Title = "Habits"
	chainPanel.FgColor = gocui.ColorWhite
	chainPanel.SelBgColor = gocui.ColorBlue
	chainPanel.InactiveViewSelBgColor = gocui.ColorDefault | gocui.AttrBold
	gui.ChainPanel = NewChainPanelContext(chainPanel, gui, gui.HabitService)
	chainPanel.Visible = false
	chainPanel.CanScrollPastBottom = true
	chainPanel.Highlight = true

	return nil
}
