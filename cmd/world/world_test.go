package world

import (
	"image/color"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/window"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"go.uber.org/zap"
)

var testApp fyne.App

func initTestLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

func setup() {
	testApp = app.New()
}

func cleanup() {
	if testApp != nil {
		testApp.Quit()
	}
}

type AvatarMock struct {
	pos  common.Coordinates
	char rune
}

func (av AvatarMock) GetPos() common.Coordinates          { return av.pos }
func (av AvatarMock) GetPreviousPos() common.Coordinates  { return av.pos }
func (av AvatarMock) GetID() string                       { return "" }
func (av AvatarMock) Highlight()                          {}
func (av AvatarMock) GetFlag() string                     { return "" }
func (av AvatarMock) GetType() string                     { return "" }
func (av AvatarMock) GetName() string                     { return "" }
func (av AvatarMock) GetForegroundColor() color.Color     { return color.White }
func (av AvatarMock) GetBackgroundColor() color.Color     { return color.White }
func (av AvatarMock) GetViewableRange() window.Dimensions { return window.Dimensions{} }
func (av AvatarMock) GetCharacter() string                { return string(av.char) }

func TestWorldInit(t *testing.T) {
	setup()
	t.Cleanup(cleanup)
	c := common.Coordinates{X: 10, Y: 10}
	logger := initTestLogger()
	world := Init(logger)
	world.SetPositionType(c, 99)
	tt := world.GetPositionType(c)
	if tt != 99 {
		t.Fatalf("SetPositionType not set")
	}
}

func TestPaint(t *testing.T) {
	setup()
	t.Cleanup(cleanup)
	avatar := AvatarMock{pos: common.Coordinates{X: 100, Y: 100}, char: '@'}
	logger := initTestLogger()
	world := Init(logger)
	world.Paint(avatar, []entities.AvatarReadOnly{}, avatar)
}
