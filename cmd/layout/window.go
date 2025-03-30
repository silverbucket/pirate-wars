package layout

import (
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

//	func SetWindowSize(width int, height int) {
//		const startingInfoPaneSize = 5
//		const infoPaneSizeIncrements = 1
//
//		if width < 80 || height < 24 {
//			fmt.Println("Window size too small. Minimum 80x24")
//			os.Exit(1)
//		}
//
//		// padding needed for consistent rendering among different terminals
//		usefulWidth := width - 4
//
//		scale := (usefulWidth - 80) / 20
//		scale++
//		InfoPaneSize = startingInfoPaneSize + (infoPaneSizeIncrements * scale)
//		Dimensions.Width = (usefulWidth / 3) - InfoPaneSize
//		Dimensions.Height = height - scale
//		CalcMiniMapFactor(scale)
//	}
var MiniMapFactor int

func CalcMiniMapFactor(scale int) {
	MiniMapFactor = 24
	for i := 1; i < scale; i++ {
		MiniMapFactor = MiniMapFactor - 2
	}
}
func GetMiniMapScale(c common.Coordinates) common.Coordinates {
	return common.Coordinates{c.X / MiniMapFactor, c.Y / MiniMapFactor}
}

func CalcViewport(pos common.Coordinates) Region {
	// viewable range is based on columns in grid and ratio of ViewableArea
	vr := GetColGridDimensions()

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
	if bottom >= common.WorldHeight {
		bottom = common.WorldHeight
		top = common.WorldHeight - vr.Height
	}
	if right >= common.WorldWidth {
		right = common.WorldWidth
		left = common.WorldWidth - vr.Width
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
