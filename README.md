# Habheat
A simple terminal UI for habit tracking with Github like heat map.

![green theme](docs/assets/heatmap-green-scheme.png)

- [Table of contents](#table-of-contents)
- [Features](#features)
  - [Year Selection](#year-selection)
  - [Habheat Grid](#habheat-grid)
    - [Create Habit](#create-habit)
    - [Remove Habit](#remove-habit)
    - [Toggle Habit](#toggle-habit)
    - [Update Habit](#update-habit)
- [Installation](#installation)
  - [Binary Releases](#binary-releases)
  - [Go](#go)
  - [Manual](#manual)
- [Usage](#usage)
- [Configuration](#configuration)
  - [Custom Theme](#custom-theme)
  - [Keybindings](#keybinding)

## Features

### Year Selection
Press `space` to select the year. The `Default` option shows 12 months starting today and going backwards. For example, for *July 13 2024*, the grid will show 53 weeks starting from the current week to the 2023.

### Habheat Grid
The grid displays colors based on the habit completion ratio. There are [built-in color schemes](#built-in-color-schemes) for the grid, and you can also [create your own](#custom-color-scheme). Each cell in the grid represents a day. You can navigate through the grid and see the habits for any day by pressing `space`. This will open a popup where you can edit the habits. 

#### Create Habit
You can create a new habit by pressing `n` on the habit popup. It will ask for the title of the habit, after writing your title you can press `enter` to confirm.

#### Remove Habit
Press `r` on a habit to remove it. Be careful, as this action cannot be undone.

#### Toggle Habit
Press `space` on a habit to toggle its completion status. This will affect the color in the heat map.

#### Update Habit
Press `u` on a habit to edit its title.

## Installation

### Binary Releases
For Windows, Mac OS(10.12+) or Linux, you can download a binary release [here](https://github.com/metagunner/habheat/releases).

### Go
```sh
go install github.com/metagunner/habheat@latest
```

### Manual

You'll need to [install Go](https://golang.org/doc/install)

```
git clone https://github.com/metagunner/habheat.git
cd habheat
go install
```

You can also use `go run main.go` to run the application.

## Usage
Call `habheat` in your terminal.

```sh
$ habheat
```

## Configuration

Default path for the config file and the database:

- Linux: `~/.config/habheat/config.yml`
- MacOS: `~/Library/Application\ Support/habheat/config.yml`
- Windows: `%LOCALAPPDATA%\habheat\config.yml` (default location, but it will also be found in `%APPDATA%\habheat\config.yml`

<!-- START CONFIG YAML: AUTOMATICALLY GENERATED DO NOT UPDATE MANUALLY -->
### Defaults
```yaml
# Config relating to the Habheat UI
gui:
    # Config relating to colors and styles.
    theme:
        # Selected heat map color scheme.
        selected: green

        # Available heat map color schemes.
        colorSchemes:
            # Color scheme name. Set this to "selected" property.
            green:
                # Value for cell that is in the future or not applicable for the year.
                invalidDayValue: '  '

                # Value for a cell with no habits.
                noHabitsValue: '  '

                # Value for a cell with no completed habits.
                zeroCompletedHabitValue: '  '

                # Color shades with ANSI codes, from less to more
                statusValues:
                    1: "\e[48;5;22m  \e[0m"
                    2: "\e[48;5;29m  \e[0m"
                    3: "\e[48;5;34m  \e[0m"
                    4: "\e[48;5;40m  \e[0m"
                    5: "\e[48;5;118m  \e[0m"

                # Value for the cursor
                cursorValue: "\e[48;5;196m  \e[0m"

            # Blue color scheme
            ice:
                invalidDayValue: '  '
                noHabitsValue: '  '
                zeroCompletedHabitValue: '  '
                statusValues:
                    1: "\e[48;5;61m  \e[0m"
                    2: "\e[48;5;67m  \e[0m"
                    3: "\e[48;5;68m  \e[0m"
                    4: "\e[48;5;74m  \e[0m"
                    5: "\e[48;5;75m  \e[0m"
                cursorValue: "\e[48;5;196m  \e[0m"

            # Purple color scheme
            purple:
                invalidDayValue: '  '
                noHabitsValue: '  '
                zeroCompletedHabitValue: '  '
                statusValues:
                    1: "\e[48;5;55m  \e[0m"
                    2: "\e[48;5;92m  \e[0m"
                    3: "\e[48;5;93m  \e[0m"
                    4: "\e[48;5;129m  \e[0m"
                    5: "\e[48;5;135m  \e[0m"
                cursorValue: "\e[48;5;196m  \e[0m"

            # Yellow color scheme
            yellow:
                invalidDayValue: '  '
                noHabitsValue: '  '
                zeroCompletedHabitValue: '  '
                statusValues:
                    1: "\e[48;5;142m  \e[0m"
                    2: "\e[48;5;178m  \e[0m"
                    3: "\e[48;5;184m  \e[0m"
                    4: "\e[48;5;220m  \e[0m"
                    5: "\e[48;5;226m  \e[0m"
                cursorValue: "\e[48;5;196m  \e[0m"

        # Border color of the focused window
        activeBorderColor:
            - green
            - bold
            
        # Border color of non-focused windows
        inactiveBorderColor:
            - default
```

### Built-in Color Schemes
These are all the available color schemes. The default one is *green*. You can change the color scheme in the configuration.
```yaml
gui:
    theme:
        selected: purple # The color scheme you want
```
![green theme](docs/assets/heatmap-green-scheme.png)
![ice theme](docs/assets/heatmap-ice-scheme.png)
![purple theme](docs/assets/heatmap-purple-scheme.png)
![yellow theme](docs/assets/heatmap-yellow-scheme.png)

### Customization
The built-in color schemes use ANSI color codes to support most terminals. Depending on your terminal's capabilities, you can also use Unicode characters.

#### Custom Color Scheme
"Selected" property should be same with the name of the color scheme.
```yaml
gui:
  theme:
    selected: custom
    colorSchemes: 
      custom:
        invalidDayValue: "🚫"
        noHabitsValue: "  "
        zeroCompletedHabitValue: "  "
        statusValues:
          1: 😭
          2: 🥺
          3: 😎
          4: 😳
          5: 🤩
        cursorValue: 🖖

```
![custom theme](docs/assets/custom-scheme.png)

## Color Attributes

The available color attributes are:

**Colors**

- black
- red
- green
- yellow
- blue
- magenta
- cyan
- white

**Modifiers**

- bold
- default
- reverse # useful for high-contrast
- underline

## Keybindings
At the moment, there are no custom keybindings. It is on the to-do list.
<!-- For all keybinding options check [Keybindings](./Keybindings.md). -->