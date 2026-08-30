package tavern

import (
	"strings"
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/town"
)

var flavorRumors = []string{
	"hidden cove",
	"navy patrols",
	"merchant fleet",
}

func TestPickRumorShipIntelNamesRealNPC(t *testing.T) {
	cfg := economy.DefaultConfig()
	dest := town.NewTownForTest(common.Coordinates{X: 2, Y: 3}, cfg)

	trader := npc.Npc{}
	trader.SetTestState("Anne Bonny", dest, economy.GoodCloth, 12)
	npcs := npc.TestNpcsWith(trader)

	rumor := PickRumor(cfg, npcs, town.TestTownsWith(dest), 0)
	if !strings.Contains(rumor, "Anne Bonny") {
		t.Fatalf("rumor = %q, want trader name", rumor)
	}
	if !strings.Contains(rumor, dest.GetName()) {
		t.Fatalf("rumor = %q, want destination %q", rumor, dest.GetName())
	}
	if !strings.Contains(rumor, "Cloth") {
		t.Fatalf("rumor = %q, want cargo good", rumor)
	}
	assertNoFlavorRumor(t, rumor)
}

func TestPickRumorTownShortageNamesRealTownAndGood(t *testing.T) {
	cfg := economy.DefaultConfig()
	shortTown := town.NewTownForTest(common.Coordinates{X: 5, Y: 6}, cfg)
	shortTown.Market().SetStock(economy.GoodPowder, cfg.PowderStockMin, cfg)

	rumor := PickRumor(cfg, npc.TestNpcsWith(), town.TestTownsWith(shortTown), 1)
	if !strings.Contains(rumor, shortTown.GetName()) {
		t.Fatalf("rumor = %q, want town %q", rumor, shortTown.GetName())
	}
	if !strings.Contains(rumor, "Powder") {
		t.Fatalf("rumor = %q, want short good", rumor)
	}
	assertNoFlavorRumor(t, rumor)
}

func TestPickRumorNeverUsesFlavorLines(t *testing.T) {
	cfg := economy.DefaultConfig()
	dest := town.NewTownForTest(common.Coordinates{X: 1, Y: 1}, cfg)
	shortTown := town.NewTownForTest(common.Coordinates{X: 9, Y: 9}, cfg)
	shortTown.Market().SetStock(economy.GoodRum, cfg.RumStockMin, cfg)

	trader := npc.Npc{}
	trader.SetTestState("Calico Jack", dest, economy.GoodRum, 5)
	npcs := npc.TestNpcsWith(trader)
	towns := town.TestTownsWith(dest, shortTown)

	for seed := 0; seed < 8; seed++ {
		rumor := PickRumor(cfg, npcs, towns, seed)
		assertNoFlavorRumor(t, rumor)
	}
}

func assertNoFlavorRumor(t *testing.T, rumor string) {
	t.Helper()
	lower := strings.ToLower(rumor)
	for _, flavor := range flavorRumors {
		if strings.Contains(lower, flavor) {
			t.Fatalf("rumor = %q contains flavor line %q", rumor, flavor)
		}
	}
}
