package terrain

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"pirate-wars/cmd/common"
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
