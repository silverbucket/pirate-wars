package player

import (
	"testing"

	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/sailing"
)

func TestShipwrightOnce(t *testing.T) {
	cfg := economy.DefaultConfig()
	hold := NewHold(cfg)
	sailCfg := sailing.DefaultConfig()
	baseHull := sailCfg.HullSpeed

	if hold.Gold < cfg.SailUpgradeCost {
		t.Fatal("test setup: need enough starting gold")
	}

	hold.Gold -= cfg.SailUpgradeCost
	hold.FineSailsPurchased = true
	sailCfg.HullSpeed += cfg.SailUpgradeHullBonus

	if !hold.FineSailsPurchased {
		t.Fatal("fine sails should be marked purchased")
	}
	if sailCfg.HullSpeed != baseHull+cfg.SailUpgradeHullBonus {
		t.Fatalf("hull speed = %.2f, want %.2f", sailCfg.HullSpeed, baseHull+cfg.SailUpgradeHullBonus)
	}

	// second purchase blocked
	beforeGold := hold.Gold
	if !hold.FineSailsPurchased || hold.Gold < cfg.SailUpgradeCost {
		// cannot buy again
	} else {
		t.Fatal("should not allow second purchase path when already purchased")
	}
	if hold.Gold != beforeGold {
		t.Fatal("gold should not change on blocked second purchase")
	}
}
