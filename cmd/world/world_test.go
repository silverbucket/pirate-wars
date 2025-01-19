package world

import (
	"go.uber.org/zap"
	"pirate-wars/cmd/common"
	"testing"
)

func initTestLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

func cleanup() {
}

type AvatarMock struct {
	pos  common.Coordinates
	char rune
}

func (av AvatarMock) Render() string {
	return ""
}
func (av AvatarMock) GetPos() common.Coordinates { return av.pos }
func (av AvatarMock) GetID() string              { return "" }
func (av AvatarMock) Highlight()                 {}
func (av AvatarMock) GetFlag() string            { return "" }
func (av AvatarMock) GetType() string            { return "" }
func (av AvatarMock) GetName() string            { return "" }
func (av AvatarMock) GetForegroundColor() string { return "" }

func TestWorldInit(t *testing.T) {
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
	t.Cleanup(cleanup)
	avatar := AvatarMock{pos: common.Coordinates{X: 100, Y: 100}, char: '@'}
	logger := initTestLogger()
	world := Init(logger)
	world.Paint(avatar, []common.AvatarReadOnly{}, avatar, 1)
}
