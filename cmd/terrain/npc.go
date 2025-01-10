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
	t.Logger.Infof("Creating NPC at %d, %d", pos.X, pos.Y)
	firstTown := t.GetRandomTown()
	secondTown := t.GetRandomTown()
	for {
		// ensure towns are unique
		if secondTown.pos.X == firstTown.pos.X && secondTown.pos.Y == firstTown.pos.Y {
			secondTown = t.GetRandomTown()
		} else {
			break
		}
	}
	npc := Npc{
		avatar: CreateAvatar(pos, '⏏', ColorPossibilities[rand.Intn(len(ColorPossibilities)-1)]),
		agenda: Agenda{
			goal:        GoalTypeTrade,
			tradeTarget: 0,
			tadeRoute:   []Town{firstTown, secondTown},
		},
	}
	t.Npcs = append(t.Npcs, npc)
}

func (t *Terrain) InitNpcs() {
	for i := 0; i < common.TotalNpcs; i++ {
		t.CreateNpc()
	}
}

func (t *Terrain) CalcNpcMovements() {
	for _, npc := range t.Npcs {
		if rand.Intn(100) > ChanceToMove {
			continue
		}

		target := npc.avatar.GetPos()
		town := npc.agenda.tadeRoute[npc.agenda.tradeTarget]

		// if we're already at our destination, flip our trade route
		if town.heatMap[target.X][target.Y] < 3 {
			npc.agenda.tradeTarget = npc.agenda.tradeTarget ^ 1
			town = npc.agenda.tadeRoute[npc.agenda.tradeTarget]
		}

		// find next move by cost on heatmap
		var lowestCost = common.MaxMovementCost
		for _, dir := range common.Directions {
			newX, newY := npc.avatar.GetX()+dir.X, npc.avatar.GetY()+dir.Y
			if (newX < 0 || newX > common.WorldHeight-1) || (newY < 0 || newY > common.WorldHeight-1) {
				// don't check out of bounds
				continue
			}
			//t.Logger.Debug(fmt.Sprintf("New heatmap coordinates check [%v][%v]", newX, newY))
			//t.Logger.Debug(fmt.Sprintf("Npc at %v, %v - checking square %v, %v cost:%v [lowest cost: %v]", newPosition.X, newPosition.Y, newX, newY, town.heatMap[newX][newY], lowestCost)

			if town.heatMap[newX][newY] >= 0 && town.heatMap[newX][newY] < lowestCost {
				lowestCost = town.heatMap[newX][newY]
				target = common.Coordinates{newX, newY}
			}
		}
		if target.X == npc.avatar.GetX() && target.Y == npc.avatar.GetY() {
			t.Logger.Debug(fmt.Sprintf("NPC stuck! Travelling to town at %v, %v (cost square %v)", town.pos.X, town.pos.Y, town.heatMap[target.X][target.Y]))
		} else {
			t.Logger.Info(fmt.Sprintf("NPC moving to %v, %v", target.X, target.Y))
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
