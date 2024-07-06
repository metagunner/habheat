package gui

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/metagunner/habheath/pkg/database"
	"github.com/metagunner/habheath/pkg/models"
	"github.com/metagunner/habheath/pkg/utils"
)

type Gui struct {
	g                  *gocui.Gui
	db                 *database.DB
	ViewHeathmap       *gocui.View
	ViewYears          *gocui.View
	mustRenderYears    bool
	mustRenderHeathmap bool
}

func NewGui(db *database.DB) *Gui {
	return &Gui{db: db}
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

	gui.g.FgColor = gocui.ColorDefault
	gui.g.SelFgColor = gocui.ColorGreen | gocui.AttrBold
	gui.g.FrameColor = gocui.ColorDefault
	gui.g.SelFrameColor = gocui.ColorGreen | gocui.AttrBold

	gui.g.SetManager(gocui.ManagerFunc(gui.layout))

	if err := gui.createAllViews(); err != nil {
		return err
	}

	gui.mustRenderHeathmap = true
	gui.mustRenderYears = true

	gui.setKeybindings()

	return gui.g.MainLoop()
}

func (gui *Gui) setKeybindings() error {
	var err error
	err = gui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		return err
	}

	err = gui.g.SetKeybinding("", 1, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gui.nextWindow("heathmap")
	})
	if err != nil {
		return err
	}

	err = gui.g.SetKeybinding("", 2, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gui.nextWindow("years")
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

	if _, err := gui.g.SetCurrentView("years"); err != nil {
		return err
	}

	gui.renderYears()
	gui.renderHeathmap()

	return nil
}

func (gui *Gui) renderYears() error {
	if !gui.mustRenderYears {
		return nil
	}
	gui.mustRenderYears = false

	view := gui.ViewYears
	view.Clear()

	currentYear := time.Now().Year()
	var years []int
	var dayStr string
	var ts time.Time
	err := gui.db.QueryRow(context.Background(), `SELECT day FROM habit ORDER BY day LIMIT 1`).Scan(&dayStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Fprintln(view, currentYear)
			return nil
		}

		return err
	}

	ts, _ = time.Parse(time.RFC3339, dayStr)
	years = utils.GetYearsBetween(ts.Year(), currentYear)
	for _, year := range years {
		fmt.Fprintln(view, year)
	}

	return nil
}

func (gui *Gui) renderHeathmap() error {
	if !gui.mustRenderHeathmap {
		return nil
	}
	gui.mustRenderHeathmap = false

	view := gui.ViewHeathmap
	view.Clear()

	fmt.Fprintln(gui.ViewHeathmap, "======Available Colors======")
	asd := AvailableColors()
	for _, t := range asd {
		for _, c := range t {
			fmt.Fprint(view, c)
		}
		fmt.Fprintln(view, "")
	}

	now := time.Now().UTC()
	lastDay := utils.CreateDate(now.Year(), now.Month(), now.Day())
	months := utils.GetMonths(lastDay)
	for _, month := range months {
		monthName := fmt.Sprintf("   %s   ", month.Format("Jan"))
		fmt.Fprint(view, monthName)
	}
	fmt.Fprintln(view)

	firstDay := months[0]

	var dates []time.Time
	// start from the sunday
	for d := firstDay; !d.After(lastDay); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Sunday {
			firstDay = d
			break
		}
	}
	for d := firstDay; !d.After(lastDay); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d)
	}

	// TODO: store a struct with other related data!
	days := make(map[time.Weekday][]time.Time)
	for _, date := range dates {
		days[date.Weekday()] = append(days[date.Weekday()], date)
	}

	weekDays := []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	}

	var weeks [][]time.Time
	for _, weekDay := range weekDays {
		if dates, ok := days[weekDay]; ok {
			weeks = append(weeks, dates)
		}
	}

	habitService := database.NewHabitService(gui.db)
	habits, _, err := habitService.HeatMap(context.Background(), firstDay, lastDay)
	if err != nil {
		panic(err)
	}

	labels := []string{"   ", "Mon", "   ", "Wed", "   ", "Fri", "   "}
	printLabel := 0
	colors := AvailableColors()
	for _, week := range weeks {
		fmt.Fprintf(view, "%s  ", labels[printLabel])
		for _, day := range week {
			habitDay, ok := habits[day]
			if ok {
				shade := GetTheShade(habitDay)
				fmt.Fprintf(view, "%s", colors["greens"][shade])
			} else {
				fmt.Fprint(view, "  ")
			}
		}
		printLabel += 1
		fmt.Fprintln(view)
	}

	return nil
}

func GetTheShade(heathmap *models.HeathMap) int {
	if heathmap.CompletedHabits == 0 {
		return 0
	}

	//offset := 1
	// 0 1 2 3 4 5+
	//slots := make([]int, 5)
	//for i := 1; i <= 5; i++ {
	//	slots[i-1] = offset * i
	//}

	completionRatio := float64(heathmap.CompletedHabits) / float64(heathmap.TotalNumberOfHabits)
	shade := int(math.Round(completionRatio * 5))

	//	for index, value := range slots {
	//		if input <= value {
	//			return index + 1
	//		}
	//	}
	//return len(slots)
	if shade < 1 {
		shade = 1
	}
	return shade
}

func AvailableColors() map[string]map[int]string {
	space := "  "
	end := "\033[0m"

	// normal colors
	g1 := "\033[48;5;118m" // bright green (less commits)
	g2 := "\033[48;5;40m"
	g3 := "\033[48;5;34m"
	g4 := "\033[48;5;29m"
	g5 := "\033[48;5;22m" // dark green   (more commits)

	// colors for those with dark terminal schemes
	r1 := "\033[48;5;52m" // dark red (less commits)
	r2 := "\033[48;5;88m"
	r3 := "\033[48;5;124m"
	r4 := "\033[48;5;160m"
	r5 := "\033[48;5;196m" // bright red   (more commits)

	status_values := map[string]map[int]string{
		"greens": {
			1: g1 + space + end,
			2: g2 + space + end,
			3: g3 + space + end,
			4: g4 + space + end,
			5: g5 + space + end,
		},
		"reds": {
			1: r1 + space + end,
			2: r2 + space + end,
			3: r3 + space + end,
			4: r4 + space + end,
			5: r5 + space + end,
		},
	}
	return status_values
}
