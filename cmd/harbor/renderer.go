package harbor

import (
	"image"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/window"

	"github.com/hajimehoshi/ebiten/v2"
)

// Renderer draws the painted harbor with Ebiten: the 1536×1024 backdrop is
// blitted as a single image (never sliced, never upscaled from the tileset)
// and the 8-way ship sheets are sub-imaged per facing.
type Renderer struct {
	mask         *Mask
	backdrop     *ebiten.Image
	playerFrames [8]*ebiten.Image
	npcFrames    [8]*ebiten.Image
	viewW        int
	viewH        int
}

// NewRenderer uploads the verified harbor art to the GPU. assets and mask must be non-nil.
func NewRenderer(assets *AssetSet, mask *Mask) *Renderer {
	r := &Renderer{
		mask:     mask,
		backdrop: ebiten.NewImageFromImage(assets.Backdrop),
		viewW:    window.ViewPort.Dimensions.Width,
		viewH:    window.ViewPort.Dimensions.Height,
	}
	sliceShipSheet(assets.PlayerShip, r.playerFrames[:])
	sliceShipSheet(assets.NPCShip, r.npcFrames[:])
	return r
}

// Facings are laid out left to right: N NE E SE S SW W NW.
func sliceShipSheet(sheet image.Image, out []*ebiten.Image) {
	full := ebiten.NewImageFromImage(sheet)
	for i := range out {
		rect := image.Rect(i*ShipSpriteCell, 0, (i+1)*ShipSpriteCell, ShipSpriteCell)
		out[i] = full.SubImage(rect).(*ebiten.Image)
	}
}

// CameraRect returns the top-left harbor pixel for the follow-cam, clamped to the painting.
func CameraRect(playerPos common.Coordinates) (x, y int) {
	px, py, ok := CellCenterPixel(playerPos)
	if !ok {
		return 0, 0
	}
	vw := window.ViewPort.Dimensions.Width
	vh := window.ViewPort.Dimensions.Height
	camX := px - vw/2
	camY := py - vh/2
	if camX < 0 {
		camX = 0
	}
	if camY < 0 {
		camY = 0
	}
	maxX := PixelWidth - vw
	maxY := PixelHeight - vh
	if maxX < 0 {
		maxX = 0
	}
	if maxY < 0 {
		maxY = 0
	}
	if camX > maxX {
		camX = maxX
	}
	if camY > maxY {
		camY = maxY
	}
	return camX, camY
}

// Draw renders the harbor view (backdrop + ships) into dst.
func (r *Renderer) Draw(dst *ebiten.Image, player entities.AvatarReadOnly, npcs []entities.AvatarReadOnly) {
	if r == nil {
		return
	}
	camX, camY := CameraRect(player.GetPos())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(camX), -float64(camY))
	dst.DrawImage(r.backdrop, op)

	for _, n := range npcs {
		if InRegion(n.GetPos()) {
			r.drawShip(dst, n.GetPos(), n.GetFacing(), r.npcFrames[:], camX, camY)
		}
	}
	r.drawShip(dst, player.GetPos(), player.GetFacing(), r.playerFrames[:], camX, camY)
}

func (r *Renderer) drawShip(dst *ebiten.Image, pos common.Coordinates, facing common.Facing, frames []*ebiten.Image, camX, camY int) {
	px, py, ok := CellCenterPixel(pos)
	if !ok {
		return
	}
	idx := int(facing)
	if idx < 0 || idx >= len(frames) || frames[idx] == nil {
		return
	}
	half := ShipSpriteCell / 2
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(px-half-camX), float64(py-half-camY))
	dst.DrawImage(frames[idx], op)
}

// Mask returns the harbor collision mask.
func (r *Renderer) Mask() *Mask {
	if r == nil {
		return nil
	}
	return r.mask
}
