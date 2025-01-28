package npc

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/screen"
)

type Avatar struct {
	pos     common.Coordinates
	char    rune
	fgColor string
	bgColor string
	blink   bool
}

func (a *Avatar) SetPos(c common.Coordinates) {
	a.pos = c
}

func (a *Avatar) GetPos() common.Coordinates {
	return a.pos
}

func (a *Avatar) SetBlink(b bool) {
	a.blink = b
}

func (a *Avatar) SetBackgroundColor(c string) {
	a.bgColor = c
}

func (a *Avatar) GetViewableRange() screen.ViewRange {
	return screen.ViewRange{
		Width:  20,
		Height: 20,
	}
}

func (a *Avatar) Render() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(a.fgColor)).
		Background(lipgloss.Color(a.bgColor)).
		Bold(true).
		Blink(a.blink).
		PaddingLeft(1).PaddingRight(1).Margin(0).
		Render(fmt.Sprintf("%c", a.char))
}

func CreateAvatar(coordinates common.Coordinates, c rune, color ColorScheme) Avatar {
	return Avatar{pos: coordinates, char: c, fgColor: color.Foreground, bgColor: color.Background}
}

var ColorPossibilities = []ColorScheme{
	{"9", "0"},   // strong red
	{"10", "0"},  // green
	{"11", "0"},  // yellow
	{"14", "0"},  // bright cyan
	{"15", "0"},  // off-white
	{"46", "0"},  // blue/green
	{"69", "0"},  // faded cyan
	{"86", "0"},  // light cyan
	{"93", "0"},  // fuchsia
	{"172", "0"}, // off pink
	{"193", "0"}, // light green
	{"201", "0"}, // pink
	{"207", "0"}, // light pink
	{"211", "0"}, // lighter pink
	{"218", "0"}, // light pink/white
	{"222", "0"}, // light yellow/orange
	{"230", "0"}, // yellow/white
	{"253", "0"}, // grey
	{"255", "0"}, // white
}

type ColorScheme struct {
	Foreground string
	Background string
}
