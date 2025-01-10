package terrain

import (
	"fmt"
	"math/rand"
	"pirate-wars/cmd/common"
)

type ColorScheme struct {
	Foreground string
	Background string
}

// ChanceToMove Percentage chance an NPC will calculate movement per tick
const ChanceToMove = 50

const GoalTypeTrade = 1

type Agenda struct {
	goal        int
	tradeTarget int
	tadeRoute   []Town
}

type Npc struct {
	id     string
	avatar *Avatar
	agenda Agenda
}

var ColorPossibilities = []ColorScheme{
	{"#AA0000", "#000000"},
	{"#00AA00", "#000000"},
	{"#00AA00", "#000000"},
	{"#00FF00", "#000000"},
	{"#0000FF", "#000000"},
	{"#FFFF00", "#000000"},
	{"#FF00FF", "#000000"},
	{"#00FFFF", "#000000"},
}

func (t *Terrain) CreateNpc() {
	pos := t.RandomPositionDeepWater()
	firstTown := t.GetRandomTown()
	secondTown := t.GetRandomTown()
	for {
		// ensure towns are unique
		if secondTown.GetY() == firstTown.GetY() && secondTown.GetY() == firstTown.GetY() {
			secondTown = t.GetRandomTown()
		} else {
			break
		}
	}
	npc := Npc{
		id:     common.GenID(pos),
		avatar: CreateAvatar(pos, '⏏', ColorPossibilities[rand.Intn(len(ColorPossibilities)-1)]),
		agenda: Agenda{
			goal:        GoalTypeTrade,
			tradeTarget: 0,
			tadeRoute:   []Town{firstTown, secondTown},
		},
	}
	t.Logger.Infof("[%v] NPC created at %d, %d", npc.id, pos.X, pos.Y)
	t.Npcs = append(t.Npcs, npc)
}

func (t *Terrain) InitNpcs() {
	for i := 0; i < common.TotalNpcs; i++ {
		t.CreateNpc()
	}
}

func (t *Terrain) CalcNpcMovements() {
	for i := range t.Npcs {
		if rand.Intn(100) > ChanceToMove {
			continue
		}

		npc := &t.Npcs[i]
		target := npc.avatar.GetPos()
		town := &npc.agenda.tadeRoute[npc.agenda.tradeTarget]

		// if we're already at our destination, flip our trade route
		if town.GetHeatmapCost(target) < 3 {
			oldTown := npc.agenda.tadeRoute[npc.agenda.tradeTarget]
			npc.agenda.tradeTarget = npc.agenda.tradeTarget ^ 1
			town = &npc.agenda.tadeRoute[npc.agenda.tradeTarget]
			t.Logger.Info(fmt.Sprintf("[%v] NPC movement trade switch from town %v to town %v", npc.id, oldTown.GetPos(), town.GetPos()))
		}

		// find next move by cost on heatmap
		var lowestCost = common.MaxMovementCost
		for _, dir := range common.Directions {
			n := common.Coordinates{npc.avatar.GetX() + dir.X, npc.avatar.GetY() + dir.Y}
			if !common.Inbounds(n) {
				// don't check out of bounds
				continue
			}

			//t.Logger.Debug(fmt.Sprintf("New heatmap coordinates check [%v][%v]", newX, newY))
			//t.Logger.Debug(fmt.Sprintf("Npc at %v, %v - checking square %v, %v cost:%v [lowest cost: %v]", newPosition.X, newPosition.Y, newX, newY, town.heatMap[newX][newY], lowestCost)
			cost := town.GetHeatmapCost(n)
			if cost >= 0 && cost < lowestCost {
				lowestCost = cost
				target = n
			}
		}

		if target.X == npc.avatar.GetX() && target.Y == npc.avatar.GetY() {
			t.Logger.Debug(fmt.Sprintf("[%v] NPC stuck! Travelling to town at %v (cost %v)", npc.id, town.GetPos(), town.GetHeatmapCost(target)))
		} else {
			pos := npc.avatar.GetPos()
			t.Logger.Info(fmt.Sprintf("[%v] NPC moving from %v to %v", npc.id, pos, target))
			if !t.isPositionAdjacent(pos, target) {
				t.Logger.Warn(fmt.Sprintf("[%v] NPC warp!", npc.id))
			}
			npc.avatar.SetPos(target)
		}
	}
}

func (t *Terrain) GetNpcAvatars() []Avatar {
	var avs []Avatar
	for npc := range t.Npcs {
		avs = append(avs, *t.Npcs[npc].avatar)
	}
	return avs
}
