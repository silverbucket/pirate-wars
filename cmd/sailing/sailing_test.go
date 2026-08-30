package sailing

import (
	"os"
	"path/filepath"
	"pirate-wars/cmd/common"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TickMS != 250 {
		t.Fatalf("tick_ms = %d, want 250", cfg.TickMS)
	}
	if cfg.HullSpeed != 1.0 {
		t.Fatalf("hull_speed = %v, want 1.0", cfg.HullSpeed)
	}
	if cfg.NPCSkipPercent != 50 {
		t.Fatalf("npc_skip_percent = %d, want 50", cfg.NPCSkipPercent)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sailing.cfg")
	content := "tick_ms=100\nhull_speed=2.5\npoint_of_sail_run=0.8\nnpc_skip_percent=25\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := LoadConfig(path)
	if cfg.TickMS != 100 {
		t.Fatalf("tick_ms = %d, want 100", cfg.TickMS)
	}
	if cfg.HullSpeed != 2.5 {
		t.Fatalf("hull_speed = %v, want 2.5", cfg.HullSpeed)
	}
	if cfg.PointOfSailRun != 0.8 {
		t.Fatalf("point_of_sail_run = %v, want 0.8", cfg.PointOfSailRun)
	}
	if cfg.NPCSkipPercent != 25 {
		t.Fatalf("npc_skip_percent = %d, want 25", cfg.NPCSkipPercent)
	}
}

func TestLoadConfigMissingFileUsesDefaults(t *testing.T) {
	cfg := LoadConfig(filepath.Join(t.TempDir(), "missing.cfg"))
	defaults := DefaultConfig()
	if cfg.TickMS != defaults.TickMS {
		t.Fatalf("expected default tick_ms")
	}
}

func TestPointOfSailMultipliers(t *testing.T) {
	cfg := DefaultConfig()
	cases := []struct {
		heading, wind common.Facing
		want          PointOfSail
		mult          float64
	}{
		{common.FacingE, common.FacingE, PointRun, cfg.PointOfSailRun},
		{common.FacingE, common.FacingSE, PointBroad, cfg.PointOfSailBroad},
		{common.FacingE, common.FacingS, PointBeam, cfg.PointOfSailBeam},
		{common.FacingE, common.FacingSW, PointClose, cfg.PointOfSailClose},
		{common.FacingE, common.FacingW, PointIrons, cfg.PointOfSailIrons},
	}
	for _, tc := range cases {
		got := ClassifyPointOfSail(tc.heading, tc.wind)
		if got != tc.want {
			t.Fatalf("heading %v wind %v = %v, want %v", tc.heading, tc.wind, got, tc.want)
		}
		if m := cfg.PointOfSailMultiplier(got); m != tc.mult {
			t.Fatalf("multiplier for %v = %v, want %v", got, m, tc.mult)
		}
	}
}

func TestWindDriftBounds(t *testing.T) {
	cfg := DefaultConfig()
	cfg.WindDriftTicks = 1
	cfg.WindStrengthMin = 1
	cfg.WindStrengthMax = 3
	w := NewWind(cfg)

	for i := 0; i < 200; i++ {
		w.Tick()
		if w.Strength < cfg.WindStrengthMin || w.Strength > cfg.WindStrengthMax {
			t.Fatalf("strength %d out of bounds [%d,%d]", w.Strength, cfg.WindStrengthMin, cfg.WindStrengthMax)
		}
		if int(w.Facing) < 0 || int(w.Facing) > 7 {
			t.Fatalf("facing %d out of 8-way range", w.Facing)
		}
	}
}

func TestCollisionSoftBump(t *testing.T) {
	posA := common.Coordinates{X: 5, Y: 5}
	posB := common.Coordinates{X: 6, Y: 5}
	occupancy := NewOccupancy(map[string]common.Coordinates{
		"a": posA,
		"b": posB,
	})

	final, moved := TryStep(posA, common.FacingE, "a", occupancy, func(c common.Coordinates) bool { return true })
	if moved {
		t.Fatal("should not move into occupied cell")
	}
	if !common.CoordsMatch(final, posA) {
		t.Fatalf("bumped ship should stay at %v, got %v", posA, final)
	}

	free := common.Coordinates{X: 8, Y: 5}
	final, moved = TryStep(free, common.FacingW, "c", occupancy, func(c common.Coordinates) bool { return true })
	want := common.Coordinates{X: 7, Y: 5}
	if !moved || !common.CoordsMatch(final, want) {
		t.Fatalf("open cell move = %v moved=%v, want %v true", final, moved, want)
	}
}

func TestWakeAftPosition(t *testing.T) {
	pos := common.Coordinates{X: 10, Y: 10}
	aft := WakeAftPosition(pos, common.FacingE)
	want := common.Coordinates{X: 9, Y: 10}
	if !common.CoordsMatch(aft, want) {
		t.Fatalf("aft of east heading = %v, want %v", aft, want)
	}
}

func TestShouldMoveThresholds(t *testing.T) {
	if ShouldMove(0) {
		t.Fatal("zero speed should not move")
	}
	if !ShouldMove(1.5) {
		t.Fatal("speed >= 1 should always move")
	}
}
