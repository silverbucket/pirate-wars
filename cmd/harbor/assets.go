package harbor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
)

// Official Media Arts pack hashes (Game Designer signed).
const (
	BackdropSHA256    = "83b588416cd4dc13a0d98b28e75471c1cb1c64becb6d47bd153af70b7b455ed7"
	MaskSHA256        = "e2d79c0c09f6c1c220c6283a5d19edee81ac2d44c94fb9ff04a8c24ba3310bd3"
	PlayerSheetSHA256 = "13f37014ce960486109ca85c4b0c95c11c000c00360bdc33b3d290f7856b9d65"
	NPCSheetSHA256    = "8e3c5d2e18a3b07735db7789b2c29239a5d03a0e5d0c00a0ecae0b53b791e03d"
)

const (
	BackdropPath   = "assets/harbor-midday-1536x1024.png"
	MaskPath       = "assets/harbor-mask-1536x1024.png"
	PlayerShipPath = "assets/ship-player-white-8way-160.png"
	NPCShipPath    = "assets/ship-npc-red-8way-160.png"
)

// ShipSpriteCell is the width/height of each facing cell on the 8-way sheets.
const ShipSpriteCell = 160

// AssetSet holds decoded harbor art after hash verification.
type AssetSet struct {
	Backdrop   image.Image
	Mask       image.Image
	PlayerShip image.Image
	NPCShip    image.Image
}

// VerifyFileSHA256 returns an error when path is missing or hash mismatches.
func VerifyFileSHA256(path, want string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("harbor asset missing %s: %w", path, err)
	}
	sum := sha256.Sum256(data)
	got := hex.EncodeToString(sum[:])
	if got != want {
		return fmt.Errorf("harbor asset hash mismatch %s: got %s want %s", path, got, want)
	}
	return nil
}

// LoadAssets reads and verifies all four harbor PNGs from baseDir (repo root when "").
func LoadAssets(baseDir string) (*AssetSet, error) {
	paths := map[string]string{
		"backdrop":    filepath.Join(baseDir, BackdropPath),
		"mask":        filepath.Join(baseDir, MaskPath),
		"player ship": filepath.Join(baseDir, PlayerShipPath),
		"npc ship":    filepath.Join(baseDir, NPCShipPath),
	}
	hashes := map[string]string{
		"backdrop":    BackdropSHA256,
		"mask":        MaskSHA256,
		"player ship": PlayerSheetSHA256,
		"npc ship":    NPCSheetSHA256,
	}

	set := &AssetSet{}
	var err error

	for name, path := range paths {
		if err = VerifyFileSHA256(path, hashes[name]); err != nil {
			return nil, err
		}
	}

	set.Backdrop, err = loadPNG(paths["backdrop"])
	if err != nil {
		return nil, err
	}
	set.Mask, err = loadPNG(paths["mask"])
	if err != nil {
		return nil, err
	}
	set.PlayerShip, err = loadPNG(paths["player ship"])
	if err != nil {
		return nil, err
	}
	set.NPCShip, err = loadPNG(paths["npc ship"])
	if err != nil {
		return nil, err
	}

	if b := set.Backdrop.Bounds(); b.Dx() != PixelWidth || b.Dy() != PixelHeight {
		return nil, fmt.Errorf("backdrop size %dx%d, want %dx%d", b.Dx(), b.Dy(), PixelWidth, PixelHeight)
	}
	if b := set.Mask.Bounds(); b.Dx() != PixelWidth || b.Dy() != PixelHeight {
		return nil, fmt.Errorf("mask size %dx%d, want %dx%d", b.Dx(), b.Dy(), PixelWidth, PixelHeight)
	}
	pb := set.PlayerShip.Bounds()
	if pb.Dx() != 8*ShipSpriteCell || pb.Dy() != ShipSpriteCell {
		return nil, fmt.Errorf("player ship sheet %dx%d, want %dx%d", pb.Dx(), pb.Dy(), 8*ShipSpriteCell, ShipSpriteCell)
	}
	nb := set.NPCShip.Bounds()
	if nb.Dx() != 8*ShipSpriteCell || nb.Dy() != ShipSpriteCell {
		return nil, fmt.Errorf("npc ship sheet %dx%d, want %dx%d", nb.Dx(), nb.Dy(), 8*ShipSpriteCell, ShipSpriteCell)
	}

	return set, nil
}

func loadPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}
