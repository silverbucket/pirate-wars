package npc

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/resources"
	"pirate-wars/cmd/town"
	"pirate-wars/cmd/window"
	"pirate-wars/cmd/world"
	"sort"

	"go.uber.org/zap"
)

// ChanceToMove Percentage chance an NPC will calculate movement per tick
const ChanceToMove = 50
const GoalTypeTrade = 1

type Agenda struct {
	goal        int
	tradeTarget int
	tadeRoute   []town.Town
}

type Npc struct {
	id     string
	name   string
	eType  string
	flag   string
	ship   common.ShipType
	logger *zap.SugaredLogger
	avatar entities.Avatar
	agenda Agenda
}

type Npcs struct {
	logger *zap.SugaredLogger
	list   []Npc
}

func (n *Npc) GetName() string {
	return n.name
}
func (n *Npc) GetType() string {
	return n.eType
}
func (n *Npc) GetFlag() string {
	return n.flag
}

func (n *Npc) GetPos() common.Coordinates {
	return n.avatar.GetPos()
}

func (n *Npc) GetPreviousPos() common.Coordinates {
	return n.avatar.GetPreviousPos()
}

func (n *Npc) SetPos(p common.Coordinates) {
	n.avatar.SetPos(p)
}

func (n *Npc) GetID() string {
	return n.avatar.GetID()
}

func (n *Npc) GetTileImage() image.Image {
	return n.avatar.GetTileImage()
}

func (n *Npc) GetViewableRange() window.Dimensions {
	return window.Dimensions{Width: 20, Height: 20}
}

func (n *Npc) Highlight(b bool) {
	n.avatar.Highlight(b)
}

func (n *Npc) IsHighlighted() bool {
	return n.avatar.IsHighlighted()
}

func (n *Npc) GetColor() color.Color {
	return n.avatar.GetColor()
}

func (ns *Npcs) ForEach(fn func(n Npc)) {
	for _, n := range ns.list {
		fn(n)
	}
}

func (ns *Npcs) Create(towns *town.Towns, world *world.MapView) {
	pos := world.RandomPositionDeepWater()
	tradeTowns := []town.Town{}

	tryCount := 0
	for {
		tryCount++

		newTown, _ := towns.GetRandomTown()
		// ensure towns are unique
		if len(tradeTowns) > 2 {
			break
		} else if len(tradeTowns) == 2 {
			if common.CoordsMatch(newTown.GetPos(), tradeTowns[0].GetPos()) || !newTown.AccessibleFrom(pos) {
				// either same town, or inaccessible from position
				if tryCount > 20 {
					// abort creation
					//ns.logger.Info(fmt.Sprintf("Failed creating npc at position %d, skipping [town: %v, accessible?: %v]", pos, newTown.GetPos(), newTown.AccessibleFrom(pos)))
					return
				}
				// try again
				continue
			}
		}
		tradeTowns = append(tradeTowns, newTown)
	}

	// c := entities.ColorPossibilities[rand.Intn(len(entities.ColorPossibilities)-1)]
	flag := common.GetRandomFlag()

	npc := Npc{
		eType:  "NPC",
		logger: ns.logger,
		name:   common.GenerateCaptainName(),
		flag:   flag.Name,
		ship:   flag.Ship,
		avatar: entities.CreateAvatar(pos, resources.GetShipTile(flag.Ship), flag.Color),
		agenda: Agenda{
			goal:        GoalTypeTrade,
			tradeTarget: 0,
			tadeRoute:   tradeTowns,
		},
	}
	ns.logger.Infof("[%v] NPC created at %d, %d", npc.id, pos.X, pos.Y)
	ns.list = append(ns.list, npc)
}

func Init(towns *town.Towns, world *world.MapView, logger *zap.SugaredLogger) *Npcs {
	ns := Npcs{
		logger: logger,
	}
	for i := 0; i < common.TotalNpcs; i++ {
		ns.Create(towns, world)
	}
	logger.Infof("NPCs initialized: %d", len(ns.list))
	return &ns
}

func (ns *Npcs) CalcMovements() {
	ns.logger.Infof("Calculating NPC movements: %d", len(ns.list))
	for i := range ns.list {
		if rand.Intn(100) > ChanceToMove {
			continue
		}

		npc := &ns.list[i]
		targetTown := &npc.agenda.tadeRoute[npc.agenda.tradeTarget]

		// if we're already at our destination, flip our trade route
		if targetTown.HeatMap.GetCost(npc.avatar.GetPos()) < 3 {
			oldTown := npc.agenda.tadeRoute[npc.agenda.tradeTarget]
			npc.agenda.tradeTarget = npc.agenda.tradeTarget ^ 1
			targetTown = &npc.agenda.tadeRoute[npc.agenda.tradeTarget]
			ns.logger.Info(fmt.Sprintf("[%v] NPC movement trade route switch town %v to town %v", npc.id, oldTown.GetPos(), targetTown.GetPos()))
		}

		// find next move by cost on heatmap
		opts := []town.DirectionCost{}
		for _, dir := range common.Directions {
			n := common.AddDirection(npc.GetPos(), dir)
			if !common.Inbounds(n) {
				// don't check out of bounds
				continue
			}
			//t.Logger.Debug(fmt.Sprintf("New heatmap coordinates check [%v][%v]", newX, newY))
			//t.Logger.Debug(fmt.Sprintf("Npc at %v, %v - checking square %v, %v cost:%v [lowest cost: %v]", newPosition.X, newPosition.Y, newX, newY, town.HeatMap[newX][newY], lowestCost)
			opts = append(opts, town.DirectionCost{Pos: n, Cost: targetTown.HeatMap.GetCost(n)})
		}

		pick := town.DecideDirection(opts, targetTown.GetPos())
		target := pick.Pos
		cost := pick.Cost
		npcpos := npc.GetPos()

		if target.X == npcpos.X && target.Y == npcpos.Y {
			ns.logger.Debug(fmt.Sprintf("[%v] NPC stuck at %+v! Travelling to town at %v (cost %v)", npc.id, npcpos, targetTown.GetPos(), cost))
		} else {
			ns.logger.Info(fmt.Sprintf("[%v] NPC moving from %v to %v (cost %v) (color: %v)", npc.id, npcpos, target, cost, npc.GetColor()))
			if !common.IsPositionAdjacent(npcpos, target) {
				ns.logger.Debug(fmt.Sprintf("[%v] NPC warp! from %v to %v", npc.id, npcpos, target))
			}
			npc.SetPos(target)
		}
	}
}

func (ns *Npcs) GetList() []Npc {
	return ns.list
}

func (ns *Npcs) GetVisible(c common.Coordinates, vr window.Dimensions) Npcs {
	vp := window.GetViewportRegion(c)
	viewable := map[int]Npc{}
	keys := []int{}
	for _, npc := range ns.list {
		p := npc.GetPos()
		if vp.IsPositionWithin(p) {
			keys = append(keys, p.X)
			viewable[p.X] = npc
		}
	}
	sorted := Npcs{}
	sort.Ints(keys)
	for _, key := range keys {
		sorted.list = append(sorted.list, viewable[key])
	}
	return sorted
}
