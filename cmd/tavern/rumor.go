package tavern

import (
	"fmt"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/hail"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/town"
)

// PickRumor returns one live intel line for a tavern visit.
// Prefers ship intel when traders are sailing with cargo; otherwise town shortages.
func PickRumor(cfg economy.Config, npcs *npc.Npcs, towns *town.Towns, seed int) string {
	ship := collectShipRumors(npcs)
	short := collectTownShortRumors(cfg, towns)

	switch {
	case len(ship) > 0 && len(short) == 0:
		return ship[seed%len(ship)]
	case len(short) > 0 && len(ship) == 0:
		return short[seed%len(short)]
	case len(ship) > 0 && len(short) > 0:
		if seed%2 == 0 {
			return ship[seed%len(ship)]
		}
		return short[seed%len(short)]
	default:
		return "The tavern has no fresh leads."
	}
}

func collectShipRumors(npcs *npc.Npcs) []string {
	if npcs == nil {
		return nil
	}
	var out []string
	list := npcs.GetList()
	for i := range list {
		n := list[i]
		if n.TraderAmount() <= 0 {
			continue
		}
		p := hail.PayloadFromNPC(&n)
		out = append(out, fmt.Sprintf("%s sails for %s carrying %s.", p.Name, p.Dest, p.Cargo))
	}
	return out
}

func collectTownShortRumors(cfg economy.Config, towns *town.Towns) []string {
	if towns == nil {
		return nil
	}
	var out []string
	for _, t := range towns.GetTowns() {
		market := t.Market()
		for _, g := range economy.AllGoods {
			if economy.IsShortOnGood(cfg, market.Stock(g), g) {
				out = append(out, fmt.Sprintf("%s is short on %s (stock %d).", t.GetName(), g.Label(), market.Stock(g)))
			}
		}
	}
	return out
}
