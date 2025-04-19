package window

import (
	"pirate-wars/cmd/common"
)

type Dimensions struct {
	Width  int
	Height int
}

type Region struct {
	X, Y       int
	Cols, Rows int
}

type DimensionsAndRegion struct {
	Dimensions Dimensions
	Region     Region
}

var Window Dimensions = Dimensions{
	Width:  1024,
	Height: 768,
}

var SidePanel Dimensions = Dimensions{
	Width:  170,
	Height: Window.Height,
}

var ActionMenu Dimensions = Dimensions{
	Width:  Window.Width - SidePanel.Width,
	Height: 68,
}

var viewPortWidth = Window.Width - SidePanel.Width
var viewPortHeight = Window.Height - ActionMenu.Height + 28
var ViewPort DimensionsAndRegion = DimensionsAndRegion{
	Dimensions: Dimensions{
		Width:  viewPortWidth,
		Height: viewPortHeight,
	},
	Region: Region{
		Cols: viewPortWidth / CellSize,
		Rows: viewPortHeight / CellSize,
	},
}

var MiniMapArea Dimensions = Dimensions{
	Width:  700,
	Height: 700,
}

var CellSize = 20

func GetViewportRegion(pos common.Coordinates) Region {
	// viewable range is based on columns in grid and ratio of ViewableArea
	vp := Region{
		Cols: ViewPort.Region.Cols,
		Rows: ViewPort.Region.Rows,
		X:    int(pos.X - ViewPort.Region.Cols/2),
		Y:    int(pos.Y - ViewPort.Region.Rows/2),
	}

	if vp.X < 0 {
		vp.X = 0
	} else if vp.X+vp.Cols > common.WorldCols {
		vp.X = common.WorldCols - vp.Cols
	}
	if vp.Y < 0 {
		vp.Y = 0
	} else if vp.Y+vp.Rows > common.WorldRows {
		vp.Y = common.WorldRows - vp.Rows
	}
	return vp
}

func (v *Region) IsPositionWithin(c common.Coordinates) bool {
	if (v.X <= c.X && c.X <= v.X+v.Cols) && (v.Y <= c.Y && c.Y <= v.Y+v.Rows) {
		return true
	}
	return false
}
