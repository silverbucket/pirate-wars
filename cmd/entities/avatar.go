package entities

import (
	"image"
	"image/color"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/window"
)

type Avatar struct {
	id        string
	pos       common.Coordinates
	prevPos   common.Coordinates
	image     image.Image
	color     color.Color
	blink     bool
	alternate bool
}

type AvatarReadOnly interface {
	GetID() string
	GetPos() common.Coordinates
	GetPreviousPos() common.Coordinates
	GetTileImage() image.Image
	GetViewableRange() window.Dimensions
	IsHighlighted() bool
	GetColor() color.Color
}

func (a *Avatar) GetID() string {
	return a.id
}

func (a *Avatar) SetPos(c common.Coordinates) {
	if !common.CoordsMatch(a.pos, c) {
		a.prevPos = a.pos
		a.pos = c
	}
}

func (a *Avatar) GetPos() common.Coordinates {
	return a.pos
}

func (a *Avatar) GetPreviousPos() common.Coordinates {
	return a.prevPos
}

func (a *Avatar) Highlight(b bool) {
	a.blink = b
	a.alternate = b
}

func (a *Avatar) IsHighlighted() bool {
	return a.blink
}

func (a *Avatar) GetColor() color.Color {
	if a.blink {
		if !a.alternate {
			a.alternate = true
			return color.RGBA{0, 0, 0, 0}
		}
	}
	a.alternate = false
	return a.color
}

func (a *Avatar) GetViewableRange() window.Dimensions {
	return window.Dimensions{
		Width:  20,
		Height: 20,
	}
}

func (a *Avatar) GetTileImage() image.Image {
	return a.image
}

func CreateAvatar(pos common.Coordinates, i image.Image, c color.Color) Avatar {
	return Avatar{
		id:  common.GenID(pos),
		pos: pos, image: i, color: c,
	}
}

var ColorPossibilities = []color.Color{
	color.RGBA{255, 0, 0, 255},     // strong red
	color.RGBA{0, 255, 0, 255},     // green
	color.RGBA{255, 255, 0, 255},   // yellow
	color.RGBA{65, 253, 254, 255},  // bright cyan
	color.RGBA{248, 240, 227, 255}, // off-white
	color.RGBA{13, 152, 186, 200},  // blue/green
	color.RGBA{0, 255, 255, 125},   // faded cyan
	color.RGBA{224, 255, 255, 255}, // light cyan
	color.RGBA{255, 119, 255, 255}, // fuchsia
	color.RGBA{255, 192, 203, 125}, // off pink
	color.RGBA{144, 238, 144, 125}, // light green
	color.RGBA{227, 61, 148, 255},  // pink
	color.RGBA{255, 182, 193, 255}, // light pink
	color.RGBA{255, 202, 213, 200}, // lighter pink
	color.RGBA{255, 222, 233, 255}, // light pink/white
	color.RGBA{255, 213, 128, 255}, // light yellow/orange
	color.RGBA{255, 255, 125, 255}, // yellow/white
	color.RGBA{125, 125, 125, 255}, // grey
	color.RGBA{255, 255, 255, 255}, // white
}
