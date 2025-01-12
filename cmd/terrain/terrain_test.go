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
func (av AvatarMock) GetPos() CoordinatesMock { return av.pos }

func TestPaint(t *testing.T) {
	t.Cleanup(cleanup)
	avatar := AvatarMock{pos: CoordinatesMock{X: 100, Y: 100}, char: '@'}
	logger := initTestLogger()
	tr := Init(logger)
	tr.GenerateWorld()
	tr.GenerateTowns()
	tr.World.Paint(avatar, []AvatarMock{}, avatar)
}
