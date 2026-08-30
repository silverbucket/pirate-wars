package economy

import (
	"testing"
)

func TestClockWrapsAtTicksPerDay(t *testing.T) {
	cfg := DefaultConfig()
	clock := NewClock(cfg.TicksPerDay)
	atStart := clock.TimeOfDay()

	for i := 0; i < cfg.TicksPerDay; i++ {
		clock.Tick()
	}
	afterFullDay := clock.TimeOfDay()

	if atStart != afterFullDay {
		t.Fatalf("time should wrap after %d ticks: start %s != after full day %s", cfg.TicksPerDay, atStart, afterFullDay)
	}
	if clock.CurrentTick() != cfg.TicksPerDay {
		t.Fatalf("tick = %d, want %d", clock.CurrentTick(), cfg.TicksPerDay)
	}
}

func TestCargoCapacityTwenty(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.CargoCapacity != 20 {
		t.Fatalf("cargo capacity = %d, want 20", cfg.CargoCapacity)
	}
	cargo := NewCargoHold(cfg.CargoCapacity)
	if cargo.Add(GoodRum, 15) != 15 {
		t.Fatal("expected first add to succeed")
	}
	if cargo.Add(GoodPowder, 10) != 5 {
		t.Fatalf("expected only 5 free slots, got add result wrong total=%d", cargo.Total())
	}
	if cargo.Total() != 20 {
		t.Fatalf("total cargo = %d, want 20", cargo.Total())
	}
}

func TestBuySellUpdatesGoldAndStock(t *testing.T) {
	cfg := DefaultConfig()
	market := NewTownMarket(cfg, nil)
	market.stock[GoodRum] = 10
	market.prices[GoodRum] = cfg.RumPrice
	cargo := NewCargoHold(cfg.CargoCapacity)
	gold := cfg.StartingGold

	bought := BuyFromTown(&market, &cargo, &gold, GoodRum, 3)
	if bought != 3 || gold != cfg.StartingGold-3*cfg.RumPrice || market.Stock(GoodRum) != 7 {
		t.Fatalf("buy mismatch: bought=%d gold=%d stock=%d", bought, gold, market.Stock(GoodRum))
	}

	sold := SellToTown(&market, &cargo, &gold, GoodRum, 2)
	if sold != 2 || gold != cfg.StartingGold-1*cfg.RumPrice || market.Stock(GoodRum) != 9 {
		t.Fatalf("sell mismatch: sold=%d gold=%d stock=%d", sold, gold, market.Stock(GoodRum))
	}
}

func TestLoadConfigTicksPerDay(t *testing.T) {
	cfg := LoadConfig("economy.cfg")
	if cfg.TicksPerDay != 240 {
		t.Fatalf("ticks_per_day = %d, want 240", cfg.TicksPerDay)
	}
	if cfg.StartingGold != 50 {
		t.Fatalf("starting_gold = %d, want 50", cfg.StartingGold)
	}
}
