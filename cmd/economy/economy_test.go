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
	market.SetStock(GoodRum, 20, cfg)
	cargo := NewCargoHold(cfg.CargoCapacity)
	gold := cfg.StartingGold
	buyPrice := market.BuyPrice(GoodRum)

	bought := BuyFromTown(&market, &cargo, &gold, cfg, GoodRum, 3)
	if bought != 3 || gold != cfg.StartingGold-3*buyPrice || market.Stock(GoodRum) != 17 {
		t.Fatalf("buy mismatch: bought=%d gold=%d stock=%d", bought, gold, market.Stock(GoodRum))
	}

	sellPrice := market.SellPrice(GoodRum, cfg)
	currentBuy := market.BuyPrice(GoodRum)
	if sellPrice != currentBuy*cfg.SellPercent/100 {
		t.Fatalf("sell price = %d, want %d (80%% of buy %d)", sellPrice, currentBuy*cfg.SellPercent/100, currentBuy)
	}

	sold := SellToTown(&market, &cargo, &gold, cfg, GoodRum, 2)
	wantGold := cfg.StartingGold - 3*buyPrice + 2*sellPrice
	if sold != 2 || gold != wantGold || market.Stock(GoodRum) != 19 {
		t.Fatalf("sell mismatch: sold=%d gold=%d stock=%d", sold, gold, market.Stock(GoodRum))
	}
}

func TestSameTownBuyThenSellLosesGold(t *testing.T) {
	cfg := DefaultConfig()
	market := NewTownMarket(cfg, nil)
	market.SetStock(GoodRum, 20, cfg)
	cargo := NewCargoHold(cfg.CargoCapacity)
	gold := cfg.StartingGold
	startGold := gold

	buyPrice := market.BuyPrice(GoodRum)
	if BuyFromTown(&market, &cargo, &gold, cfg, GoodRum, 3) != 3 {
		t.Fatal("expected to buy 3 rum")
	}
	sellPrice := market.SellPrice(GoodRum, cfg)
	if SellToTown(&market, &cargo, &gold, cfg, GoodRum, 3) != 3 {
		t.Fatal("expected to sell 3 rum")
	}
	if gold >= startGold {
		t.Fatalf("round-trip should lose gold: start=%d end=%d", startGold, gold)
	}
	expectedGold := startGold - 3*buyPrice + 3*sellPrice
	if gold != expectedGold {
		t.Fatalf("gold = %d, want %d after round-trip loss", gold, expectedGold)
	}
}

func TestBuyPriceMovesWithStock(t *testing.T) {
	cfg := DefaultConfig()
	market := NewTownMarket(cfg, nil)
	min, max := cfg.StockRange(GoodRum)

	market.SetStock(GoodRum, min, cfg)
	scarce := market.BuyPrice(GoodRum)
	market.SetStock(GoodRum, max, cfg)
	plenty := market.BuyPrice(GoodRum)

	if scarce <= plenty {
		t.Fatalf("scarce buy price %d should exceed plenty buy price %d", scarce, plenty)
	}
	wantScarce := BuyPriceAtStock(cfg, GoodRum, min, min, max)
	wantPlenty := BuyPriceAtStock(cfg, GoodRum, max, min, max)
	if scarce != wantScarce || plenty != wantPlenty {
		t.Fatalf("prices scarce=%d plenty=%d, want %d and %d", scarce, plenty, wantScarce, wantPlenty)
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
	if cfg.SellPercent != 80 {
		t.Fatalf("sell_percent = %d, want 80", cfg.SellPercent)
	}
	if cfg.PriceBandLowPercent != 80 || cfg.PriceBandHighPercent != 120 {
		t.Fatalf("price bands = %d/%d, want 80/120", cfg.PriceBandLowPercent, cfg.PriceBandHighPercent)
	}
}
