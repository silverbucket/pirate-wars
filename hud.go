package main

import (
	"fmt"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/sailing"
	"pirate-wars/cmd/window"
)

const (
	placeholderHealth = 100
	placeholderCargo  = 250
	placeholderGold   = 0
)

var debugOverlayVisible = false

func shipStatusText(speed float64, wind *sailing.Wind) string {
	windLabel := "—"
	windStrength := 0
	if wind != nil {
		windLabel = wind.Label()
		windStrength = wind.Strength
	}
	return fmt.Sprintf(
		"Galleon\nHealth: %d\nSpeed: %.2f\nWind: %s (%d)\nCargo: %d\nGold: %d\n",
		placeholderHealth,
		speed,
		windLabel,
		windStrength,
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

func debugOverlayText(playerPos common.Coordinates, wind *sailing.Wind) string {
	windInfo := "wind: n/a"
	if wind != nil {
		windInfo = fmt.Sprintf("wind: %s strength %d", wind.Label(), wind.Strength)
	}
	return fmt.Sprintf(
		"Window: %dx%d\nViewport: %dx%d (%dx%d tiles)\nMap: %dx%d\nPlayer: %+v\n%s\n",
		window.Window.Width,
		window.Window.Height,
		window.ViewPort.Dimensions.Width,
		window.ViewPort.Dimensions.Height,
		window.ViewPort.Region.Cols,
		window.ViewPort.Region.Rows,
		common.WorldCols,
		common.WorldRows,
		playerPos,
		windInfo,
	)
}

func examineActionBarLabel(focused entities.ViewableEntity) string {
	if name := focused.GetName(); name != "" {
		return name
	}
	return "Examine"
}
