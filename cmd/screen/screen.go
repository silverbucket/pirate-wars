package screen

import (
	"fmt"
	"os"
)

type ViewRange struct {
	Width  int
	Height int
}

var Dimensions ViewRange

var MiniMapFactor int
var InfoPaneSize int

func SetWindowSize(width int, height int) {
	const startingInfoPaneSize = 5
	const infoPaneSizeIncrements = 1

	if width < 80 || height < 24 {
		fmt.Println("Window size too small. Minimum 80x24")
		os.Exit(1)
	}

	// padding needed for consistent rendering among different terminals
	usefulWidth := width - 4

	scale := (usefulWidth - 80) / 20
	scale++
	InfoPaneSize = startingInfoPaneSize + (infoPaneSizeIncrements * scale)
	Dimensions.Width = (usefulWidth / 3) - InfoPaneSize
	Dimensions.Height = height - scale
	CalcMiniMapFactor(scale)
}

func CalcMiniMapFactor(scale int) {
	MiniMapFactor = 24
	for i := 1; i < scale; i++ {
		MiniMapFactor = MiniMapFactor - 2
	}
}
