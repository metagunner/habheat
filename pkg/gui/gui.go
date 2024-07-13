package gui

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/metagunner/habheat/pkg/app"
	"github.com/metagunner/habheat/pkg/config"
	"github.com/metagunner/habheat/pkg/database"
	"github.com/metagunner/habheat/pkg/models"
	"github.com/metagunner/habheat/pkg/utils"
	"github.com/samber/lo"
)

type Gui struct {
	g                 *gocui.Gui
	db                *database.DB
	ViewHeatmap       *gocui.View
	YearsSelectList   *SelectList
	ChainPanel        *ChainPanelContext
	HabitsPanel       *HabitPanelContext
	mustRenderHeatmap bool
	HabitService      models.HabitService
	heatmapFirstDate  time.Time
	heatmapLastDate   time.Time
	Config            *config.UserConfig
	StatusView        *gocui.View
	version           string
}

type HeatGrid struct {
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
	grid                [][]*HeatGrid
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

	gui.mustRenderHeatmap = true

	if _, err := gui.g.SetCurrentView("years"); err != nil {
		return err
	}

	gui.setKeybindings()

	newVersionAvailable = app.CheckForNewUpdate(gui.version)

	return gui.g.MainLoop()
}

func (gui *Gui) setKeybindings() error {
	var err error
	err = gui.g.SetKeybinding("", config.GetKey(gui.Config.Keybinding.Universal.Quit), gocui.ModNone, quit)
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
		return gui.nextWindow("heatmap")
	})
	if err != nil {
		return err
	}

	heatmapKeys := gui.Config.Keybinding.Heatmap
	gui.g.SetKeybinding("heatmap", config.GetKey(heatmapKeys.UpAlt), gocui.ModNone, moveCursor(gui, -1, 0))
	gui.g.SetKeybinding("heatmap", config.GetKey(heatmapKeys.DownAlt), gocui.ModNone, moveCursor(gui, 1, 0))
	gui.g.SetKeybinding("heatmap", config.GetKey(heatmapKeys.LeftAlt), gocui.ModNone, moveCursor(gui, 0, -1))
	gui.g.SetKeybinding("heatmap", config.GetKey(heatmapKeys.RightAlt), gocui.ModNone, moveCursor(gui, 0, 1))
	gui.g.SetKeybinding("heatmap", config.GetKey(heatmapKeys.Up), gocui.ModNone, moveCursor(gui, -1, 0))
	gui.g.SetKeybinding("heatmap", config.GetKey(heatmapKeys.Down), gocui.ModNone, moveCursor(gui, 1, 0))
	gui.g.SetKeybinding("heatmap", config.GetKey(heatmapKeys.Left), gocui.ModNone, moveCursor(gui, 0, -1))
	gui.g.SetKeybinding("heatmap", config.GetKey(heatmapKeys.Right), gocui.ModNone, moveCursor(gui, 0, 1))
	gui.g.SetKeybinding("heatmap", config.GetKey(gui.Config.Keybinding.Universal.Select), gocui.ModNone, gui.wrappedHandler(gui.ChainPanel.OpenChainPanel))

	err = gui.g.SetKeybinding("years", config.GetKey(gui.Config.Keybinding.Universal.Select), gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		selected := gui.YearsSelectList.GetSelected().option
		gui.reInitGrid(selected)
		gui.renderHeatmap()
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

	gui.renderHeatmap()
	gui.g.SetViewOnTop("colors")
	gui.renderVersion()
	return nil
}

func (gui *Gui) renderVersion() {
	gui.StatusView.Clear()
	newVersionText := lo.Ternary(newVersionAvailable, "new version available!", "")
	fmt.Fprintf(gui.StatusView, "%s %s", gui.version, newVersionText)
}

