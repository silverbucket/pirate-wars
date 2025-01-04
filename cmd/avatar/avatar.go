package avatar

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"pirate-wars/cmd/common"
)

type Type struct {
	pos  common.Coordinates
	char rune
}

func (a *Type) GetX() int {
	return a.pos.X
}

func (a *Type) GetMiniMapX() int {
	return a.pos.X / common.MiniMapFactor
}

func (a *Type) SetX(x int) {
	a.pos.X = x
}

func (a *Type) GetY() int {
	return a.pos.Y
}

func (a *Type) GetMiniMapY() int {
	return a.pos.Y / common.MiniMapFactor
}

func (a *Type) SetY(y int) {
	a.pos.Y = y
}

func (a *Type) SetXY(c common.Coordinates) {
	a.pos.X = c.X
	a.pos.Y = c.Y
}

func (a *Type) Render() string {
	return fmt.Sprintf(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#ffffff")).
			Blink(true).
			PaddingLeft(1).PaddingRight(1).Margin(0).
			Render("%c"), a.char)
}

func Create(coordinates common.Coordinates, c rune) Type {
	return Type{pos: coordinates, char: c}
}
