package resources

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"pirate-wars/cmd/common"
)

// TileSize represents the size of each tile in the tileset
const TileSize = 32

// TileMapping maps terrain types to tile coordinates in the tileset
var TileMapping = map[int]image.Point{
	common.TerrainTypeDeepWater:    {X: 0, Y: 1}, // Deep water
	common.TerrainTypeOpenWater:    {X: 1, Y: 1}, // Open water
	common.TerrainTypeShallowWater: {X: 2, Y: 1}, // Shallow water
	common.TerrainTypeBeach:        {X: 1, Y: 0}, // Beach
	common.TerrainTypeLowland:      {X: 0, Y: 0}, // Lowland
	common.TerrainTypeHighland:     {X: 2, Y: 0}, // Highland
	common.TerrainTypeRock:         {X: 3, Y: 0}, // Rock
	common.TerrainTypePeak:         {X: 4, Y: 0}, // Peak
	common.TerrainTypeTown:         {X: 3, Y: 1}, // Town
	common.TerrainTypeGhostTown:    {X: 4, Y: 1}, // Ghost town
	common.TerrainTypeLowlandBrush: {X: 5, Y: 0}, // Lowland brush
	common.ShipWhite:               {X: 0, Y: 2},
	common.ShipPirate:              {X: 1, Y: 2},
	common.ShipRed:                 {X: 2, Y: 2},
	common.ShipGreen:               {X: 3, Y: 2},
	common.ShipBlue:                {X: 4, Y: 2},
	common.ShipYellow:              {X: 5, Y: 2},
}

var (
	tilesetCache image.Image
	tileCache    = make(map[int]image.Image)
)

func getTileByRegion(idx int) image.Image {
	// Get the tile coordinates from the mapping
	tileCoords, ok := TileMapping[idx]
	if !ok {
		fmt.Printf("No mapping found for ship type %v, using default\n", idx)
		tileCoords = TileMapping[common.ShipWhite]
	}

	// Check if we have this tile cached
	if cached, ok := tileCache[idx]; ok {
		return cached
	}
	fmt.Printf("Getting tile %v at coordinates (%d,%d)\n", idx, tileCoords.X, tileCoords.Y)

	tileset := getTileset()

	// Create a new RGBA image for the tile
	tileImg := image.NewRGBA(image.Rect(0, 0, TileSize, TileSize))

	// Calculate the source coordinates in the tileset
	srcX := tileCoords.X * TileSize
	srcY := tileCoords.Y * TileSize

	// Copy pixels directly from the tileset to our tile image
	for y := 0; y < TileSize; y++ {
		for x := 0; x < TileSize; x++ {
			srcPixel := tileset.At(srcX+x, srcY+y)
			tileImg.Set(x, y, srcPixel)
		}
	}

	// Cache the tile
	tileCache[idx] = tileImg
	return tileImg
}

func GetShipTile(s common.ShipType) image.Image {
	idx := int(s)
	return getTileByRegion(idx)
}

// GetTerrainTile returns the image for a specific terrain type
func GetTerrainTile(tt common.TerrainType) image.Image {
	idx := int(tt)
	return getTileByRegion(idx)
}

// loadTilesetImage loads the tileset image from the bundled resources
func loadTilesetImage() (image.Image, error) {
	// Get the tileset data from the bundled resource
	tilesetData := resourcePirateWarsTilesetPng.StaticContent

	// Decode the PNG data into an image
	img, err := png.Decode(bytes.NewReader(tilesetData))
	if err != nil {
		return nil, fmt.Errorf("error decoding tileset: %v", err)
	}

	return img, nil
}

func getTileset() image.Image {
	if tilesetCache == nil {
		var err error
		tilesetCache, err = loadTilesetImage()
		if err != nil {
			fmt.Printf("Error loading tileset: %v\n", err)
			return image.NewRGBA(image.Rect(0, 0, TileSize, TileSize))
		}
		fmt.Printf("Loaded tileset image with bounds: %v\n", tilesetCache.Bounds())
	}
	return tilesetCache
}
