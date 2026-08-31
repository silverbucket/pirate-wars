package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/sailing"
)

// playpassCfg loads the branch sailing.cfg defaults (do not mutate).
func playpassCfg(t *testing.T) sailing.Config {
	t.Helper()
	cfg := sailing.LoadConfig("sailing.cfg")
	if cfg.WindDriftTicks != 100 || cfg.TickMS != 250 {
		t.Fatalf("unexpected cfg: drift=%d tick_ms=%d", cfg.WindDriftTicks, cfg.TickMS)
	}
	return cfg
}

func alwaysPassable(common.Coordinates) bool { return true }

type tickSim struct {
	cfg       sailing.Config
	wind      *sailing.Wind
	avatar    entities.Avatar
	occupancy sailing.Occupancy
}

func newTickSim(cfg sailing.Config, wind *sailing.Wind, pos common.Coordinates, heading common.Facing) *tickSim {
	a := entities.CreateAvatar(pos, common.ShipWhite, entities.ColorPossibilities[0])
	a.SetHeading(heading)
	a.SetSail(sailing.SailFull)
	ships := map[string]common.Coordinates{a.GetID(): pos}
	return &tickSim{
		cfg:       cfg,
		wind:      wind,
		avatar:    a,
		occupancy: sailing.NewOccupancy(ships),
	}
}

// playerTick mirrors resolvePlayerMovement in sailing_tick.go.
func (s *tickSim) playerTick() (speed float64, stepped bool) {
	speed = s.cfg.EffectiveSpeed(s.avatar.GetFacing(), s.avatar.GetSail(), s.wind)
	s.avatar.SetLastSpeed(speed)
	if !s.avatar.AccumulateSpeed(speed) {
		return speed, false
	}
	newPos, moved := sailing.TryStep(
		s.avatar.GetPos(),
		s.avatar.GetFacing(),
		s.avatar.GetID(),
		s.occupancy,
		alwaysPassable,
	)
	if moved {
		oldKey := common.CoordToKey(s.avatar.GetPos())
		s.avatar.SetPos(newPos)
		delete(s.occupancy, oldKey)
		s.occupancy[common.CoordToKey(newPos)] = s.avatar.GetID()
	}
	return speed, moved
}

func parseHUDSpeed(text string) (float64, bool) {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "Speed:") {
			v, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimPrefix(line, "Speed:")), 64)
			return v, err == nil
		}
	}
	return 0, false
}

func parseHUDWind(text string) (label string, strength int, ok bool) {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "Wind:") {
			rest := strings.TrimSpace(strings.TrimPrefix(line, "Wind:"))
			open := strings.Index(rest, "(")
			close := strings.Index(rest, ")")
			if open < 0 || close < open {
				return "", 0, false
			}
			label = strings.TrimSpace(rest[:open])
			s, err := strconv.Atoi(strings.TrimSpace(rest[open+1 : close]))
			return label, s, err == nil
		}
	}
	return "", 0, false
}

