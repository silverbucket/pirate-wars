package npc

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/sailing"
	"pirate-wars/cmd/town"
	"pirate-wars/cmd/window"
	"pirate-wars/cmd/world"
	"sort"

	"go.uber.org/zap"
)

// ChanceToMove is deprecated; NPC skip percent now comes from sailing.cfg.
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
		avatar: entities.CreateAvatar(pos, flag.Ship, flag.Color),
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

func (n *Npc) GetFacing() common.Facing {
	return n.avatar.GetFacing()
}

func (n *Npc) MovedThisTick() bool {
	return n.avatar.MovedThisTick()
}

func (n *Npc) ClearMovedFlag() {
	n.avatar.ClearMovedFlag()
}

func (ns *Npcs) ClearMovedFlags() {
	for i := range ns.list {
		ns.list[i].ClearMovedFlag()
	}
}

func (ns *Npcs) ResolveMovements(cfg sailing.Config, wind *sailing.Wind, world *world.MapView, occupancy sailing.Occupancy) {
	for i := range ns.list {
		if rand.Intn(100) < cfg.NPCSkipPercent {
			continue
		}

		npc := &ns.list[i]
		targetTown := &npc.agenda.tadeRoute[npc.agenda.tradeTarget]

		if targetTown.HeatMap.GetCost(npc.avatar.GetPos()) < 3 {
			oldTown := npc.agenda.tadeRoute[npc.agenda.tradeTarget]
			npc.agenda.tradeTarget = npc.agenda.tradeTarget ^ 1
			targetTown = &npc.agenda.tadeRoute[npc.agenda.tradeTarget]
			ns.logger.Info(fmt.Sprintf("[%v] NPC movement trade route switch town %v to town %v", npc.id, oldTown.GetPos(), targetTown.GetPos()))
		}

		opts := []town.DirectionCost{}
		for _, dir := range common.Directions {
			n := common.AddDirection(npc.GetPos(), dir)
			if !common.Inbounds(n) {
				continue
			}
			opts = append(opts, town.DirectionCost{Pos: n, Cost: targetTown.HeatMap.GetCost(n)})
		}

		pick := town.DecideDirection(opts, targetTown.GetPos())
		target := pick.Pos
		npcpos := npc.GetPos()

		if target.X == npcpos.X && target.Y == npcpos.Y {
			ns.logger.Debug(fmt.Sprintf("[%v] NPC stuck at %+v! Travelling to town at %v (cost %v)", npc.id, npcpos, targetTown.GetPos(), pick.Cost))
			continue
		}

		dx := target.X - npcpos.X
		dy := target.Y - npcpos.Y
		if facing, ok := common.FacingFromDelta(dx, dy); ok {
			npc.avatar.SetHeading(facing)
		}

		speed := cfg.EffectiveSpeed(npc.avatar.GetFacing(), npc.avatar.GetSail(), wind)
		if !sailing.ShouldMove(speed) {
			continue
		}

		newPos, moved := sailing.TryStep(npcpos, npc.avatar.GetFacing(), npc.GetID(), occupancy, world.IsPassableByBoat)
		if moved {
			ns.logger.Info(fmt.Sprintf("[%v] NPC moving from %v to %v (color: %v)", npc.id, npcpos, newPos, npc.GetColor()))
			npc.SetPos(newPos)
			occupancy[common.CoordToKey(newPos)] = npc.GetID()
			delete(occupancy, common.CoordToKey(npcpos))
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
