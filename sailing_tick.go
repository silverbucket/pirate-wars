package main

import (
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/hail"
	"pirate-wars/cmd/sailing"
)

func (m *GameState) resolveSailingTick() {
	m.wind.Tick()

	m.player.ClearMovedFlag()
	m.npcs.ClearMovedFlags()

	occupancy := m.buildOccupancy()

	m.resolvePlayerMovement(occupancy)
	m.npcs.ResolveMovements(m.sailingCfg, m.economyCfg, m.wind, m.world, m.towns, occupancy)
}

func (m *GameState) buildOccupancy() sailing.Occupancy {
	ships := make(map[string]common.Coordinates)
	ships[m.player.GetID()] = m.player.GetPos()
	for _, n := range m.npcs.GetList() {
		ships[n.GetID()] = n.GetPos()
	}
	return sailing.NewOccupancy(ships)
}

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
		if npc := m.npcs.GetByID(occupantID); npc != nil {
			m.openHail(hail.PayloadFromNPC(npc))
		}
		return
	}

	newPos, moved := sailing.TryStep(
		from,
		m.player.GetFacing(),
		m.player.GetID(),
		occupancy,
		m.sailingWorld.IsPassableByBoat,
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
	m.player.SetHeading(common.TackPort(m.player.GetFacing()))
}

// tackPlayerStarboard turns the ship one octant to starboard.
func tackPlayerStarboard(m *GameState) {
	m.player.SetHeading(common.TackStarboard(m.player.GetFacing()))
}

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
