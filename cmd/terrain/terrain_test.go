package terrain

import (
	"go.uber.org/zap"
	"testing"
)

func initTestLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

func cleanup() {
}

type CoordinatesMock struct {
	X int
	Y int
}

type AvatarMock struct {
	pos  CoordinatesMock
	char rune
}

func (av AvatarMock) Render() string {
	return ""
}
func (av AvatarMock) GetX() int        { return av.pos.X }
func (av AvatarMock) GetMiniMapX() int { return av.pos.X }
func (av AvatarMock) GetY() int        { return av.pos.X }
func (av AvatarMock) GetMiniMapY() int { return av.pos.X }

func TestPaint(t *testing.T) {
	t.Cleanup(cleanup)
	avatar := AvatarMock{pos: CoordinatesMock{X: 100, Y: 100}, char: '@'}
	logger := initTestLogger()
	tr := Init(logger)
	tr.GenerateWorld()
	tr.GenerateTowns()
	tr.World.Paint(avatar, []AvatarReadOnly{})
}
