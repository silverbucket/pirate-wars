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
		town := &npc.agenda.tadeRoute[npc.agenda.tradeTarget]

		// if we're already at our destination, flip our trade route
		if town.heatMap.GetCost(npc.avatar.GetPos()) < 3 {
			oldTown := npc.agenda.tadeRoute[npc.agenda.tradeTarget]
			npc.agenda.tradeTarget = npc.agenda.tradeTarget ^ 1
			town = &npc.agenda.tadeRoute[npc.agenda.tradeTarget]
			t.Logger.Info(fmt.Sprintf("[%v] NPC movement trade route switch town %v to town %v", npc.id, oldTown.GetPos(), town.GetPos()))
		}

		// find next move by cost on heatmap
		opts := []DirectionCost{}
		for _, dir := range common.Directions {
			n := common.AddDirection(npc.avatar.GetPos(), dir)
			if !common.Inbounds(n) {
				// don't check out of bounds
				continue
			}
			//t.Logger.Debug(fmt.Sprintf("New heatmap coordinates check [%v][%v]", newX, newY))
			//t.Logger.Debug(fmt.Sprintf("Npc at %v, %v - checking square %v, %v cost:%v [lowest cost: %v]", newPosition.X, newPosition.Y, newX, newY, town.heatMap[newX][newY], lowestCost)
			opts = append(opts, DirectionCost{n, town.heatMap.GetCost(n)})
		}

		pick := decideDirection(opts, town.GetPos())
		target := pick.pos
		cost := pick.cost
		npcpos := npc.avatar.GetPos()

		if target.X == npc.avatar.GetX() && target.Y == npc.avatar.GetY() {
			t.Logger.Debug(fmt.Sprintf("[%v] NPC stuck at %v! Travelling to town at %v (cost %v)", npc.id, npcpos, town.GetPos(), cost))
		} else {
			t.Logger.Info(fmt.Sprintf("[%v] NPC moving from %v to %v (cost %v)", npc.id, npcpos, target, cost))
			if !common.IsPositionAdjacent(npcpos, target) {
				t.Logger.Warn(fmt.Sprintf("[%v] NPC warp!", npc.id))
			}
			npc.avatar.SetPos(target)
		}
	}
}

func (t *Terrain) GetNpcAvatars() []AvatarReadOnly {
	var avs []AvatarReadOnly
	for npc := range t.Npcs {
		avs = append(avs, t.Npcs[npc].avatar)
	}
	return avs
}

func decideDirection(o []DirectionCost, dest common.Coordinates) DirectionCost {
	lowestCost := common.MaxMovementCost
	choice := DirectionCost{}
	for _, e := range o {
		if e.cost <= lowestCost {
			lowestCost = e.cost
			//possibilities = append(possibilities, e.pos)
			choice = e
		}
	}
	//return common.ClosestTo(dest, possibilities)
	return choice
}
