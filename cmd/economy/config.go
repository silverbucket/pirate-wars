package economy

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Config holds tunable economy parameters loaded from economy.cfg.
type Config struct {
	TicksPerDay          int
	StartingGold         int
	CargoCapacity        int
	RumPrice             int
	PowderPrice          int
	ClothPrice           int
	RumStockMin          int
	RumStockMax          int
	PowderStockMin       int
	PowderStockMax       int
	ClothStockMin        int
	ClothStockMax        int
	SailUpgradeCost      int
	SailUpgradeHullBonus float64
	TraderCargoMin       int
	TraderCargoMax       int
	SellPercent          int
	PriceBandLowPercent  int
	PriceBandHighPercent int
	ShortStockSlack      int
}

// DefaultConfig returns designer-approved defaults when economy.cfg is missing.
func DefaultConfig() Config {
	return Config{
		TicksPerDay:          240,
		StartingGold:         50,
		CargoCapacity:        20,
		RumPrice:             5,
		PowderPrice:          8,
		ClothPrice:           6,
		RumStockMin:          10,
		RumStockMax:          30,
		PowderStockMin:       8,
		PowderStockMax:       25,
		ClothStockMin:        12,
		ClothStockMax:        28,
		SailUpgradeCost:      40,
		SailUpgradeHullBonus: 0.10,
		TraderCargoMin:       5,
		TraderCargoMax:       15,
		SellPercent:          80,
		PriceBandLowPercent:  80,
		PriceBandHighPercent: 120,
		ShortStockSlack:      2,
	}
}

// LoadConfig reads economy.cfg from path, falling back to defaults for missing keys.
func LoadConfig(path string) Config {
	cfg := DefaultConfig()
	f, err := os.Open(path)
	if err != nil {
		return cfg
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		applyConfigValue(&cfg, key, val)
	}
	return cfg
}

func applyConfigValue(cfg *Config, key, val string) {
	switch key {
	case "ticks_per_day":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.TicksPerDay = n
		}
	case "starting_gold":
		if n, err := strconv.Atoi(val); err == nil && n >= 0 {
			cfg.StartingGold = n
		}
	case "cargo_capacity":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.CargoCapacity = n
		}
	case "rum_price":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.RumPrice = n
		}
	case "powder_price":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.PowderPrice = n
		}
	case "cloth_price":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.ClothPrice = n
		}
	case "rum_stock_min":
		if n, err := strconv.Atoi(val); err == nil && n >= 0 {
			cfg.RumStockMin = n
		}
	case "rum_stock_max":
		if n, err := strconv.Atoi(val); err == nil && n >= 0 {
			cfg.RumStockMax = n
		}
	case "powder_stock_min":
		if n, err := strconv.Atoi(val); err == nil && n >= 0 {
			cfg.PowderStockMin = n
		}
	case "powder_stock_max":
		if n, err := strconv.Atoi(val); err == nil && n >= 0 {
			cfg.PowderStockMax = n
		}
	case "cloth_stock_min":
		if n, err := strconv.Atoi(val); err == nil && n >= 0 {
			cfg.ClothStockMin = n
		}
	case "cloth_stock_max":
		if n, err := strconv.Atoi(val); err == nil && n >= 0 {
			cfg.ClothStockMax = n
		}
	case "sail_upgrade_cost":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.SailUpgradeCost = n
		}
	case "sail_upgrade_hull_bonus":
		if f, err := strconv.ParseFloat(val, 64); err == nil && f > 0 {
			cfg.SailUpgradeHullBonus = f
		}
	case "trader_cargo_min":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.TraderCargoMin = n
		}
	case "trader_cargo_max":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.TraderCargoMax = n
		}
	case "sell_percent":
		if n, err := strconv.Atoi(val); err == nil && n > 0 && n <= 100 {
			cfg.SellPercent = n
		}
	case "price_band_low_percent":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.PriceBandLowPercent = n
		}
	case "price_band_high_percent":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.PriceBandHighPercent = n
		}
	case "short_stock_slack":
		if n, err := strconv.Atoi(val); err == nil && n >= 0 {
			cfg.ShortStockSlack = n
		}
	}
}

func (c Config) Price(g Good) int {
	switch g {
	case GoodRum:
		return c.RumPrice
	case GoodPowder:
		return c.PowderPrice
	case GoodCloth:
		return c.ClothPrice
	default:
		return 0
	}
}

func (c Config) StockRange(g Good) (min, max int) {
	switch g {
	case GoodRum:
		return c.RumStockMin, c.RumStockMax
	case GoodPowder:
		return c.PowderStockMin, c.PowderStockMax
	case GoodCloth:
		return c.ClothStockMin, c.ClothStockMax
	default:
		return 0, 0
	}
}
