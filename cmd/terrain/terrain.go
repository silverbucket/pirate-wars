package terrain

import (
	"fmt"
	"image/color"
	"pirate-wars/cmd/common"
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

type Terrain struct {
	Cells [common.WorldCols][common.WorldRows]Type
}

type Type int

type TypeQualities struct {
	symbol       rune
	color        color.RGBA
	Passable     bool
	RequiresBoat bool
}

var TypeLookup = map[Type]TypeQualities{
	TypeDeepWater:    {symbol: ' ', color: color.RGBA{2, 0, 121, 255}, Passable: true, RequiresBoat: true},
	TypeOpenWater:    {symbol: ' ', color: color.RGBA{0, 19, 222, 255}, Passable: true, RequiresBoat: true},
	TypeShallowWater: {symbol: ' ', color: color.RGBA{0, 33, 243, 255}, Passable: true, RequiresBoat: true},
	TypeBeach:        {symbol: ' ', color: color.RGBA{205, 170, 109, 125}, Passable: true, RequiresBoat: false},
	TypeLowland:      {symbol: ' ', color: color.RGBA{65, 152, 10, 255}, Passable: true, RequiresBoat: false},
	TypeHighland:     {symbol: ' ', color: color.RGBA{192, 155, 40, 255}, Passable: true, RequiresBoat: false},
	TypeRock:         {symbol: ' ', color: color.RGBA{150, 150, 150, 255}, Passable: true, RequiresBoat: false},
	TypePeak:         {symbol: ' ', color: color.RGBA{229, 229, 229, 255}, Passable: false, RequiresBoat: false},
	TypeTown:         {symbol: '⩎', color: color.RGBA{246, 104, 94, 255}, Passable: true, RequiresBoat: false},
	TypeGhostTown:    {symbol: '⩎', color: color.RGBA{147, 62, 56, 255}, Passable: true, RequiresBoat: false},
}

func (tt *Type) GetBackgroundColor() color.RGBA {
	return TypeLookup[*tt].color
}

func (tt *Type) GetCharacter() string {
	return fmt.Sprintf("%c", TypeLookup[*tt].symbol)
}
