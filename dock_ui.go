package main

import (
	"fmt"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/dock"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/hail"
	"pirate-wars/cmd/tavern"
	"pirate-wars/cmd/town"
	"pirate-wars/cmd/world"
)

type dockPage int

const (
	dockPageMenu dockPage = iota
	dockPageMerchant
	dockPageTavern
	dockPageShipwright
)

func (gs *GameState) dockOverlayScreen() overlayScreen {
	if gs.dockTown == nil {
		return overlayScreen{
			title: "Dock",
			rows: []overlayRow{
				{text: "No town nearby.", dim: true},
				{buttons: []overlayAction{{label: barLabelFor(dockKeyMap(), "Leave dock"), do: gs.closeDock}}},
			},
		}
	}
	switch gs.dockPage {
	case dockPageMerchant:
		return gs.dockMerchantScreen(gs.dockTown)
	case dockPageTavern:
		return gs.dockTavernScreen()
	case dockPageShipwright:
		return gs.dockShipwrightScreen()
	default:
		return gs.dockMenuScreen()
	}
}

func (gs *GameState) dockMenuScreen() overlayScreen {
	return overlayScreen{
		title: fmt.Sprintf("Dock - %s", gs.dockTown.GetName()),
		rows: []overlayRow{
			{buttons: []overlayAction{{label: "Merchant", do: func() { gs.openDockPage(dockPageMerchant) }}}},
			{buttons: []overlayAction{{label: "Tavern", do: gs.openTavern}}},
			{buttons: []overlayAction{{label: "Shipwright", do: func() { gs.openDockPage(dockPageShipwright) }}}},
			{buttons: []overlayAction{{label: barLabelFor(dockKeyMap(), "Leave dock"), do: gs.closeDock}}},
		},
	}
}

func (gs *GameState) dockMerchantScreen(t *town.Town) overlayScreen {
	market := t.Market()
	rows := []overlayRow{
		{
			text: fmt.Sprintf("Gold: %d   Cargo: %d/%d",
				gs.hold.Gold, gs.hold.Cargo.Total(), gs.hold.Cargo.Capacity()),
			dim: true,
		},
	}
	for _, g := range economy.AllGoods {
		good := g
		rows = append(rows, overlayRow{
			text: fmt.Sprintf("%-10s stock %3d  buy %4d  sell %4d",
				good.Label(), market.Stock(good), market.BuyPrice(good), market.SellPrice(good, gs.economyCfg)),
			buttons: []overlayAction{
				{label: "Buy", do: func() {
					economy.BuyFromTown(market, &gs.hold.Cargo, &gs.hold.Gold, gs.economyCfg, good, 1)
				}},
				{label: "Sell", do: func() {
					economy.SellToTown(market, &gs.hold.Cargo, &gs.hold.Gold, gs.economyCfg, good, 1)
				}},
			},
		})
	}
	rows = append(rows, overlayRow{buttons: []overlayAction{{label: "Back", do: func() { gs.openDockPage(dockPageMenu) }}}})
	return overlayScreen{title: fmt.Sprintf("Merchant - %s", t.GetName()), rows: rows}
}

func (gs *GameState) dockTavernScreen() overlayScreen {
	rumor := gs.tavernRumor
	if rumor == "" {
		rumor = "The tavern has no fresh leads."
	}
	rows := []overlayRow{}
	for _, line := range wrapText(rumor, 68) {
		rows = append(rows, overlayRow{text: line})
	}
	rows = append(rows, overlayRow{buttons: []overlayAction{{label: "Back", do: func() { gs.openDockPage(dockPageMenu) }}}})
	return overlayScreen{title: "Tavern", rows: rows}
}

func (gs *GameState) dockShipwrightScreen() overlayScreen {
	status := fmt.Sprintf("Fine sails - %d gold (+%.2f hull speed, once)",
		gs.economyCfg.SailUpgradeCost, gs.economyCfg.SailUpgradeHullBonus)
	rows := []overlayRow{{text: status}}
	if gs.hold.FineSailsPurchased {
		rows = []overlayRow{{text: "Fine sails already fitted.", dim: true}}
	} else if gs.hold.Gold < gs.economyCfg.SailUpgradeCost {
		rows = append(rows, overlayRow{text: "Not enough gold.", dim: true})
	} else {
		rows = append(rows, overlayRow{buttons: []overlayAction{{label: "Buy fine sails", do: gs.buyFineSails}}})
	}
	rows = append(rows, overlayRow{buttons: []overlayAction{{label: "Back", do: func() { gs.openDockPage(dockPageMenu) }}}})
	return overlayScreen{title: "Shipwright", rows: rows}
}

func (gs *GameState) buyFineSails() {
	if gs.hold.FineSailsPurchased || gs.hold.Gold < gs.economyCfg.SailUpgradeCost {
		return
	}
	gs.hold.Gold -= gs.economyCfg.SailUpgradeCost
	gs.hold.FineSailsPurchased = true
	gs.sailingCfg.HullSpeed += gs.economyCfg.SailUpgradeHullBonus
}

func (gs *GameState) openDockPage(p dockPage) {
	gs.dockPage = p
}

func (gs *GameState) openTavern() {
	gs.tavernRumor = tavern.PickRumor(gs.economyCfg, gs.npcs, gs.towns, int(gs.clock.CurrentTick()))
	gs.dockPage = dockPageTavern
}

// adjacentDockTown returns the dockable town next to the player, tolerating a
// partially built game state (no player or world yet).
func (gs *GameState) adjacentDockTown() *town.Town {
	if gs == nil || gs.player == nil || gs.world == nil {
		return nil
	}
	return dock.AdjacentTown(gs.player.GetPos(), gs.world, gs.towns)
}

func (gs *GameState) enterDock(t *town.Town) {
	gs.dockTown = t
	gs.dockPage = dockPageMenu
	ViewType = world.ViewTypeDock
}

func (gs *GameState) tryOpenDock() bool {
	if ViewType != world.ViewTypeMainMap {
		return false
	}
	t := gs.adjacentDockTown()
	if t == nil {
		return false
	}
	gs.openDock(t)
	return true
}

// openDockIfAdjacent docks from an explicit position and world, for tests and for
// callers that are not the player's own state.
func openDockIfAdjacent(gs *GameState, pos common.Coordinates, bw interface {
	IsPassableByBoat(common.Coordinates) bool
}, towns *town.Towns) bool {
	if ViewType != world.ViewTypeMainMap {
		return false
	}
	t := dock.AdjacentTown(pos, bw, towns)
	if t == nil {
		return false
	}
	gs.enterDock(t)
	return true
}

func (gs *GameState) openDock(t *town.Town) {
	gs.enterDock(t)
}

func (gs *GameState) closeDock() {
	gs.dockTown = nil
	gs.dockPage = dockPageMenu
	gs.tavernRumor = ""
	ViewType = world.ViewTypeMainMap
}

func (gs *GameState) openHail(payload hail.Payload) {
	gs.hailData = payload
	ViewType = world.ViewTypeHail
}

func (gs *GameState) closeHail() {
	gs.hailData = hail.Payload{}
	ViewType = world.ViewTypeMainMap
}
