package config

type UserConfig struct {
	Gui        GuiConfig        `yaml:"gui"`
	Keybinding KeybindingConfig `yaml:"keybinding"`
}

type GuiConfig struct {
	Theme ThemeConfig `yaml:"theme"`
}

type ThemeConfig struct {
	Selected            string                        `yaml:"selected"`
	ColorSchemes        map[string]HeatmapColorScheme `yaml:"colorSchemes"`
	ActiveBorderColor   []string                      `yaml:"activeBorderColor"`
	InactiveBorderColor []string                      `yaml:"inactiveBorderColor"`
}

type HeatmapColorScheme struct {
	InvalidDayValue         string         `yaml:"invalidDayValue"`
	NoHabitsValue           string         `yaml:"noHabitsValue"`
	ZeroCompletedHabitValue string         `yaml:"zeroCompletedHabitValue"`
	StatusValues            map[int]string `yaml:"statusValues"`
	CursorValue             string         `yaml:"cursorValue"`
}

type KeybindingConfig struct {
	Universal KeybindingUniversalConfig `yaml:"universal"`
	Heatmap   KeybindingHeatmapConfig   `yaml:"heatmap"`
}

type KeybindingUniversalConfig struct {
	Quit        string `yaml:"quit"`
	PrevItem    string `yaml:"prevItem"`
	NextItem    string `yaml:"nextItem"`
	PrevItemAlt string `yaml:"prevItemAlt"`
	NextItemAlt string `yaml:"nextItemAlt"`
	Select      string `yaml:"select"`
	Confirm     string `yaml:"confirm"`
	Close       string `yaml:"close"`
}

type KeybindingHeatmapConfig struct {
	Right       string `yaml:"right"`
	Left        string `yaml:"left"`
	Up          string `yaml:"up"`
	Down        string `yaml:"down"`
	RightAlt    string `yaml:"rightAlt"`
	LeftAlt     string `yaml:"leftAlt"`
	UpAlt       string `yaml:"upAlt"`
	DownAlt     string `yaml:"downAlt"`
	EditHabit   string `yaml:"editHabit"`
	ToggleHabit string `yaml:"toggleHabit"`
	CreateHabit string `yaml:"createHabit"`
	DeleteHabit string `yaml:"deleteHabit"`
}

const (
	space = "  "
	end   = "\033[0m"
)

func GetDefaultConfig() *UserConfig {
	return &UserConfig{
		Gui: GuiConfig{
			Theme: ThemeConfig{
				ActiveBorderColor:   []string{"green", "bold"},
				InactiveBorderColor: []string{"default"},
				Selected:            "green",
				ColorSchemes: map[string]HeatmapColorScheme{
					"green": {
						InvalidDayValue: "  ", NoHabitsValue: "  ", ZeroCompletedHabitValue: "  ", StatusValues: map[int]string{
							1: "\033[48;5;22m" + space + end,
							2: "\033[48;5;29m" + space + end,
							3: "\033[48;5;34m" + space + end,
							4: "\033[48;5;40m" + space + end,
							5: "\033[48;5;118m" + space + end,
						},
						CursorValue: "\033[48;5;196m" + space + end,
					},
					"purple": {
						InvalidDayValue: "  ", NoHabitsValue: "  ", ZeroCompletedHabitValue: "  ", StatusValues: map[int]string{
							1: "\033[48;5;55m" + space + end,
							2: "\033[48;5;92m" + space + end,
							3: "\033[48;5;93m" + space + end,
							4: "\033[48;5;129m" + space + end,
							5: "\033[48;5;135m" + space + end,
						},
						CursorValue: "\033[48;5;196m" + space + end,
					},
					"yellow": {
						InvalidDayValue: "  ", NoHabitsValue: "  ", ZeroCompletedHabitValue: "  ", StatusValues: map[int]string{
							1: "\033[48;5;142m" + space + end,
							2: "\033[48;5;178m" + space + end,
							3: "\033[48;5;184m" + space + end,
							4: "\033[48;5;220m" + space + end,
							5: "\033[48;5;226m" + space + end,
						},
						CursorValue: "\033[48;5;196m" + space + end,
					},
					"ice": {
						InvalidDayValue: "  ", NoHabitsValue: "  ", ZeroCompletedHabitValue: "  ", StatusValues: map[int]string{
							1: "\033[48;5;61m" + space + end,
							2: "\033[48;5;67m" + space + end,
							3: "\033[48;5;68m" + space + end,
							4: "\033[48;5;74m" + space + end,
							5: "\033[48;5;75m" + space + end,
						},
						CursorValue: "\033[48;5;196m" + space + end,
					},
				},
			},
		},
		Keybinding: KeybindingConfig{
			Universal: KeybindingUniversalConfig{
				Quit:        "q",
				PrevItem:    "<up>",
				NextItem:    "<down>",
				PrevItemAlt: "k",
				NextItemAlt: "j",
				Select:      "<space>",
				Confirm:     "<enter>",
				Close:       "<esc>",
			},
			Heatmap: KeybindingHeatmapConfig{
				Right:       "l",
				Left:        "h",
				Up:          "k",
				Down:        "j",
				RightAlt:    "<right>",
				LeftAlt:     "<left>",
				UpAlt:       "<up>",
				DownAlt:     "<down>",
				EditHabit:   "u",
				ToggleHabit: "<space>",
				CreateHabit: "n",
				DeleteHabit: "r",
			},
		},
	}
}
