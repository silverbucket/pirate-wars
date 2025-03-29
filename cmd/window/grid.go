package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

var viewableAreaGridCols = 30

// customGridLayoutNoPadding defines a grid layout with no padding between cells
type customGridLayoutNoPadding struct {
	cols int // Number of columns
}

// NewCustomGridLayoutNoPadding creates a new instance of the custom grid layout
func newViewableAreaLayout() fyne.Layout {
	return &customGridLayoutNoPadding{cols: viewableAreaGridCols}
}

// Layout arranges the cells in a grid with no padding
func (g *customGridLayoutNoPadding) Layout(cells []fyne.CanvasObject, size fyne.Size) {
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

	// Calculate rows based on visible cells and columns
	rows := (visibleCount + g.cols - 1) / g.cols // Ceiling division
	cellWidth := int(size.Width / float32(g.cols))
	cellHeight := int(size.Height / float32(rows))

	// Position each cell without padding
	i := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < g.cols; col++ {
			if i >= visibleCount {
				break
			}
			if !cells[i].Visible() {
				i++
				continue
			}

			x := col * cellWidth
			y := row * cellHeight
			cells[i].Move(fyne.NewPos(float32(x), float32(y)))
			cells[i].Resize(fyne.NewSize(float32(cellWidth), float32(cellHeight)))
			i++
		}
	}
}

// MinSize calculates the minimum size required for the layout
func (g *customGridLayoutNoPadding) MinSize(cells []fyne.CanvasObject) fyne.Size {
	if g.cols < 1 {
		g.cols = 1
	}

	// Find the maximum min width and height of individual cells
	maxCellWidth := float32(0)
	maxCellHeight := float32(0)
	visibleCount := 0

	for _, obj := range cells {
		if obj.Visible() {
			minSize := obj.MinSize()
			if minSize.Width > maxCellWidth {
				maxCellWidth = minSize.Width
			}
			if minSize.Height > maxCellHeight {
				maxCellHeight = minSize.Height
			}
			visibleCount++
		}
	}

	if visibleCount == 0 {
		return fyne.NewSize(0, 0)
	}

	// Calculate rows based on visible objects
	rows := (visibleCount + g.cols - 1) / g.cols
	return fyne.NewSize(maxCellWidth*float32(g.cols), maxCellHeight*float32(rows))
}

func CreateGridContainer(cells []fyne.CanvasObject) *fyne.Container {
	return container.New(
		newViewableAreaLayout(),
		cells...,
	)
}

func GetCellList() []fyne.CanvasObject {
	cg := getColgridDimensions()
	return make([]fyne.CanvasObject, cg.Width*cg.Height)
}

func getColgridDimensions() Dimensions {
	return Dimensions{
		Width:  viewableAreaGridCols,
		Height: int(float32(viewableAreaGridCols) * (float32(ViewableArea.Height) / float32(ViewableArea.Width))),
	}
}
