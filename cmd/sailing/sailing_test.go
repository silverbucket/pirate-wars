package sailing

import (
	"math"
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
	if cfg.HullSpeed != 0.55 {
		t.Fatalf("hull_speed = %v, want 0.55", cfg.HullSpeed)
	}
	if cfg.WindDriftTicks != 100 {
		t.Fatalf("wind_drift_ticks = %d, want 100", cfg.WindDriftTicks)
	}
	if cfg.NPCSkipPercent != 0 {
		t.Fatalf("npc_skip_percent = %d, want 0", cfg.NPCSkipPercent)
	}
	if cfg.WindFactorLight != 0.55 || cfg.WindFactorFresh != 1.0 || cfg.WindFactorStrong != 1.25 {
		t.Fatalf("wind factors = %v/%v/%v, want 0.55/1.0/1.25",
			cfg.WindFactorLight, cfg.WindFactorFresh, cfg.WindFactorStrong)
	}
}

func TestWindStrengthFactorTable(t *testing.T) {
	cfg := DefaultConfig()
	cases := []struct {
		strength int
		want     float64
	}{
		{0, 0},
		{1, 0.55},
		{2, 1.0},
		{3, 1.25},
	}
	for _, tc := range cases {
		got := cfg.WindStrengthFactor(tc.strength)
		if got != tc.want {
			t.Fatalf("strength %d factor = %v, want %v", tc.strength, got, tc.want)
		}
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sailing.cfg")
	content := "tick_ms=100\nhull_speed=2.5\npoint_of_sail_run=0.8\nnpc_skip_percent=25\nwind_factor_fresh=0.9\nwind_drift_ticks=90\n"
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
	if cfg.WindFactorFresh != 0.9 {
		t.Fatalf("wind_factor_fresh = %v, want 0.9", cfg.WindFactorFresh)
	}
	if cfg.WindDriftTicks != 90 {
		t.Fatalf("wind_drift_ticks = %d, want 90", cfg.WindDriftTicks)
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

func TestEffectiveSpeedUsesWindFactorTable(t *testing.T) {
	cfg := DefaultConfig()
	wind := &Wind{cfg: cfg, Facing: common.FacingE, Strength: 2}
	speed := cfg.EffectiveSpeed(common.FacingE, SailFull, wind)
	want := cfg.HullSpeed * cfg.PointOfSailRun * cfg.WindFactorFresh * cfg.SailFull
	if math.Abs(speed-want) > 1e-9 {
		t.Fatalf("effective speed = %v, want %v", speed, want)
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

func TestAccumulateSpeed(t *testing.T) {
	var progress float64
	if AccumulateSpeed(&progress, 0) {
		t.Fatal("zero tick speed should not advance")
	}
	if AccumulateSpeed(&progress, 0.3) {
		t.Fatal("0.3 should not trigger step yet")
	}
	if progress != 0.3 {
		t.Fatalf("progress = %v, want 0.3", progress)
	}
	if !AccumulateSpeed(&progress, 0.8) {
		t.Fatal("0.3 + 0.8 should trigger step")
	}
	if math.Abs(progress-0.1) > 1e-9 {
		t.Fatalf("carry progress = %v, want 0.1", progress)
	}
}

func TestRunFreshWindDoesNotStepEveryTick(t *testing.T) {
	cfg := DefaultConfig()
	wind := &Wind{cfg: cfg, Facing: common.FacingE, Strength: 2}
	tickSpeed := cfg.EffectiveSpeed(common.FacingE, SailFull, wind)
	if tickSpeed >= 1.0 {
		t.Fatalf("run+fresh tick speed = %v, want < 1 so accumulation gates steps", tickSpeed)
	}
}
