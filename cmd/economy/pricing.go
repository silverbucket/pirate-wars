package economy

// BuyPriceAtStock returns the town buy price for a good at the given stock level.
// Scarce stock (near min) yields a higher price; plentiful stock (near max) yields a lower price.
func BuyPriceAtStock(cfg Config, g Good, stock, min, max int) int {
	base := cfg.Price(g)
	if base <= 0 {
		return 0
	}
	if max < min {
		max = min
	}
	if stock < min {
		stock = min
	}
	if stock > max {
		stock = max
	}
	if max == min {
		return scalePrice(base, cfg.PriceBandHighPercent)
	}
	ratio := float64(stock-min) / float64(max-min)
	low := float64(cfg.PriceBandLowPercent) / 100.0
	high := float64(cfg.PriceBandHighPercent) / 100.0
	factor := high + (low-high)*ratio
	return scalePrice(base, int(factor*100+0.5))
}

func scalePrice(base, percent int) int {
	if percent <= 0 {
		return base
	}
	p := base * percent / 100
	if p < 1 {
		return 1
	}
	return p
}

// SellPriceAtBuy returns what the town pays per unit when selling to them.
func SellPriceAtBuy(buyPrice, sellPercent int) int {
	if buyPrice <= 0 || sellPercent <= 0 {
		return 0
	}
	return buyPrice * sellPercent / 100
}

// IsShortOnGood reports whether stock is at or near the configured minimum.
func IsShortOnGood(cfg Config, stock int, g Good) bool {
	min, _ := cfg.StockRange(g)
	return stock <= min+cfg.ShortStockSlack
}
