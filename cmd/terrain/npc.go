package terrain

import (
	"fmt"
	"math/rand"
	"pirate-wars/cmd/common"
	"sort"
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
	avatar Avatar
	agenda Agenda
}

type Npcs []Npc

var ColorPossibilities = []ColorScheme{
	{"9", "0"},   // strong red
	{"10", "0"},  // green
	{"11", "0"},  // yellow
	{"14", "0"},  // bright cyan
	{"15", "0"},  // off-white
	{"46", "0"},  // blue/green
	{"69", "0"},  // faded cyan
	{"86", "0"},  // light cyan
	{"93", "0"},  // fuchsia
	{"172", "0"}, // off pink
	{"193", "0"}, // light green
	{"201", "0"}, // pink
	{"207", "0"}, // light pink
	{"211", "0"}, // lighter pink
	{"218", "0"}, // light pink/white
	{"222", "0"}, // light yellow/orange
	{"230", "0"}, // yellow/white
	{"253", "0"}, // grey
	{"255", "0"}, // white
}

func (n *Npc) GetPos() common.Coordinates {
	return n.avatar.pos
}

func (n *Npc) SetPos(p common.Coordinates) {
	n.avatar.pos = p
}

func (n *Npc) GetID() string {
	return n.id
}

func (n *Npc) SetID(s string) {
	n.id = s
}

func (n *Npc) GetForegroundColor() string {
	return n.avatar.fgColor
}

func (n *Npc) Render() string {
	return n.avatar.Render()
}

func (n *Npc) GetBackgroundColor() string {
	return n.avatar.bgColor
}

func (n *Npc) Highlight() {
	n.avatar.SetBlink(true)
	n.avatar.SetBackgroundColor("7")
}

func (ns *Npcs) ForEach(fn func(n Npc)) {
	for _, n := range *ns {
		fn(n)
	}
}

func (t *Terrain) CreateNpc() {
	pos := t.RandomPositionDeepWater()
	tradeTowns := []Town{}

	tryCount := 0
	for {
		tryCount++

		town := t.GetRandomTown()
		// ensure towns are unique
		if len(tradeTowns) > 2 {
			break
		} else if len(tradeTowns) == 2 {
			if common.CoordsMatch(town.GetPos(), tradeTowns[0].GetPos()) || !town.AccessibleFrom(pos) {
				// either same town, or inaccessible from position
				if tryCount > 20 {
					// abort creation
					t.Logger.Info(fmt.Sprintf("Failed creating npc at position %d, skipping [town: %v, accessible?: %v]", pos, town.GetPos(), town.AccessibleFrom(pos)))
					return
				}
				// try again
				continue
			}
		}
		tradeTowns = append(tradeTowns, town)
	}

	color := ColorPossibilities[rand.Intn(len(ColorPossibilities)-1)]
	npc := Npc{
		id:     common.GenID(pos),
		avatar: CreateAvatar(pos, '⏏', color),
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
	t.Logger.Infof("Npcs initialization completed: %d", len(t.Npcs))
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
			n := common.AddDirection(npc.GetPos(), dir)
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
		npcpos := npc.GetPos()

		if target.X == npcpos.X && target.Y == npcpos.Y {
			t.Logger.Debug(fmt.Sprintf("[%v] NPC stuck at %v! Travelling to town at %v (cost %v)", npc.id, npcpos, town.GetPos(), cost))
		} else {
			t.Logger.Info(fmt.Sprintf("[%v] NPC moving from %v to %v (cost %v) (bg color: %v)", npc.id, npcpos, target, cost, npc.GetBackgroundColor()))
			if !common.IsPositionAdjacent(npcpos, target) {
				t.Logger.Debug(fmt.Sprintf("[%v] NPC warp! from %v to %v", npc.id, npcpos, target))
			}
			npc.SetPos(target)
		}
	}
}

func (t *Terrain) GetNpcs() Npcs {
	var avs Npcs
	for _, npc := range t.Npcs {
		avs = append(avs, npc)
	}
	return avs
}

func (t *Terrain) GetVisibleNpcs(c common.Coordinates) Npcs {
	v := common.GetViewableArea(c)
	viewable := map[int]Npc{}
	keys := []int{}
	for _, npc := range t.Npcs {
		p := npc.GetPos()
		if common.IsPositionWithin(p, v) {
			keys = append(keys, p.X)
			viewable[p.X] = npc
		}
	}
	sorted := Npcs{}
	sort.Ints(keys)
	for _, key := range keys {
		sorted = append(sorted, viewable[key])
	}
	return sorted
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
