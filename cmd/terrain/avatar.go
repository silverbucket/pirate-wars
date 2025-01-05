package terrain

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"pirate-wars/cmd/common"
)

type Avatar struct {
	pos     common.Coordinates
	char    rune
	fgColor lipgloss.Color
	bgColor lipgloss.Color
}

func (a *Avatar) GetX() int {
	return a.pos.X
}

func (a *Avatar) GetMiniMapX() int {
	return a.pos.X / common.MiniMapFactor
}

func (a *Avatar) SetX(x int) {
	a.pos.X = x
}

func (a *Avatar) GetY() int {
	return a.pos.Y
}

func (a *Avatar) GetMiniMapY() int {
	return a.pos.Y / common.MiniMapFactor
}

func (a *Avatar) SetY(y int) {
	a.pos.Y = y
}

func (a *Avatar) SetXY(c common.Coordinates) {
	a.pos.X = c.X
	a.pos.Y = c.Y
}

func (a *Avatar) Render() string {
	return fmt.Sprintf(
		lipgloss.NewStyle().
			Foreground(a.fgColor).
			Background(a.bgColor).
			Bold(true).
			Blink(true).
			PaddingLeft(1).PaddingRight(1).Margin(0).
			Render("%c"), a.char)
}

func CreateAvatar(coordinates common.Coordinates, c rune, color ColorScheme) Avatar {
	return Avatar{pos: coordinates, char: c, fgColor: lipgloss.Color(color.Foreground), bgColor: lipgloss.Color(color.Background)}
}