func (gui *Gui) renderHeatmap() error {
	v := gui.ViewHeatmap
	v.Clear()

	months := utils.GetMonths(gui.heatmapLastDate)
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

		g.renderHeatmap()

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

	gui.heatmapFirstDate = today
	gui.heatmapLastDate = startDate

	// Create the grid
	grid = make([][]*HeatGrid, 7)
	for i := range grid {
		grid[i] = make([]*HeatGrid, 53)
	}

	defaultTheme := gui.Config.Gui.Theme.Selected
	theme := gui.Config.Gui.Theme.ColorSchemes[defaultTheme]
	heatmaps, _, err := gui.HabitService.HeatMap(context.Background(), startDate, today)
	if err != nil {
		panic(err)
	}

	// Fill the grid with colored boxes
	today = utils.CreateDate(now.Year(), now.Month(), now.Day())
	currentDate := startDate
	for col := 0; col < 53; col++ {
		for row := 0; row < 7; row++ {
			heatGrid := &HeatGrid{
				row:    row,
				column: col,
				key:    currentDate,
			}
			if currentDate.After(today) {
				// grid[row][col] = fmt.Sprintf("\x1b[48;5;52m%d \x1b[0m", currentDate.Day())
				heatGrid.shade = theme.InvalidDayValue
			} else {
				heatmap, ok := heatmaps[currentDate]
				if ok {
					shadeIndex := GetTheShade(heatmap)
					heatGrid.rank = shadeIndex
					colorCode, haveShade := theme.StatusValues[shadeIndex]
					if !haveShade {
						colorCode = theme.ZeroCompletedHabitValue
					}
					heatGrid.shade = colorCode
					heatGrid.haveInfo = true
					heatGrid.totalNumberOfHabits = heatmap.TotalNumberOfHabits
					heatGrid.completedHabits = heatmap.CompletedHabits
				} else {
					heatGrid.shade = theme.NoHabitsValue
				}
			}
			grid[row][col] = heatGrid
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

	gui.heatmapFirstDate = to
	gui.heatmapLastDate = from

	// Create the grid
	grid = make([][]*HeatGrid, 7)
	for i := range grid {
		grid[i] = make([]*HeatGrid, 53)
	}

	defaultTheme := gui.Config.Gui.Theme.Selected
	theme := gui.Config.Gui.Theme.ColorSchemes[defaultTheme]
	heatmaps, _, err := gui.HabitService.HeatMap(context.Background(), from, to)
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
			heatGrid := &HeatGrid{
				row:    row,
				column: col,
				key:    currentDate,
			}
			if (col == 0 && startWeekday > row) || currentDate.Year() != from.Year() {
				heatGrid.shade = theme.InvalidDayValue // not belongs to year, the january starts from the wednesday
				heatGrid.key = time.Time{}
			} else {
				if currentDate.After(today) {
					heatGrid.shade = theme.InvalidDayValue // the now is 2024 july 8, after this time all of them are null
				} else {
					heatmap, ok := heatmaps[currentDate]
					if ok {
						shadeIndex := GetTheShade(heatmap)
						heatGrid.rank = shadeIndex
						colorCode, haveShade := theme.StatusValues[shadeIndex]
						if !haveShade {
							colorCode = theme.ZeroCompletedHabitValue
						}
						heatGrid.shade = colorCode
						heatGrid.haveInfo = true
						heatGrid.totalNumberOfHabits = heatmap.TotalNumberOfHabits
						heatGrid.completedHabits = heatmap.CompletedHabits
					} else {
						heatGrid.shade = theme.NoHabitsValue
					}
				}
				currentDate = currentDate.AddDate(0, 0, 1) // move to the next day
			}
			grid[row][col] = heatGrid
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

func GetTheShade(heatmap *models.HeatMap) int {
	if heatmap.CompletedHabits == 0 {
		return 0
	}

	completionRatio := float64(heatmap.CompletedHabits) / float64(heatmap.TotalNumberOfHabits)
	shade := int(math.Round(completionRatio * 5))

	if shade < 1 {
		shade = 1
	}
	return shade
}

func (gui *Gui) GetDateFromHeatmapCursor() time.Time {
	return grid[cursorY][cursorX].key
}
