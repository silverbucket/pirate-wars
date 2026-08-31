package resources

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
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

const shipBaseRow = 2

var (
	tilesetCache          image.Image
	tileCache             = make(map[int]image.Image)
	shipTileCache         = make(map[int]image.Image)
	regionTileCache       = make(map[image.Point]image.Image)
	highlightOverlayCache = make(map[int]image.Image)
	playerMarkerCache     = make(map[int]image.Image)
)

func getTileByRegion(idx int) image.Image {
	// Get the tile coordinates from the mapping
	tileCoords, ok := TileMapping[idx]
	if !ok {
		tileCoords = TileMapping[common.ShipWhite]
	}

	// Check if we have this tile cached
	if cached, ok := tileCache[idx]; ok {
		return cached
	}

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

func shipColumn(s common.ShipType) int {
	tileCoords, ok := TileMapping[int(s)]
	if !ok {
		tileCoords = TileMapping[common.ShipWhite]
	}
	return tileCoords.X
}

func shipTileCacheKey(s common.ShipType, facing common.Facing) int {
	return int(s)*100 + int(facing)
}

func tileRegionInBounds(col, row int) bool {
	tileset := getTileset()
	bounds := tileset.Bounds()
	srcX := col * TileSize
	srcY := row * TileSize
	return srcX+TileSize <= bounds.Dx() && srcY+TileSize <= bounds.Dy() && srcX >= 0 && srcY >= 0
}

// extractTileAt memoizes tiles so callers can rely on stable image identity
// (the Ebiten texture cache keys on it).
func extractTileAt(col, row int) image.Image {
	if !tileRegionInBounds(col, row) {
		return nil
	}

	// A coordinate key rather than arithmetic: col*1000+row carried a silent
	// "fewer than 1000 rows" assumption, and packing into an int drops col
	// entirely where int is 32 bits.
	regionKey := image.Point{X: col, Y: row}
	if cached, ok := regionTileCache[regionKey]; ok {
		return cached
	}

	tileset := getTileset()
	tileImg := image.NewRGBA(image.Rect(0, 0, TileSize, TileSize))
	srcX := col * TileSize
	srcY := row * TileSize

	for y := 0; y < TileSize; y++ {
		for x := 0; x < TileSize; x++ {
			tileImg.Set(x, y, tileset.At(srcX+x, srcY+y))
		}
	}

	regionTileCache[regionKey] = tileImg
	return tileImg
}

func isTileNearlyEmpty(img image.Image) bool {
	if img == nil {
		return true
	}

	coloredPixels := 0
	for y := 0; y < TileSize; y++ {
		for x := 0; x < TileSize; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a < 0x1000 {
				continue
			}
			if r>>8 > 20 || g>>8 > 20 || b>>8 > 20 {
				coloredPixels++
			}
		}
	}

	return coloredPixels < 10
}

func GetShipTile(s common.ShipType, facing common.Facing) image.Image {
	cacheKey := shipTileCacheKey(s, facing)
	if cached, ok := shipTileCache[cacheKey]; ok {
		return cached
	}

	col := shipColumn(s)
	row := shipBaseRow + int(facing)
	tile := extractTileAt(col, row)

	if facing != common.FacingN && (tile == nil || isTileNearlyEmpty(tile)) {
		tile = GetShipTile(s, common.FacingN)
	}

	shipTileCache[cacheKey] = tile
	return tile
}

// GetTerrainTile returns the image for a specific terrain type
func GetTerrainTile(tt common.TerrainType) image.Image {
	idx := int(tt)
	return getTileByRegion(idx)
}

// loadTilesetImage loads the tileset image from the bundled resources
func loadTilesetImage() (image.Image, error) {
	// Get the tileset data from the bundled resource
	tilesetData := tilesetPNGData

	// Decode the PNG data into an image
	img, err := png.Decode(bytes.NewReader(tilesetData))
	if err != nil {
		return nil, fmt.Errorf("error decoding tileset: %v", err)
	}

	return img, nil
}

// GetHighlightOverlay returns a pulsing highlight frame for examined entities.
func GetHighlightOverlay(size int) image.Image {
	if cached, ok := highlightOverlayCache[size]; ok {
		return cached
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for x := 0; x < size; x++ {
		img.Set(x, 0, image.White)
		img.Set(x, size-1, image.White)
	}
	for y := 0; y < size; y++ {
		img.Set(0, y, image.White)
		img.Set(size-1, y, image.White)
	}

	highlightOverlayCache[size] = img
	return img
}

// CompositeWithHighlight draws a highlight frame over a base tile image.
func CompositeWithHighlight(base image.Image, size int) image.Image {
	result := image.NewRGBA(image.Rect(0, 0, size, size))
	if base != nil {
		draw.Draw(result, result.Bounds(), base, image.Point{}, draw.Over)
	}
	draw.Draw(result, result.Bounds(), GetHighlightOverlay(size), image.Point{}, draw.Over)
	return result
}

func getTileset() image.Image {
	if tilesetCache == nil {
		var err error
		tilesetCache, err = loadTilesetImage()
		if err != nil {
			return image.NewRGBA(image.Rect(0, 0, TileSize, TileSize))
		}
	}
	return tilesetCache
}
