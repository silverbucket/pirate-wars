package resources

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"pirate-wars/cmd/terrain"
)

// TileSize represents the size of each tile in the tileset
const TileSize = 32

// TileMapping maps terrain types to tile coordinates in the tileset
var TileMapping = map[terrain.Type]image.Point{
	terrain.TypeDeepWater:    {X: 0, Y: 1}, // Deep water
	terrain.TypeOpenWater:    {X: 1, Y: 1}, // Open water
	terrain.TypeShallowWater: {X: 2, Y: 1}, // Shallow water
	terrain.TypeBeach:        {X: 1, Y: 0}, // Beach
	terrain.TypeLowland:      {X: 0, Y: 0}, // Lowland
	terrain.TypeHighland:     {X: 2, Y: 0}, // Highland
	terrain.TypeRock:         {X: 3, Y: 0}, // Rock
	terrain.TypePeak:         {X: 4, Y: 0}, // Peak
	terrain.TypeTown:         {X: 3, Y: 1}, // Town
	terrain.TypeGhostTown:    {X: 4, Y: 1}, // Ghost town
}

var (
	tilesetCache image.Image
	tileCache    = make(map[terrain.Type]image.Image)
)

// GetTileImage returns the image for a specific terrain type
func GetTileImage(tt terrain.Type) image.Image {
	// Check if we have this tile cached
	if cached, ok := tileCache[tt]; ok {
		return cached
	}

	// Get the tile coordinates from the mapping
	tileCoords, ok := TileMapping[tt]
	if !ok {
		fmt.Printf("No mapping found for terrain type %v, using deep water\n", tt)
		tileCoords = TileMapping[terrain.TypeDeepWater]
	}

	fmt.Printf("Getting tile for terrain type %v at coordinates (%d,%d)\n", tt, tileCoords.X, tileCoords.Y)

	// Load or get cached tileset
	if tilesetCache == nil {
		var err error
		tilesetCache, err = loadTilesetImage()
		if err != nil {
			fmt.Printf("Error loading tileset: %v\n", err)
			return image.NewRGBA(image.Rect(0, 0, TileSize, TileSize))
		}
		fmt.Printf("Loaded tileset image with bounds: %v\n", tilesetCache.Bounds())
	}

	// Create a new RGBA image for the tile
	tileImg := image.NewRGBA(image.Rect(0, 0, TileSize, TileSize))

	// Calculate the source coordinates in the tileset
	srcX := tileCoords.X * TileSize
	srcY := tileCoords.Y * TileSize

	// Copy pixels directly from the tileset to our tile image
	for y := 0; y < TileSize; y++ {
		for x := 0; x < TileSize; x++ {
			srcPixel := tilesetCache.At(srcX+x, srcY+y)
			tileImg.Set(x, y, srcPixel)
		}
	}

	// Cache the tile
	tileCache[tt] = tileImg

	// Debug: Check if the tile has any non-black pixels
	hasNonBlack := false
	for y := 0; y < TileSize; y++ {
		for x := 0; x < TileSize; x++ {
			r, g, b, a := tileImg.At(x, y).RGBA()
			if r > 0 || g > 0 || b > 0 || a > 0 {
				hasNonBlack = true
				break
			}
		}
		if hasNonBlack {
			break
		}
	}
	fmt.Printf("Tile for terrain type %v has non-black pixels: %v\n", tt, hasNonBlack)

	return tileImg
}

// loadTilesetImage loads the tileset image from the bundled resources
func loadTilesetImage() (image.Image, error) {
	// Get the tileset data from the bundled resource
	tilesetData := resourceAssetsTilesetPng.StaticContent

	// Decode the PNG data into an image
	img, err := png.Decode(bytes.NewReader(tilesetData))
	if err != nil {
		return nil, fmt.Errorf("error decoding tileset: %v", err)
	}

	return img, nil
}
