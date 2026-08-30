package main

import (
	"fmt"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/window"
)

const (
	placeholderHealth = 100
	placeholderSpeed  = 5
	placeholderCargo  = 250
	placeholderGold   = 0
)

var debugOverlayVisible = false

func shipStatusText() string {
	return fmt.Sprintf(
		"Galleon\nHealth: %d\nSpeed: %d\nCargo: %d\nGold: %d\n",
		placeholderHealth,
		placeholderSpeed,
		placeholderCargo,
		placeholderGold,
	)
}

func examinePanelText(examine entities.ViewableEntity) string {
	return fmt.Sprintf(
		"Captain: %s\nType: %s\nFlag: %s\n",
		examine.GetName(),
		examine.GetType(),
		examine.GetFlag(),
	)
}

func debugOverlayText(playerPos common.Coordinates) string {
	return fmt.Sprintf(
		"Window: %dx%d\nViewport: %dx%d (%dx%d tiles)\nMap: %dx%d\nPlayer: %+v\n",
		window.Window.Width,
		window.Window.Height,
		window.ViewPort.Dimensions.Width,
		window.ViewPort.Dimensions.Height,
		window.ViewPort.Region.Cols,
		window.ViewPort.Region.Rows,
		common.WorldCols,
		common.WorldRows,
		playerPos,
	)
}

func examineActionBarLabel(focused entities.ViewableEntity) string {
	if name := focused.GetName(); name != "" {
		return name
	}
	return "Examine"
}
