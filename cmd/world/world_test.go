package world

import (
	"image"
	"image/color"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/window"
	"testing"

	"go.uber.org/zap"
)

func initTestLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

type AvatarMock struct {
	pos  common.Coordinates
	char rune
}

func (av AvatarMock) GetPos() common.Coordinates          { return av.pos }
func (av AvatarMock) GetPreviousPos() common.Coordinates  { return av.pos }
func (av AvatarMock) GetID() string                       { return "" }
func (av AvatarMock) GetFlag() string                     { return "" }
func (av AvatarMock) GetType() string                     { return "" }
func (av AvatarMock) GetName() string                     { return "" }
func (av AvatarMock) GetColor() color.Color               { return color.White }
func (av AvatarMock) GetFacing() common.Facing            { return common.FacingN }
func (av AvatarMock) MovedThisTick() bool                 { return false }
func (av AvatarMock) GetViewableRange() window.Dimensions { return window.Dimensions{} }
func (av AvatarMock) GetCharacter() string                { return string(av.char) }
func (av AvatarMock) GetTileImage() image.Image {
	return image.NewRGBA(image.Rect(0, 0, window.CellSize, window.CellSize))
}

func (av AvatarMock) IsHighlighted() bool { return false }
func (av AvatarMock) Highlight(b bool)    {}

func TestWorldInit(t *testing.T) {
	c := common.Coordinates{X: 10, Y: 10}
	logger := initTestLogger()
	world := Init(logger)
	world.SetPositionType(c, 99)
	tt := world.GetPositionType(c)
	if tt != 99 {
		t.Fatalf("SetPositionType not set")
	}
}

// TestMinimapImage exercises the headless minimap path (no GPU context needed).
func TestMinimapImage(t *testing.T) {
	avatar := AvatarMock{pos: common.Coordinates{X: 100, Y: 100}, char: '@'}
	logger := initTestLogger()
	world := Init(logger)
	img := world.MinimapImage(avatar.GetPos(), entities.ViewableEntities{})
	if img == nil {
		t.Fatal("expected minimap image")
	}
	if img.Bounds().Dx() != window.MiniMapArea.Width {
		t.Fatalf("minimap width = %d, want %d", img.Bounds().Dx(), window.MiniMapArea.Width)
	}
}
