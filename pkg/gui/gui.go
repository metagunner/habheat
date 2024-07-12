package gui

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/metagunner/habheath/pkg/app"
	"github.com/metagunner/habheath/pkg/config"
	"github.com/metagunner/habheath/pkg/database"
	"github.com/metagunner/habheath/pkg/models"
	"github.com/metagunner/habheath/pkg/utils"
	"github.com/samber/lo"
)

type Gui struct {
	g                  *gocui.Gui
	db                 *database.DB
	ViewHeathmap       *gocui.View
	YearsSelectList    *SelectList
	ChainPanel         *ChainPanelContext
	HabitsPanel        *HabitPanelContext
	mustRenderHeathmap bool
	HabitService       models.HabitService
	heathmapFirstDate  time.Time
	heathmapLastDate   time.Time
	Config             *config.UserConfig
	StatusView         *gocui.View
	version            string
}

type HeathGrid struct {
	row                 int
	column              int
	key                 time.Time
	rank                int
	shade               string
	totalNumberOfHabits int
	completedHabits     int
	haveInfo            bool
}

var (
	cursorX             int
	cursorY             int
	grid                [][]*HeathGrid
	newVersionAvailable bool
)

func NewGui(config *config.UserConfig, db *database.DB, version string) *Gui {
	return &Gui{Config: config, db: db, version: version}
}

func (gui *Gui) initGocui() (*gocui.Gui, error) {
	g, err := gocui.NewGui(gocui.NewGuiOpts{
		OutputMode:      gocui.OutputTrue,
		SupportOverlaps: false,
	})
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (gui *Gui) Run() error {
	g, err := gui.initGocui()
	if err != nil {
		log.Println(err)
	}

	gui.g = g
	defer gui.g.Close()

	gui.g.FgColor = config.GetGocuiStyle(gui.Config.Gui.Theme.InactiveBorderColor)
	gui.g.SelFgColor = config.GetGocuiStyle(gui.Config.Gui.Theme.ActiveBorderColor)
	gui.g.FrameColor = config.GetGocuiStyle(gui.Config.Gui.Theme.InactiveBorderColor)
	gui.g.SelFrameColor = config.GetGocuiStyle(gui.Config.Gui.Theme.ActiveBorderColor)

	gui.g.SetManager(gocui.ManagerFunc(gui.layout))

	gui.HabitService = database.NewHabitService(gui.db)
	gui.initializeGrid()

	if err := gui.createAllViews(); err != nil {
		return err
	}

	gui.mustRenderHeathmap = true

	if _, err := gui.g.SetCurrentView("years"); err != nil {
		return err
	}

	gui.setKeybindings()

	newVersionAvailable = app.CheckForNewUpdate(gui.version)

	return gui.g.MainLoop()
}

func (gui *Gui) setKeybindings() error {
	var err error
	err = gui.g.SetKeybinding("", 'q', gocui.ModNone, quit)
	if err != nil {
		return err
	}

	err = gui.g.SetKeybinding("", '1', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gui.nextWindow("years")
	})
	if err != nil {
		return err
	}

	err = gui.g.SetKeybinding("", '2', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gui.nextWindow("heathmap")
	})
	if err != nil {
		return err
	}

	gui.g.SetKeybinding("heathmap", gocui.KeyArrowUp, gocui.ModNone, moveCursor(gui, -1, 0))
	gui.g.SetKeybinding("heathmap", gocui.KeyArrowDown, gocui.ModNone, moveCursor(gui, 1, 0))
	gui.g.SetKeybinding("heathmap", gocui.KeyArrowLeft, gocui.ModNone, moveCursor(gui, 0, -1))
	gui.g.SetKeybinding("heathmap", gocui.KeyArrowRight, gocui.ModNone, moveCursor(gui, 0, 1))
	gui.g.SetKeybinding("heathmap", 'k', gocui.ModNone, moveCursor(gui, -1, 0))
	gui.g.SetKeybinding("heathmap", 'j', gocui.ModNone, moveCursor(gui, 1, 0))
	gui.g.SetKeybinding("heathmap", 'h', gocui.ModNone, moveCursor(gui, 0, -1))
	gui.g.SetKeybinding("heathmap", 'l', gocui.ModNone, moveCursor(gui, 0, 1))
	gui.g.SetKeybinding("heathmap", 'a', gocui.ModNone, gui.wrappedHandler(gui.ChainPanel.OpenChainPanel))

	err = gui.g.SetKeybinding("years", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		selected := gui.YearsSelectList.GetSelected().option
		gui.reInitGrid(selected)
		gui.renderHeathmap()
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) nextWindow(viewName string) error {
	if _, err := gui.g.SetCurrentView(viewName); err != nil {
		return err
	}
	_, err := gui.g.SetViewOnTop(viewName)
	if err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// layout is called for every screen re-render e.g. when the screen is resized
func (gui *Gui) layout(g *gocui.Gui) error {
	g.Highlight = true

	gui.YearsSelectList.Render()

	gui.renderHeathmap()
	gui.g.SetViewOnTop("colors")
	gui.renderVersion()
	return nil
}

func (gui *Gui) renderVersion() {
	gui.StatusView.Clear()
	newVersionText := lo.Ternary(newVersionAvailable, "new version available!", "")
	fmt.Fprintf(gui.StatusView, "%s %s", gui.version, newVersionText)
}

func (gui *Gui) renderHeathmap() error {
	v := gui.ViewHeathmap
	v.Clear()

	months := utils.GetMonths(gui.heathmapLastDate)
	for _, month := range months {
		monthName := fmt.Sprintf("   %s   ", month.Format("Jan"))
		fmt.Fprint(v, monthName)
	}
	fmt.Fprintln(v)

	defaultTheme := gui.Config.Gui.Theme.Selected
	theme := gui.Config.Gui.Theme.ColorSchemes[defaultTheme]

	labels := []string{"   ", "Mon", "   ", "Wed", "   ", "Fri", "   "}
	printLabel := 0
	// Print the grid
	for _, row := range grid {
		fmt.Fprintf(v, "%s  ", labels[printLabel])
		for _, slot := range row {
			if slot.row == cursorY && slot.column == cursorX {
				fmt.Fprintf(v, "%s", theme.CursorValue)
			} else {
				fmt.Fprintf(v, "%s", slot.shade)
			}
		}
		printLabel++
		fmt.Fprintln(v)
	}

	fmt.Fprintln(v)
	info := grid[cursorY][cursorX]
	if !info.key.IsZero() {
		if info.haveInfo {
			fmt.Fprintf(v, "%d/%d habits on %s %s", info.completedHabits, info.totalNumberOfHabits, info.key.Format("Jan"), utils.GetOrdinalSuffix(info.key.Day()))
		} else {
			fmt.Fprintf(v, "No habits on %s %s", info.key.Format("Jan"), utils.GetOrdinalSuffix(info.key.Day()))
		}
	}

	return nil
}

func (gui *Gui) wrappedHandler(f func() error) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return f()
	}
}

