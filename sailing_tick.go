package main

import (
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/sailing"
)

func (m *GameState) resolveSailingTick() {
	m.wind.Tick()

	m.player.ClearMovedFlag()
	m.npcs.ClearMovedFlags()

	occupancy := m.buildOccupancy()

	m.resolvePlayerMovement(occupancy)
	m.npcs.ResolveMovements(m.sailingCfg, m.wind, m.world, occupancy)
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
	if !sailing.ShouldMove(speed) {
		return
	}
	newPos, moved := sailing.TryStep(
		m.player.GetPos(),
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

func setPlayerHeading(m *GameState, d common.Coordinates) {
	if f, ok := common.FacingFromDirection(d); ok {
		m.player.SetHeading(f)
	}
}

func setPlayerSail(m *GameState, sail sailing.SailSetting) {
	m.player.SetSail(sail)
}

func cyclePlayerSail(m *GameState) {
	m.player.SetSail(m.player.GetSail().Next())
}
