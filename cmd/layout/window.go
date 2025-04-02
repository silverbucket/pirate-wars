package layout

import (
	"math"
	"pirate-wars/cmd/common"
)

type Dimensions struct {
	Width  int
	Height int
}

type Region struct {
	Top    int
	Left   int
	Bottom int
	Right  int
}

var Window Dimensions = Dimensions{
	Width:  1024,
	Height: 768,
}

var InfoPane Dimensions = Dimensions{
	Width:  100,
	Height: Window.Height,
}

var ActionMenu Dimensions = Dimensions{
	Width:  Window.Width,
	Height: 50,
}

var ViewableArea Dimensions = Dimensions{
	Width:  Window.Width - InfoPane.Width,
	Height: Window.Height - ActionMenu.Height,
}

var MiniMapArea Dimensions = Dimensions{
	Width:  700,
	Height: 700,
}

var WorldViewport Region = Region{
	Top:    0,
	Left:   0,
	Bottom: common.WorldRows - 1,
	Right:  common.WorldCols - 1,
}

type MapIncrementFactor struct {
	X int
	Y int
}

var MiniMapIncrementFactor = MapIncrementFactor{
	int(math.Ceil(float64(common.WorldCols/viewableAreaGridCols))) + 2,
	int(math.Ceil(float64(common.WorldRows/viewableAreaGridCols))) + 5,
}

func GetMiniMapScale(c common.Coordinates) common.Coordinates {
	return common.Coordinates{c.X / MiniMapIncrementFactor.X, c.Y / MiniMapIncrementFactor.Y}
}

func CalcViewport(pos common.Coordinates) Region {
	// viewable range is based on columns in grid and ratio of ViewableArea
	vr := getColGridDimensions(viewableAreaGridCols)

	// center viewport on position
	left := pos.X - (vr.Width / 2)
	right := pos.X + (vr.Width / 2)

	top := pos.Y - (vr.Height / 2)
	bottom := pos.Y + (vr.Height / 2)

	// take up screen
	if right-left < vr.Width {
		left = right - vr.Width
	}
	if bottom-top < vr.Height {
		top = bottom - vr.Height
	}

	// don't slide the screen when you hit the edge
	if bottom >= common.WorldRows {
		bottom = common.WorldRows
		top = common.WorldRows - vr.Height
	}
	if right >= common.WorldCols {
		right = common.WorldCols
		left = common.WorldCols - vr.Width
	}

	if left < 0 {
		left = 0
		right = vr.Width
	}
	if top < 0 {
		top = 0
		bottom = vr.Height
	}

	return Region{top, left, bottom, right}
}

func (v *Region) IsPositionWithin(c common.Coordinates) bool {
	if (v.Left <= c.X && c.X <= v.Right) && (v.Top <= c.Y && c.Y <= v.Bottom) {
		return true
	}
	return false
}
