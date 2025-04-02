package layout

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"math"
	"pirate-wars/cmd/common"
)

var viewableAreaGridCols = 66

// gridLayoutNoPadding defines a grid layout with no padding between cells
type gridLayoutNoPadding struct {
	cols int
}

// newGridLayoutNoPadding creates a new instance of the custom grid layout
func newGridLayoutNoPadding(cols int) fyne.Layout {
	return &gridLayoutNoPadding{cols: cols}
}

// Layout arranges the cells in a grid with no padding
func (g *gridLayoutNoPadding) Layout(cells []fyne.CanvasObject, size fyne.Size) {
	if g.cols < 1 {
		g.cols = 1 // Ensure at least one column
	}

	// Count visible cells
	visibleCount := 0
	for _, obj := range cells {
		if obj.Visible() {
			visibleCount++
		}
	}

	if visibleCount == 0 {
		return
	}

	// Add a small overlap to eliminate gaps (e.g., 1 pixel)
	overlap := float32(1.0)

	// Calculate rows based on visible cells and columns
	rows := (visibleCount + g.cols - 1) / g.cols // Ceiling division

	// Calculate cell size: width = available width / cols, height = available height / rows
	// Use the smaller value to ensure square cells
	cellWidth := size.Width / float32(g.cols)
	cellHeight := size.Height / float32(rows)
	cellSize := float32(math.Min(float64(cellWidth), float64(cellHeight)))

	fmt.Println(fmt.Sprintf("cell width:%v height:%v %+v", cellWidth, cellHeight, g.cols))

	// Position each square cell without padding
	i := 0
	for row := 0; row < rows && i < visibleCount; row++ {
		for col := 0; col < g.cols && i < visibleCount; col++ {
			if !cells[i].Visible() {
				i++
				continue
			}

			x := float32(col) * cellSize
			y := float32(row) * cellSize
			cells[i].Move(fyne.NewPos(x, y))
			cells[i].Resize(fyne.NewSize(cellSize+overlap, cellSize+overlap))
			i++
		}
	}
}

// MinSize calculates the minimum size required for the layout
func (g *gridLayoutNoPadding) MinSize(cells []fyne.CanvasObject) fyne.Size {
	if g.cols < 1 {
		g.cols = 1
	}

	// Find the maximum min size of individual cells (assuming square)
	maxCellSize := float32(0)
	visibleCount := 0
	for _, obj := range cells {
		if obj.Visible() {
			minSize := obj.MinSize()
			cellSize := float32(math.Max(float64(minSize.Width), float64(minSize.Height)))
			if cellSize > maxCellSize {
				maxCellSize = cellSize
			}
			visibleCount++
		}
	}

	if visibleCount == 0 {
		return fyne.NewSize(0, 0)
	}

	// Calculate rows based on visible objects
	rows := (visibleCount + g.cols - 1) / g.cols
	return fyne.NewSize(maxCellSize*float32(g.cols), maxCellSize*float32(rows))
}

func CreateGridContainer(cells []fyne.CanvasObject, isMiniMap bool) *fyne.Container {
	cols := viewableAreaGridCols
	if isMiniMap {
		cols = ViewableArea.Width
	}
	fmt.Println(fmt.Sprintf("init with cols:%v", cols))
	return container.New(
		newGridLayoutNoPadding(cols),
		cells...,
	)
}

func GetCellList(isMiniMap bool) []fyne.CanvasObject {
	cg := Dimensions{
		Width:  common.WorldCols,
		Height: common.WorldRows,
	}
	if !isMiniMap {
		cg = getColGridDimensions(viewableAreaGridCols)
	}
	return make([]fyne.CanvasObject, cg.Width*cg.Height)
}

func getColGridDimensions(cols int) Dimensions {
	return Dimensions{
		Width:  cols,
		Height: int(float32(cols) * (float32(ViewableArea.Height) / float32(ViewableArea.Width))),
	}
}

func CreateGridSimple(cols, rows int, width, height float32) *fyne.Container {
	grid := container.NewWithoutLayout()
	cellWidth := width / float32(cols)
	cellHeight := height / float32(rows)
	fmt.Println(fmt.Sprintf("create grid w:%v h:%v", width, height))
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			var cellColor color.Color
			if (r+c)%2 == 0 {
				cellColor = color.Black
			} else {
				cellColor = color.White
			}
			rect := canvas.NewRectangle(cellColor)
			rect.Resize(fyne.NewSize(cellWidth, cellHeight))
			rect.Move(fyne.NewPos(float32(c)*cellWidth, float32(r)*cellHeight))
			grid.Add(rect)
		}
	}
	return grid
}
