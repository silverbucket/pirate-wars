package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// loadSplash returns the title image, or nil when the asset is not installed.
func loadSplash() *ebiten.Image {
	f, err := os.Open("./assets/pirate-wars.png")
	if err != nil {
		return nil
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(img)
}
