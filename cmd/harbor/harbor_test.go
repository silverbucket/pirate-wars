package harbor

import (
	"image"
	"image/color"
	"testing"

	"pirate-wars/cmd/common"
)

func TestClassifyRGB(t *testing.T) {
	cases := []struct {
		r, g, b uint8
		want    CellKind
	}{
		{47, 112, 193, KindSail},
		{207, 62, 62, KindBlock},
		{76, 177, 60, KindDock},
	}
	for _, tc := range cases {
		got := classifyRGB(tc.r, tc.g, tc.b)
		if got != tc.want {
			t.Fatalf("classifyRGB(%d,%d,%d)=%v want %v", tc.r, tc.g, tc.b, got, tc.want)
		}
	}
}

func TestInRegion(t *testing.T) {
	if !InRegion(Origin) {
		t.Fatal("origin should be in region")
	}
	if InRegion(common.Coordinates{X: Origin.X - 1, Y: Origin.Y}) {
		t.Fatal("west of origin should be out")
	}
	if !InRegion(common.Coordinates{X: Origin.X + WorldCols - 1, Y: Origin.Y + WorldRows - 1}) {
		t.Fatal("SE corner should be in region")
	}
}

func TestMaskPassabilitySynthetic(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, PixelWidth, PixelHeight))
	// Fill sail blue.
	for y := 0; y < PixelHeight; y++ {
		for x := 0; x < PixelWidth; x++ {
			img.Set(x, y, color.RGBA{R: 47, G: 112, B: 193, A: 255})
		}
	}
	// Green dock patch at town pixel center.
	tx, ty, _ := CellCenterPixel(TownPos)
	for dy := -32; dy <= 32; dy++ {
		for dx := -32; dx <= 32; dx++ {
			img.Set(tx+dx, ty+dy, color.RGBA{R: 76, G: 177, B: 60, A: 255})
		}
	}
	m := NewMask(img)
	if !m.IsPassable(TownPos) {
		t.Fatal("town pos should be passable (green)")
	}
	if !m.IsDock(TownPos) {
		t.Fatal("town pos should be dockable (green)")
	}
	// Red block at origin cell center.
	ox, oy, _ := CellCenterPixel(Origin)
	img.Set(ox, oy, color.RGBA{R: 207, G: 62, B: 62, A: 255})
	if m.IsPassable(Origin) {
		t.Fatal("red pixel should block")
	}
}

func TestCameraClamp(t *testing.T) {
	cx, cy := CameraRect(Origin)
	if cx != 0 || cy != 0 {
		t.Fatalf("NW corner cam should be 0,0 got %d,%d", cx, cy)
	}
	se := common.Coordinates{X: Origin.X + WorldCols - 1, Y: Origin.Y + WorldRows - 1}
	cx, cy = CameraRect(se)
	maxX := PixelWidth - 854
	maxY := PixelHeight - 728
	if cx != maxX || cy != maxY {
		t.Fatalf("SE corner cam want %d,%d got %d,%d", maxX, maxY, cx, cy)
	}
}

// TestNilWorldIsNilSafe covers the missing-asset path: LoadAssets fails, the
// harbor world stays a nil *World, and every method must tolerate that receiver.
func TestNilWorldIsNilSafe(t *testing.T) {
	var w *World
	if pos, ok := w.ClampSpawn(TownPos); ok {
		t.Fatalf("nil world returned spawn %+v", pos)
	}
	if w.IsPassableByBoat(TownPos) {
		t.Fatal("nil world should not be passable")
	}
	if w.IsDock(TownPos) {
		t.Fatal("nil world should not be a dock")
	}
	if !w.InHarbor(TownPos) {
		t.Fatal("InHarbor is a pure region check and should still work")
	}
}

// TestClampSpawnWithoutMask covers a constructed world whose mask never loaded.
func TestClampSpawnWithoutMask(t *testing.T) {
	w := NewWorld(nil)
	if _, ok := w.ClampSpawn(TownPos); ok {
		t.Fatal("world without a mask should not return a spawn")
	}
}

func TestOfficialAssetsIfPresent(t *testing.T) {
	set, err := LoadAssets("")
	if err != nil {
		t.Skip("official harbor assets not on disk:", err)
	}
	if set.Backdrop == nil || set.Mask == nil {
		t.Fatal("expected decoded images")
	}
}
