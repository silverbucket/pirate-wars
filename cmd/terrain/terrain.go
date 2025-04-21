package terrain

import (
	"image"
	"image/color"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/resources"
)

// Icon ideas
// Towns: ⩎
// Boats: ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
// People: 옷

type Terrain struct {
	Cells [common.WorldCols][common.WorldRows]common.TerrainType
}

type TypeQualities struct {
	color        color.RGBA
	Passable     bool
	RequiresBoat bool
	tile         image.Image
}

var TypeLookup = map[common.TerrainType]TypeQualities{
	common.TerrainTypeDeepWater:    {color: color.RGBA{2, 0, 121, 255}, tile: resources.GetTerrainTile(common.TerrainTypeDeepWater), Passable: true, RequiresBoat: true},
	common.TerrainTypeOpenWater:    {color: color.RGBA{0, 19, 222, 255}, tile: resources.GetTerrainTile(common.TerrainTypeOpenWater), Passable: true, RequiresBoat: true},
	common.TerrainTypeShallowWater: {color: color.RGBA{0, 33, 243, 255}, tile: resources.GetTerrainTile(common.TerrainTypeShallowWater), Passable: true, RequiresBoat: true},
	common.TerrainTypeBeach:        {color: color.RGBA{205, 170, 109, 125}, tile: resources.GetTerrainTile(common.TerrainTypeBeach), Passable: true, RequiresBoat: false},
	common.TerrainTypeLowland:      {color: color.RGBA{65, 152, 10, 255}, tile: resources.GetTerrainTile(common.TerrainTypeLowland), Passable: true, RequiresBoat: false},
	common.TerrainTypeHighland:     {color: color.RGBA{192, 155, 40, 255}, tile: resources.GetTerrainTile(common.TerrainTypeHighland), Passable: true, RequiresBoat: false},
	common.TerrainTypeRock:         {color: color.RGBA{150, 150, 150, 255}, tile: resources.GetTerrainTile(common.TerrainTypeRock), Passable: true, RequiresBoat: false},
	common.TerrainTypePeak:         {color: color.RGBA{229, 229, 229, 255}, tile: resources.GetTerrainTile(common.TerrainTypePeak), Passable: false, RequiresBoat: false},
	common.TerrainTypeTown:         {color: color.RGBA{246, 104, 94, 255}, tile: resources.GetTerrainTile(common.TerrainTypeTown), Passable: true, RequiresBoat: false},
	common.TerrainTypeGhostTown:    {color: color.RGBA{147, 62, 56, 255}, tile: resources.GetTerrainTile(common.TerrainTypeGhostTown), Passable: true, RequiresBoat: false},
}

func GetColor(tt common.TerrainType) color.RGBA {
	return TypeLookup[tt].color
}

func GetTile(tt common.TerrainType) image.Image {
	return TypeLookup[tt].tile
}

func IsPassable(tt common.TerrainType) bool {
	return TypeLookup[tt].Passable
}

func RequiresBoat(tt common.TerrainType) bool {
	return TypeLookup[tt].RequiresBoat
}