// TestPlayPassSailingLoop drives the sailing tick path headlessly and locks designer Q1–Q6.
func TestPlayPassSailingLoop(t *testing.T) {
	cfg := playpassCfg(t)
	wind := sailing.NewFixedWind(cfg, common.FacingE, 2) // fresh, downwind east

	// Q1: make way with no keys — full sail, heading set, wind blowing.
	runSim := newTickSim(cfg, wind, common.Coordinates{X: 20, Y: 20}, common.FacingE)
	stepsQ1 := 0
	var lastSpeedQ1 float64
	for i := 0; i < 30; i++ {
		spd, stepped := runSim.playerTick()
		lastSpeedQ1 = spd
		if stepped {
			stepsQ1++
		}
	}
	t.Logf("Q1 make way (no keys): tick_speed=%.3f steps_in_30_ticks=%d final_pos=%+v",
		lastSpeedQ1, stepsQ1, runSim.avatar.GetPos())
	if stepsQ1 == 0 {
		t.Fatal("Q1: ship did not step with full sail + run heading + wind")
	}

	// Q2: pointing upwind crawl — irons vs run.
	ironsSim := newTickSim(cfg, wind, common.Coordinates{X: 20, Y: 20}, common.FacingW)
	runSim2 := newTickSim(cfg, wind, common.Coordinates{X: 20, Y: 20}, common.FacingE)
	ironsSteps, runSteps := 0, 0
	var ironsSpeed, runSpeed float64
	for i := 0; i < 200; i++ {
		spd, stepped := ironsSim.playerTick()
		ironsSpeed = spd
		if stepped {
			ironsSteps++
		}
		spd, stepped = runSim2.playerTick()
		runSpeed = spd
		if stepped {
			runSteps++
		}
	}
	t.Logf("Q2 upwind crawl: irons_speed=%.4f irons_steps/200=%d run_speed=%.4f run_steps/200=%d ratio=%.1fx",
		ironsSpeed, ironsSteps, runSpeed, runSteps, runSpeed/ironsSpeed)
	if ironsSpeed >= runSpeed {
		t.Fatalf("Q2: irons speed %.4f should be < run %.4f", ironsSpeed, runSpeed)
	}
	if runSteps <= ironsSteps {
		t.Fatalf("Q2: run should outpace irons over 200 ticks (%d vs %d steps)", runSteps, ironsSteps)
	}

	// Q3: falling off speeds up — close → beam → run.
	headings := []struct {
		name string
		h    common.Facing
	}{
		{"close", common.FacingSW},
		{"beam", common.FacingS},
		{"run", common.FacingE},
	}
	var speedsQ3 []float64
	for _, h := range headings {
		sim := newTickSim(cfg, wind, common.Coordinates{X: 20, Y: 20}, h.h)
		spd, _ := sim.playerTick()
		speedsQ3 = append(speedsQ3, spd)
		t.Logf("Q3 %s heading speed=%.4f", h.name, spd)
	}
	if !(speedsQ3[0] < speedsQ3[1] && speedsQ3[1] < speedsQ3[2]) {
		t.Fatalf("Q3: falling off should increase speed: close=%.4f beam=%.4f run=%.4f",
			speedsQ3[0], speedsQ3[1], speedsQ3[2])
	}

	// Q4: two ships must not stack — soft bump.
	playerPos := common.Coordinates{X: 10, Y: 10}
	npcPos := common.Coordinates{X: 11, Y: 10}
	player := entities.CreateAvatar(playerPos, common.ShipWhite, entities.ColorPossibilities[0])
	player.SetHeading(common.FacingE)
	npc := entities.CreateAvatar(npcPos, common.ShipRed, entities.ColorPossibilities[1])
	occupancy := sailing.NewOccupancy(map[string]common.Coordinates{
		player.GetID(): playerPos,
		npc.GetID():    npcPos,
	})
	blockedPos, moved := sailing.TryStep(playerPos, common.FacingE, player.GetID(), occupancy, alwaysPassable)
	stacked := common.CoordsMatch(blockedPos, npcPos)
	t.Logf("Q4 collision: player@%+v npc@%+v try_step_east moved=%v final=%+v stacked=%v",
		playerPos, npcPos, moved, blockedPos, stacked)
	if moved || stacked {
		t.Fatal("Q4: player must soft-bump, not occupy NPC cell")
	}

	// Q5: HUD Speed/Wind match live sailing state (not placeholders like speed=5).
	hudSim := newTickSim(cfg, wind, common.Coordinates{X: 20, Y: 20}, common.FacingE)
	liveSpeed, _ := hudSim.playerTick()
	hudPoint := sailing.ClassifyPointOfSail(hudSim.avatar.GetFacing(), hudSim.wind.Facing)
	hudText := shipStatusText(shipStatus{
		Heading:       hudSim.avatar.GetFacing(),
		Sail:          hudSim.avatar.GetSail(),
		Wind:          hudSim.wind,
		Point:         hudPoint,
		PointMult:     cfg.PointOfSailMultiplier(hudPoint),
		Speed:         hudSim.avatar.GetLastSpeed(),
		TimeOfDay:     "06:00",
		Gold:          50,
		CargoCapacity: 20,
	})
	hudSpeed, ok := parseHUDSpeed(hudText)
	if !ok {
		t.Fatal("Q5: could not parse HUD Speed")
	}
	hudLabel, hudStrength, ok := parseHUDWind(hudText)
	if !ok {
		t.Fatal("Q5: could not parse HUD Wind")
	}
	t.Logf("Q5 HUD: live_speed=%.3f hud_speed=%.3f wind=%s(%d) hud_wind=%s(%d)",
		liveSpeed, hudSpeed, wind.Label(), wind.Strength, hudLabel, hudStrength)
	if math.Abs(hudSpeed-liveSpeed) > 1e-9 {
		t.Fatalf("Q5: HUD speed %.2f != live %.3f", hudSpeed, liveSpeed)
	}
	if hudLabel != wind.Label() || hudStrength != wind.Strength {
		t.Fatalf("Q5: HUD wind %s(%d) != live %s(%d)", hudLabel, hudStrength, wind.Label(), wind.Strength)
	}
	if strings.Contains(hudText, "Speed: 5.0") && math.Abs(liveSpeed-5.0) > 0.05 {
		t.Fatal("Q5: HUD shows placeholder speed 5.0")
	}

	// Q6: wind shifts once on a long reach (100 ticks @ wind_drift_ticks=100).
	rand.Seed(42)
	driftWind := sailing.NewWind(cfg)
	startFacing := driftWind.Facing
	startStrength := driftWind.Strength
	facings := []string{fmt.Sprintf("tick0:%s", driftWind.Label())}
	for i := 1; i <= 120; i++ {
		driftWind.Tick()
		if i == 100 || i == 120 {
			facings = append(facings, fmt.Sprintf("tick%d:%s", i, driftWind.Label()))
		}
	}
	shiftedAt100 := driftWind.Facing != startFacing || driftWind.Strength != startStrength
	t.Logf("Q6 wind drift: start=%s str=%d path=%v shifted_by_tick100=%v end=%s str=%d",
		common.FacingLabel(startFacing), startStrength, facings, shiftedAt100,
		driftWind.Label(), driftWind.Strength)
	if !shiftedAt100 {
		t.Fatal("Q6: wind did not shift after 100 ticks")
	}
}
