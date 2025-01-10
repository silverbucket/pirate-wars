package terrain

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

// Icon ideas
// Towns: ⩎
// Boats: ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
// People: 옷

const (
	TypeDeepWater    = 0
	TypeOpenWater    = 1
	TypeShallowWater = 2
	TypeBeach        = 3
	TypeLowland      = 4
	TypeHighland     = 5
	TypeRock         = 6
	TypePeak         = 7
	TypeTown         = 8
	TypeGhostTown    = 9
)

type TerrainType int

type TypeQualities struct {
	symbol       rune
	style        lipgloss.Style
	Passable     bool
	RequiresBoat bool
}

func GetType(i int) (TerrainType, error) {
	for k := range TypeLookup {
		if int(k) == i {
			return k, nil
		}
	}
	return 0, errors.New("invalid type")
}

func createTerrainItem(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Background(color).Padding(0).Margin(0)
}

var TypeLookup = map[TerrainType]TypeQualities{
	TypeDeepWater:    {symbol: '⏖', style: createTerrainItem("18"), Passable: true, RequiresBoat: true},
	TypeOpenWater:    {symbol: '⏝', style: createTerrainItem("20"), Passable: true, RequiresBoat: true},
	TypeShallowWater: {symbol: '⏑', style: createTerrainItem("26"), Passable: true, RequiresBoat: true},
	TypeBeach:        {symbol: '~', style: createTerrainItem("#dad1ad"), Passable: true, RequiresBoat: false},
	TypeLowland:      {symbol: ':', style: createTerrainItem("113"), Passable: true, RequiresBoat: false},
	TypeHighland:     {symbol: ':', style: createTerrainItem("142"), Passable: true, RequiresBoat: false},
	TypeRock:         {symbol: '%', style: createTerrainItem("244"), Passable: true, RequiresBoat: false},
	TypePeak:         {symbol: '^', style: createTerrainItem("15"), Passable: false, RequiresBoat: false},
	TypeTown:         {symbol: '⩎', style: createTerrainItem("1"), Passable: true, RequiresBoat: false},
	TypeGhostTown:    {symbol: '⩎', style: createTerrainItem("94"), Passable: true, RequiresBoat: false},
}

func (tt *TerrainType) Render() string {
	return fmt.Sprintf(TypeLookup[*tt].style.PaddingLeft(1).PaddingRight(1).Render("%c"), TypeLookup[*tt].symbol)
}