func moveCursor(g *Gui, dy, dx int) func(*gocui.Gui, *gocui.View) error {
	return func(gui *gocui.Gui, v *gocui.View) error {
		// Calculate new cursor position
		newCursorX := cursorX + dx
		newCursorY := cursorY + dy

		// Ensure cursor stays within grid bounds
		if newCursorX < 0 {
			newCursorX = 0
		} else if newCursorX >= 53 {
			newCursorX = 52
		}
		if newCursorY < 0 {
			newCursorY = 0
		} else if newCursorY >= 7 {
			newCursorY = 6
		}

		// Update cursor position
		cursorX = newCursorX
		cursorY = newCursorY

		g.renderHeathmap()

		return nil
	}
}

// Init grid for the default view
func (gui *Gui) initializeGrid() {
	// Get today's date
	now := time.Now()
	today := utils.CreateDate(now.Year(), now.Month(), now.Day())

	if today.Weekday() == time.Sunday {
		today = today.AddDate(0, 0, 7)
	}

	// Find the next upcoming Sunday after today
	for today.Weekday() != time.Sunday {
		today = today.AddDate(0, 0, 1)
	}

	// Calculate the start date: 53 weeks before the next upcoming Sunday
	startDate := today.AddDate(0, 0, -53*7)

	// Find the first Sunday on or after the start date
	firstSunday := startDate
	for firstSunday.Weekday() != time.Sunday {
		firstSunday = firstSunday.AddDate(0, 0, 1)
	}
	startDate = firstSunday

	gui.heathmapFirstDate = today
	gui.heathmapLastDate = startDate

	// Create the grid
	grid = make([][]*HeathGrid, 7)
	for i := range grid {
		grid[i] = make([]*HeathGrid, 53)
	}

	defaultTheme := gui.Config.Gui.Theme.Selected
	theme := gui.Config.Gui.Theme.ColorSchemes[defaultTheme]
	heathmaps, _, err := gui.HabitService.HeatMap(context.Background(), startDate, today)
	if err != nil {
		panic(err)
	}

	// Fill the grid with colored boxes
	today = utils.CreateDate(now.Year(), now.Month(), now.Day())
	currentDate := startDate
	for col := 0; col < 53; col++ {
		for row := 0; row < 7; row++ {
			heathGrid := &HeathGrid{
				row:    row,
				column: col,
				key:    currentDate,
			}
			if currentDate.After(today) {
				// grid[row][col] = fmt.Sprintf("\x1b[48;5;52m%d \x1b[0m", currentDate.Day())
				heathGrid.shade = theme.InvalidDayValue
			} else {
				heathmap, ok := heathmaps[currentDate]
				if ok {
					shadeIndex := GetTheShade(heathmap)
					heathGrid.rank = shadeIndex
					colorCode, haveShade := theme.StatusValues[shadeIndex]
					if !haveShade {
						colorCode = theme.ZeroCompletedHabitValue
					}
					heathGrid.shade = colorCode
					heathGrid.haveInfo = true
					heathGrid.totalNumberOfHabits = heathmap.TotalNumberOfHabits
					heathGrid.completedHabits = heathmap.CompletedHabits
				} else {
					heathGrid.shade = theme.NoHabitsValue
				}
			}
			grid[row][col] = heathGrid
			currentDate = currentDate.AddDate(0, 0, 1) // move to the next day
		}
	}
}

