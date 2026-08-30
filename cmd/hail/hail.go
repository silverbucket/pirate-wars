package hail

import (
	"fmt"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/npc"
)

// Payload is the hail overlay content for a bumped trader.
type Payload struct {
	Name    string
	Dest    string
	Cargo   string
	NpcID   string
}

func PayloadFromNPC(n *npc.Npc) Payload {
	if n == nil {
		return Payload{}
	}
	dest := "unknown"
	if t, ok := n.DestinationTown(); ok {
		dest = t.GetName()
	}
	return Payload{
		Name:  n.GetName(),
		Dest:  dest,
		Cargo: formatTraderCargo(n.TraderGood(), n.TraderAmount()),
		NpcID: n.GetID(),
	}
}

func formatTraderCargo(g economy.Good, amount int) string {
	if amount <= 0 {
		return "empty"
	}
	return fmt.Sprintf("%s x%d", g.Label(), amount)
}

func (p Payload) Text() string {
	return fmt.Sprintf("Hail: %s\nDestination: %s\nCargo: %s", p.Name, p.Dest, p.Cargo)
}
