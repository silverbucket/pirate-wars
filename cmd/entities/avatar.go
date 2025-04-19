package entities

import (
	"fmt"
	"image/color"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/window"
)

type Avatar struct {
	pos     common.Coordinates
	char    rune
	fgColor color.Color
	bgColor color.Color
	blink   bool
}

type AvatarReadOnly interface {
	GetPos() common.Coordinates
	GetForegroundColor() color.Color
	GetCharacter() string
	GetViewableRange() window.Dimensions
}

type ColorScheme struct {
	Foreground color.Color
	Background color.Color
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

func (a *Avatar) SetBackgroundColor(c color.Color) {
	a.bgColor = c
}

func (a *Avatar) SetForegroundColor(c color.Color) {
	a.fgColor = c
}

func (a *Avatar) GetBackgroundColor() color.Color {
	return a.bgColor
}

func (a *Avatar) GetForegroundColor() color.Color {
	return a.fgColor
}

func (a *Avatar) GetViewableRange() window.Dimensions {
	return window.Dimensions{
		Width:  20,
		Height: 20,
	}
}

func (a *Avatar) GetCharacter() string {
	return fmt.Sprintf("%c", a.char)
}

// func (a *Avatar) Render() *fyne.Container {
// 	return common.RenderContainer(
// 		canvas.NewRectangle(a.bgColor),
// 		canvas.NewText(fmt.Sprintf("%c", a.char), a.fgColor))
// }

func CreateAvatar(coordinates common.Coordinates, c rune, color ColorScheme) Avatar {
	return Avatar{pos: coordinates, char: c, fgColor: color.Foreground, bgColor: color.Background}
}

var black = color.RGBA{0, 0, 0, 255}
var ColorPossibilities = []ColorScheme{
	{color.RGBA{255, 0, 0, 255}, black},     // strong red
	{color.RGBA{0, 255, 0, 255}, black},     // green
	{color.RGBA{255, 255, 0, 255}, black},   // yellow
	{color.RGBA{65, 253, 254, 255}, black},  // bright cyan
	{color.RGBA{248, 240, 227, 255}, black}, // off-white
	{color.RGBA{13, 152, 186, 200}, black},  // blue/green
	{color.RGBA{0, 255, 255, 125}, black},   // faded cyan
	{color.RGBA{224, 255, 255, 255}, black}, // light cyan
	{color.RGBA{255, 119, 255, 255}, black}, // fuchsia
	{color.RGBA{255, 192, 203, 125}, black}, // off pink
	{color.RGBA{144, 238, 144, 125}, black}, // light green
	{color.RGBA{227, 61, 148, 255}, black},  // pink
	{color.RGBA{255, 182, 193, 255}, black}, // light pink
	{color.RGBA{255, 202, 213, 200}, black}, // lighter pink
	{color.RGBA{255, 222, 233, 255}, black}, // light pink/white
	{color.RGBA{255, 213, 128, 255}, black}, // light yellow/orange
	{color.RGBA{255, 255, 125, 255}, black}, // yellow/white
	{color.RGBA{125, 125, 125, 255}, black}, // grey
	{color.RGBA{255, 255, 255, 255}, black}, // white
}
