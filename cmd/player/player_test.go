package player

import (
	"testing"

	"pirate-wars/cmd/harbor"
	"pirate-wars/cmd/world"

	"go.uber.org/zap"
)

func TestCreateWithoutHarborWorld(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mv := world.Init(logger.Sugar())
	t.Cleanup(func() {
		if r := recover(); r != nil {
			t.Fatalf("Create panicked with nil harbor world: %v", r)
		}
	})
	p := Create(mv, nil)
	if p == nil {
		t.Fatal("expected player avatar")
	}
}

func TestCreateWithNilHarborWorldAfterFailedLoad(t *testing.T) {
	if _, err := harbor.LoadAssets(""); err == nil {
		t.Skip("harbor assets present; test targets missing-asset fallback")
	}
	logger, _ := zap.NewDevelopment()
	mv := world.Init(logger.Sugar())
	var hw *harbor.World
	t.Cleanup(func() {
		if r := recover(); r != nil {
			t.Fatalf("Create panicked with typed-nil harbor path: %v", r)
		}
	})
	p := Create(mv, hw)
	if p == nil {
		t.Fatal("expected player avatar")
	}
}
