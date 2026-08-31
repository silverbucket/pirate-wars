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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type dockPage int

const (
	dockPageMenu dockPage = iota
	dockPageMerchant
	dockPageTavern
	dockPageShipwright
)

func (gs *GameState) dockOverlayContent() fyne.CanvasObject {
	if gs.dockTown == nil {
		return widget.NewLabel("No town nearby.")
	}
	switch gs.dockPage {
	case dockPageMerchant:
		return gs.dockMerchantContent(gs.dockTown)
	case dockPageTavern:
		return gs.dockTavernContent()
	case dockPageShipwright:
		return gs.dockShipwrightContent()
	default:
		return gs.dockMenuContent()
	}
}

func (gs *GameState) dockMenuContent() fyne.CanvasObject {
	title := widget.NewLabel(fmt.Sprintf("Dock — %s", gs.dockTown.GetName()))
	title.TextStyle = fyne.TextStyle{Bold: true}
	merchant := widget.NewButton("Merchant", func() {
		gs.dockPage = dockPageMerchant
		gs.refreshOverlay()
	})
	tavern := widget.NewButton("Tavern", func() {
		gs.tavernRumor = tavern.PickRumor(gs.economyCfg, gs.npcs, gs.towns, int(gs.clock.CurrentTick()))
		gs.dockPage = dockPageTavern
		gs.refreshOverlay()
	})
	shipwright := widget.NewButton("Shipwright", func() {
		gs.dockPage = dockPageShipwright
		gs.refreshOverlay()
	})
	closeBtn := widget.NewButton("Leave dock", func() {
		gs.closeDock()
	})
	return container.NewVBox(title, merchant, tavern, shipwright, closeBtn)
}

func (gs *GameState) dockMerchantContent(t *town.Town) fyne.CanvasObject {
	market := t.Market()
	lines := []fyne.CanvasObject{
		widget.NewLabel(fmt.Sprintf("Merchant — %s", t.GetName())),
		widget.NewLabel(fmt.Sprintf("Your gold: %d  Cargo: %d/%d", gs.hold.Gold, gs.hold.Cargo.Total(), gs.hold.Cargo.Capacity())),
	}
	for _, g := range economy.AllGoods {
		good := g
		buyPrice := market.BuyPrice(good)
		sellPrice := market.SellPrice(good, gs.economyCfg)
		lines = append(lines, widget.NewLabel(fmt.Sprintf(
			"%s — stock %d  buy %d / sell %d gold",
			good.Label(), market.Stock(good), buyPrice, sellPrice,
		)))
		buy := widget.NewButton(fmt.Sprintf("Buy 1 %s", good.Label()), func() {
			economy.BuyFromTown(market, &gs.hold.Cargo, &gs.hold.Gold, gs.economyCfg, good, 1)
			gs.refreshOverlay()
		})
		sell := widget.NewButton(fmt.Sprintf("Sell 1 %s", good.Label()), func() {
			economy.SellToTown(market, &gs.hold.Cargo, &gs.hold.Gold, gs.economyCfg, good, 1)
			gs.refreshOverlay()
		})
		lines = append(lines, container.NewHBox(buy, sell))
	}
	back := widget.NewButton("Back", func() {
		gs.dockPage = dockPageMenu
		gs.refreshOverlay()
	})
	lines = append(lines, back)
	return container.NewVBox(lines...)
}

func (gs *GameState) dockTavernContent() fyne.CanvasObject {
	rumor := gs.tavernRumor
	if rumor == "" {
		rumor = "The tavern has no fresh leads."
	}
	back := widget.NewButton("Back", func() {
		gs.dockPage = dockPageMenu
		gs.refreshOverlay()
	})
	return container.NewVBox(
		widget.NewLabel("Tavern"),
		widget.NewLabel(rumor),
		back,
	)
}

func (gs *GameState) dockShipwrightContent() fyne.CanvasObject {
	status := fmt.Sprintf("Fine sails — %d gold (+%.2f hull speed, once)",
		gs.economyCfg.SailUpgradeCost, gs.economyCfg.SailUpgradeHullBonus)
	if gs.hold.FineSailsPurchased {
		status = "Fine sails already fitted."
	}
	buy := widget.NewButton("Buy fine sails", func() {
		gs.buyFineSails()
		gs.refreshOverlay()
	})
	if gs.hold.FineSailsPurchased || gs.hold.Gold < gs.economyCfg.SailUpgradeCost {
		buy.Disable()
	}
	back := widget.NewButton("Back", func() {
		gs.dockPage = dockPageMenu
		gs.refreshOverlay()
	})
	return container.NewVBox(widget.NewLabel("Shipwright"), widget.NewLabel(status), buy, back)
}

func (gs *GameState) buyFineSails() {
	if gs.hold.FineSailsPurchased || gs.hold.Gold < gs.economyCfg.SailUpgradeCost {
		return
	}
	gs.hold.Gold -= gs.economyCfg.SailUpgradeCost
	gs.hold.FineSailsPurchased = true
	gs.sailingCfg.HullSpeed += gs.economyCfg.SailUpgradeHullBonus
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
	t := dock.AdjacentTown(gs.player.GetPos(), gs.world, gs.towns)
	if t == nil {
		return false
	}
	gs.openDock(t)
	return true
}

func openDockIfAdjacent(gs *GameState, pos common.Coordinates, bw interface{ IsPassableByBoat(common.Coordinates) bool }, towns *town.Towns) bool {
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
	gs.showOverlay()
	gs.updatePanels(gs.currentExamineEntity())
}

func (gs *GameState) closeDock() {
	gs.dockTown = nil
	gs.dockPage = dockPageMenu
	gs.tavernRumor = ""
	ViewType = world.ViewTypeMainMap
	gs.hideOverlay()
	gs.updatePanels(gs.currentExamineEntity())
}

func (gs *GameState) openHail(payload hail.Payload) {
	gs.hailData = payload
	ViewType = world.ViewTypeHail
	gs.showOverlay()
	gs.updatePanels(gs.currentExamineEntity())
}

func (gs *GameState) closeHail() {
	gs.hailData = hail.Payload{}
	ViewType = world.ViewTypeMainMap
	gs.hideOverlay()
	gs.updatePanels(gs.currentExamineEntity())
}
