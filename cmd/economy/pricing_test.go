package economy

import "testing"

func TestSellPriceAtBuyIsEightyPercent(t *testing.T) {
	cfg := DefaultConfig()
	buy := 10
	got := SellPriceAtBuy(buy, cfg.SellPercent)
	if got != 8 {
		t.Fatalf("sell price = %d, want 8", got)
	}
}

func TestBuyPriceAtStockBands(t *testing.T) {
	cfg := DefaultConfig()
	min, max := cfg.StockRange(GoodRum)

	scarce := BuyPriceAtStock(cfg, GoodRum, min, min, max)
	plenty := BuyPriceAtStock(cfg, GoodRum, max, min, max)
	if scarce != 6 || plenty != 4 {
		t.Fatalf("scarce=%d plenty=%d, want 6 and 4 for rum base 5", scarce, plenty)
	}
}
