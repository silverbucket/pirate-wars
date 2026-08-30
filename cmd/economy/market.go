package economy

import "math/rand"

// TownMarket holds stock and buy prices for a town.
type TownMarket struct {
	stock  [3]int
	prices [3]int // buy prices, recalculated from stock
}

func NewTownMarket(cfg Config, rng *rand.Rand) TownMarket {
	if rng == nil {
		rng = rand.New(rand.NewSource(1))
	}
	m := TownMarket{}
	for _, g := range AllGoods {
		min, max := cfg.StockRange(g)
		if max < min {
			max = min
		}
		stock := min
		if max > min {
			stock = min + rng.Intn(max-min+1)
		}
		m.stock[g] = stock
	}
	m.RecalcBuyPrices(cfg)
	return m
}

func (m *TownMarket) Stock(g Good) int {
	return m.stock[g]
}

// BuyPrice returns the current buy price (player buys from town).
func (m *TownMarket) BuyPrice(g Good) int {
	return m.prices[g]
}

// Price is an alias for BuyPrice.
func (m *TownMarket) Price(g Good) int {
	return m.BuyPrice(g)
}

// SellPrice returns what the town pays per unit when the player sells.
func (m *TownMarket) SellPrice(g Good, cfg Config) int {
	return SellPriceAtBuy(m.BuyPrice(g), cfg.SellPercent)
}

func (m *TownMarket) RecalcBuyPrices(cfg Config) {
	for _, g := range AllGoods {
		min, max := cfg.StockRange(g)
		m.prices[g] = BuyPriceAtStock(cfg, g, m.stock[g], min, max)
	}
}

// SetStock sets stock for a good and recalculates buy prices.
func (m *TownMarket) SetStock(g Good, stock int, cfg Config) {
	if stock < 0 {
		stock = 0
	}
	m.stock[g] = stock
	m.RecalcBuyPrices(cfg)
}

func (m *TownMarket) AddStock(g Good, amount int) {
	if amount > 0 {
		m.stock[g] += amount
	}
}

func (m *TownMarket) RemoveStock(g Good, amount int) int {
	if amount <= 0 {
		return 0
	}
	have := m.stock[g]
	if have <= 0 {
		return 0
	}
	removed := amount
	if removed > have {
		removed = have
	}
	m.stock[g] -= removed
	return removed
}

// BuyFromTown moves good from town stock to player cargo; returns units bought.
func BuyFromTown(market *TownMarket, cargo *CargoHold, gold *int, cfg Config, g Good, qty int) int {
	if qty <= 0 || gold == nil {
		return 0
	}
	price := market.BuyPrice(g)
	if price <= 0 {
		return 0
	}
	affordable := *gold / price
	if affordable <= 0 {
		return 0
	}
	want := qty
	if want > affordable {
		want = affordable
	}
	available := market.Stock(g)
	if want > available {
		want = available
	}
	added := cargo.Add(g, want)
	if added <= 0 {
		return 0
	}
	market.RemoveStock(g, added)
	*gold -= added * price
	market.RecalcBuyPrices(cfg)
	return added
}

// SellToTown moves good from player cargo to town stock; returns units sold.
func SellToTown(market *TownMarket, cargo *CargoHold, gold *int, cfg Config, g Good, qty int) int {
	if qty <= 0 || gold == nil {
		return 0
	}
	removed := cargo.Remove(g, qty)
	if removed <= 0 {
		return 0
	}
	sellPrice := market.SellPrice(g, cfg)
	market.AddStock(g, removed)
	*gold += removed * sellPrice
	market.RecalcBuyPrices(cfg)
	return removed
}
