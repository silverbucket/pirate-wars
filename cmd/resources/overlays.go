package resources

import (
	"image"
	"image/color"
	"image/draw"
)

// GetExamineRingOverlay returns the examine ring tile, or a code-drawn fallback.
func GetExamineRingOverlay(size int) image.Image {
	if tile := getOverlayTile(ExamineRingCol, ExamineRingRow); tile != nil {
		return tile
	}
	return GetHighlightOverlay(size)
}

// GetPlayerMarkerOverlay returns the player marker tile, or a code-drawn fallback.
func GetPlayerMarkerOverlay(size int) image.Image {
	if tile := getOverlayTile(PlayerMarkerCol, PlayerMarkerRow); tile != nil {
		return tile
	}
	return codeDrawnPlayerMarker(size)
}

func getOverlayTile(col, row int) image.Image {
	if !HasExpandedTileset() {
		return nil
	}
	tile := extractTileAt(col, row)
	if tile == nil || isTileNearlyEmpty(tile) {
		return nil
	}
	return tile
}

func codeDrawnPlayerMarker(size int) image.Image {
	if cached, ok := playerMarkerCache[size]; ok {
		return cached
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	marker := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	inset := size / 8
	if inset < 2 {
		inset = 2
	}
	length := size / 4
	if length < 4 {
		length = 4
	}

	for i := 0; i < length; i++ {
		img.Set(inset+i, inset, marker)
		img.Set(inset, inset+i, marker)

		img.Set(size-inset-1-i, inset, marker)
		img.Set(size-inset-1, inset+i, marker)

		img.Set(inset+i, size-inset-1, marker)
		img.Set(inset, size-inset-1-i, marker)

		img.Set(size-inset-1-i, size-inset-1, marker)
		img.Set(size-inset-1, size-inset-1-i, marker)
	}

	playerMarkerCache[size] = img
	return img
}

// CompositeOverlays alpha-composites overlay images onto a transparent base.
func CompositeOverlays(size int, overlays ...image.Image) image.Image {
	result := image.NewRGBA(image.Rect(0, 0, size, size))
	for _, overlay := range overlays {
		if overlay == nil {
			continue
		}
		draw.Draw(result, result.Bounds(), overlay, image.Point{}, draw.Over)
	}
	return result
}

// HasTransparentCenter reports whether the center of an overlay is mostly transparent.
func HasTransparentCenter(img image.Image) bool {
	if img == nil {
		return false
	}
	bounds := img.Bounds()
	cx := bounds.Min.X + bounds.Dx()/2
	cy := bounds.Min.Y + bounds.Dy()/2
	radius := bounds.Dx() / 8
	if radius < 2 {
		radius = 2
	}

	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0x8000 {
				return false
			}
		}
	}
	return true
}
