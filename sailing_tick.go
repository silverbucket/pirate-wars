package main

import (
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/hail"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/sailing"
)

// resolveSailingTick advances wind, then the player, then every NPC.
func (m *GameState) resolveSailingTick() {
	m.wind.Tick()

	// Steering is rationed per tick, not per frame. Key presses are read at the
	// render rate, so without this a ship could be spun through all eight octants
	// between two ticks, for free — which makes the wind model decorative,
	// because the best heading is always one keypress away at no cost.
	m.turnedThisTick = false

	m.player.ClearMovedFlag()
	m.npcs.ClearMovedFlags()

	occupancy := m.buildOccupancy()

	m.resolvePlayerMovement(occupancy)
	m.npcs.ResolveMovements(m.sailingCfg, m.economyCfg, m.wind, m.world, m.towns, occupancy)
}

// buildOccupancy indexes every ship by cell so a tick can detect collisions.
func (m *GameState) buildOccupancy() sailing.Occupancy {
	ships := make(map[string]common.Coordinates)
	ships[m.player.GetID()] = m.player.GetPos()
	for _, n := range m.npcs.GetList() {
		ships[n.GetID()] = n.GetPos()
	}
	return sailing.NewOccupancy(ships)
}

// resolvePlayerMovement accumulates speed and steps the player one cell when due.
func (m *GameState) resolvePlayerMovement(occupancy sailing.Occupancy) {
	speed := m.sailingCfg.EffectiveSpeed(m.player.GetFacing(), m.player.GetSail(), m.wind)
	m.player.SetLastSpeed(speed)
	if !m.player.AccumulateSpeed(speed) {
		return
	}

	from := m.player.GetPos()
	delta := common.FacingToDelta(m.player.GetFacing())
	target := common.Coordinates{X: from.X + delta.X, Y: from.Y + delta.Y}

	if occupantID, bumped := occupancy.OccupantAt(target, m.player.GetID()); bumped {
		// Coming alongside offers a hail on the action bar rather than forcing a
		// modal open. Movement is accumulator-driven, so the player cannot stop
		// short of a ship on purpose; auto-opening meant an interruption nobody
		// asked for, and dismissing it re-triggered on the next step — roughly
		// every 400ms at full sail — until the player happened to turn away.
		m.alongsideNpcID = occupantID
		if m.npcs.GetByID(occupantID) != nil && m.lastHailedNpcID != occupantID {
			m.setNotice("Ship alongside. Press H to hail.")
		}
		return
	}
	m.alongsideNpcID = ""
	m.lastHailedNpcID = ""

	newPos, moved := sailing.TryStep(
		from,
		m.player.GetFacing(),
		m.player.GetID(),
		occupancy,
		m.world.IsPassableByBoat,
	)
	if moved {
		m.player.SetPos(newPos)
		occupancy[common.CoordToKey(newPos)] = m.player.GetID()
		delete(occupancy, common.CoordToKey(m.player.GetPreviousPos()))
	}
}

// tackPlayerPort turns the ship one octant to port. The player steers by tacking
// relative to the current heading, not by pointing at a compass direction.
func tackPlayerPort(m *GameState) {
	m.tack(common.TackPort)
}

// tackPlayerStarboard turns the ship one octant to starboard.
func tackPlayerStarboard(m *GameState) {
	m.tack(common.TackStarboard)
}

// tack applies one octant of helm, at most once per sailing tick.
func (m *GameState) tack(turn func(common.Facing) common.Facing) {
	if m.turnedThisTick {
		return
	}
	m.turnedThisTick = true
	m.player.SetHeading(turn(m.player.GetFacing()))
}

// hailTarget returns the NPC the player has come alongside, if any.
func (m *GameState) hailTarget() *npc.Npc {
	if m == nil || m.npcs == nil || m.alongsideNpcID == "" {
		return nil
	}
	return m.npcs.GetByID(m.alongsideNpcID)
}

// tryHailAdjacent opens the hail screen for the ship alongside.
func (m *GameState) tryHailAdjacent() bool {
	target := m.hailTarget()
	if target == nil {
		return false
	}
	m.lastHailedNpcID = m.alongsideNpcID
	m.openHail(hail.PayloadFromNPC(target))
	return true
}

// setPlayerSail jumps straight to a sail state, for the 1/2/3 presets.
func setPlayerSail(m *GameState, sail sailing.SailSetting) {
	m.player.SetSail(sail)
}

// trimPlayerSailMore sets one step more canvas, clamped at full.
func trimPlayerSailMore(m *GameState) {
	m.player.SetSail(m.player.GetSail().More())
}

// trimPlayerSailLess sets one step less canvas, clamped at furled.
func trimPlayerSailLess(m *GameState) {
	m.player.SetSail(m.player.GetSail().Less())
}