// Init grid selected year
func (gui *Gui) initFromTo() {
	// from=2023-01-01&to=2023-12-31
	selectedYear, _ := strconv.Atoi(gui.YearsSelectList.GetSelected().option)
	from := utils.CreateDate(selectedYear, 1, 1)
	to := utils.CreateDate(selectedYear, 12, 31)

	gui.heathmapFirstDate = to
	gui.heathmapLastDate = from

	// Create the grid
	grid = make([][]*HeathGrid, 7)
	for i := range grid {
		grid[i] = make([]*HeathGrid, 53)
	}

	defaultTheme := gui.Config.Gui.Theme.Selected
	theme := gui.Config.Gui.Theme.ColorSchemes[defaultTheme]
	heathmaps, _, err := gui.HabitService.HeatMap(context.Background(), from, to)
	if err != nil {
		panic(err)
	}

	// Fill the grid with colored boxes
	startWeekday := int(from.Weekday())
	now := time.Now()
	today := utils.CreateDate(now.Year(), now.Month(), now.Day())
	currentDate := from
	for col := 0; col < 53; col++ {
		for row := 0; row < 7; row++ {
			heathGrid := &HeathGrid{
				row:    row,
				column: col,
				key:    currentDate,
			}
			if (col == 0 && startWeekday > row) || currentDate.Year() != from.Year() {
				heathGrid.shade = theme.InvalidDayValue // not belongs to year, the january starts from the wednesday
				heathGrid.key = time.Time{}
			} else {
				if currentDate.After(today) {
					heathGrid.shade = theme.InvalidDayValue // the now is 2024 july 8, after this time all of them are null
				} else {
					heathmap, ok := heathmaps[currentDate]
					if ok {
						shadeIndex := GetTheShade(heathmap)
						heathGrid.rank = shadeIndex
						colorCode, haveShade := theme.StatusValues[shadeIndex]
						if !haveShade {
							colorCode = theme.ZeroCompletedHabitValue
						}
						heathGrid.shade = colorCode
						heathGrid.haveInfo = true
						heathGrid.totalNumberOfHabits = heathmap.TotalNumberOfHabits
						heathGrid.completedHabits = heathmap.CompletedHabits
					} else {
						heathGrid.shade = theme.NoHabitsValue
					}
				}
				currentDate = currentDate.AddDate(0, 0, 1) // move to the next day
			}
			grid[row][col] = heathGrid
		}
	}
}

func (gui *Gui) reInitGrid(selected string) error {
	if selected == "Default" {
		gui.initializeGrid()
	} else {
		_, err := strconv.Atoi(selected)
		if err != nil {
			return err
		}
		gui.initFromTo()
	}

	return nil
}

func GetTheShade(heathmap *models.HeathMap) int {
	if heathmap.CompletedHabits == 0 {
		return 0
	}

	completionRatio := float64(heathmap.CompletedHabits) / float64(heathmap.TotalNumberOfHabits)
	shade := int(math.Round(completionRatio * 5))

	if shade < 1 {
		shade = 1
	}
	return shade
}

func (gui *Gui) GetDateFromHeathmapCursor() time.Time {
	return grid[cursorY][cursorX].key
}
