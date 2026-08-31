package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
)

// The window icon is rasterised from assets/icon.svg by `make icons`; the
// PNGs are committed so a plain `go build` needs no SVG toolchain.
var (
	//go:embed assets/icon-32.png
	iconPNG32 []byte
	//go:embed assets/icon-64.png
	iconPNG64 []byte
	//go:embed assets/icon-256.png
	iconPNG256 []byte
)

// windowIcons decodes the embedded skull-and-bones at each size the window
// manager might ask for. A PNG that fails to decode is skipped rather than
// fatal: a missing icon is cosmetic.
func windowIcons() []image.Image {
	var icons []image.Image
	for _, data := range [][]byte{iconPNG32, iconPNG64, iconPNG256} {
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			continue
		}
		icons = append(icons, img)
	}
	return icons
}
