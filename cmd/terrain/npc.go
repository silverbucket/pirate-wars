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
	{"9", "#000000"},   // strong red
	{"10", "#000000"},  // green
	{"11", "#000000"},  // yellow
	{"14", "#000000"},  // bright cyan
	{"15", "#000000"},  // off-white
	{"46", "#000000"},  // blue/green
	{"65", "#000000"},  // faded cyan
	{"86", "#000000"},  // light cyan
	{"172", "#000000"}, // off pink
	{"201", "#000000"}, // pink
	{"207", "#000000"}, // light pink
	{"218", "#000000"}, // light pink/white
	{"222", "#000000"}, // light yellow/orange
	{"230", "#000000"}, // yellow/white
	{"253", "#000000"}, // grey
}

func (t *Terrain) CreateNpc() {
	pos := t.RandomPositionDeepWater()
	tradeTowns := []Town{}
	for {
		town := t.GetRandomTown()
		// ensure towns are unique
		if len(tradeTowns) > 2 {
			break
		} else if len(tradeTowns) == 2 {
			if (town.GetY() == tradeTowns[0].GetY() && town.GetY() == tradeTowns[0].GetY()) || !town.AccessibleFrom(pos) {
				// either same town, or inaccessible from position, try again
				continue
			}
		}
		tradeTowns = append(tradeTowns, town)
	}

	npc := Npc{
		id:     common.GenID(pos),
		avatar: CreateAvatar(pos, '⏏', ColorPossibilities[rand.Intn(len(ColorPossibilities)-1)]),
		agenda: Agenda{
			goal:        GoalTypeTrade,
			tradeTarget: 0,
			tadeRoute:   tradeTowns,
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
		if town.HeatMap.GetCost(npc.avatar.GetPos()) < 3 {
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
			//t.Logger.Debug(fmt.Sprintf("Npc at %v, %v - checking square %v, %v cost:%v [lowest cost: %v]", newPosition.X, newPosition.Y, newX, newY, town.HeatMap[newX][newY], lowestCost)
			opts = append(opts, DirectionCost{n, town.HeatMap.GetCost(n)})
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
	lowestCost := MaxMovementCost
	choice := DirectionCost{}
	for _, e := range o {
		if e.cost <= lowestCost && e.cost >= 0 {
			lowestCost = e.cost
			//possibilities = append(possibilities, e.pos)
			choice = e
		}
	}
	//return common.ClosestTo(dest, possibilities)
	return choice
}
