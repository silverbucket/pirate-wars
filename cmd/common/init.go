package common

import (
	"fmt"
	"os"
)

var ViewWidth int
var ViewHeight int
var MiniMapFactor int
var InfoPaneSize int

func SetWindowSize(width int, height int) {

	const startingInfoPaneSize = 5
	const infoPaneSizeIncrements = 1

	if width < 80 || height < 24 {
		fmt.Println("Window size too small. Minimum 80x24")
		os.Exit(1)
	}
	scale := (width - 80) / 20
	scale++
	InfoPaneSize = startingInfoPaneSize + (infoPaneSizeIncrements * scale)
	ViewWidth = (width / 3) - InfoPaneSize
	ViewHeight = height - 1
	//fmt.Println(fmt.Sprintf("Set viewport width to %v of total %v (subtracted info-pane: %v) (scale: %v)", ViewWidth, width, InfoPaneSize, scale))
	CalcMiniMapFactor(scale)
	//fmt.Println(fmt.Sprintf("(MiniMapFactor: %v)", MiniMapFactor))
	//os.Exit(1)
}

func CalcMiniMapFactor(scale int) {
	MiniMapFactor = 24
	for i := 1; i < scale; i++ {
		MiniMapFactor = MiniMapFactor - 2
	}
}
