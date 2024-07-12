package config

// taken from https://github.com/jesseduffield/lazygit/blob/db0a1586d99393cda79e6022f3b3b8b4138b0e8b/pkg/theme/theme.go

import (
	"github.com/jesseduffield/gocui"
)

var gocuiColorMap = map[string]gocui.Attribute{
	"default":   gocui.ColorDefault,
	"black":     gocui.ColorBlack,
	"red":       gocui.ColorRed,
	"green":     gocui.ColorGreen,
	"yellow":    gocui.ColorYellow,
	"blue":      gocui.ColorBlue,
	"magenta":   gocui.ColorMagenta,
	"cyan":      gocui.ColorCyan,
	"white":     gocui.ColorWhite,
	"bold":      gocui.AttrBold,
	"reverse":   gocui.AttrReverse,
	"underline": gocui.AttrUnderline,
}

// GetGocuiAttribute gets the gocui color attribute from the string
func GetGocuiAttribute(key string) gocui.Attribute {
	value, present := gocuiColorMap[key]
	if present {
		return value
	}
	return gocui.ColorWhite
}

// GetGocuiStyle bitwise OR's a list of attributes obtained via the given keys
func GetGocuiStyle(keys []string) gocui.Attribute {
	var attribute gocui.Attribute
	for _, key := range keys {
		attribute |= GetGocuiAttribute(key)
	}
	return attribute
}
