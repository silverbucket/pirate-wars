package player

import "pirate-wars/cmd/economy"

// Hold tracks the player's gold, cargo, and one-time upgrades.
type Hold struct {
	Gold               int
	Cargo              economy.CargoHold
	FineSailsPurchased bool
}

func NewHold(cfg economy.Config) Hold {
	return Hold{
		Gold:  cfg.StartingGold,
		Cargo: economy.NewCargoHold(cfg.CargoCapacity),
	}
}
