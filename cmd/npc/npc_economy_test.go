package npc

import (
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/town"
)

func TestTraderDumpsOneGoodAtTown(t *testing.T) {
	cfg := economy.DefaultConfig()
	marketTown := town.NewTownForTest(common.Coordinates{X: 3, Y: 4}, cfg)
	startStock := marketTown.Market().Stock(economy.GoodPowder)

	n := &Npc{}
	n.SetTestState("Trader", marketTown, economy.GoodPowder, 9)

	n.dumpCargoAtTown(&marketTown)

	if n.TraderAmount() != 0 {
		t.Fatalf("trader cargo = %d, want 0 after dump", n.TraderAmount())
	}
	if got := marketTown.Market().Stock(economy.GoodPowder); got != startStock+9 {
		t.Fatalf("town stock = %d, want %d", got, startStock+9)
	}
}
