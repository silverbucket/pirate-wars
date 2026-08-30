package sailing

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds tunable sailing parameters loaded from sailing.cfg.
type Config struct {
	TickMS            int
	HullSpeed         float64
	WindDriftTicks    int
	WindStrengthMin   int
	WindStrengthMax   int
	WindFactorLight   float64
	WindFactorFresh   float64
	WindFactorStrong  float64
	NPCSkipPercent    int
	PointOfSailRun    float64
	PointOfSailBroad  float64
	PointOfSailBeam   float64
	PointOfSailClose  float64
	PointOfSailIrons  float64
	SailFull          float64
	SailHalf          float64
	SailFurled        float64
}

// DefaultConfig returns designer-approved defaults when sailing.cfg is missing.
func DefaultConfig() Config {
	return Config{
		TickMS:           250,
		HullSpeed:        0.55,
		WindDriftTicks:   100,
		WindStrengthMin:  1,
		WindStrengthMax:  3,
		WindFactorLight:  0.55,
		WindFactorFresh:  1.0,
		WindFactorStrong: 1.25,
		NPCSkipPercent:   0,
		PointOfSailRun:   1.0,
		PointOfSailBroad: 0.9,
		PointOfSailBeam:  0.75,
		PointOfSailClose: 0.4,
		PointOfSailIrons: 0.05,
		SailFull:         1.0,
		SailHalf:         0.5,
		SailFurled:       0.0,
	}
}

// WindStrengthFactor maps wind strength 0–3 to a movement multiplier.
func (c Config) WindStrengthFactor(strength int) float64 {
	switch strength {
	case 0:
		return 0
	case 1:
		return c.WindFactorLight
	case 2:
		return c.WindFactorFresh
	case 3:
		return c.WindFactorStrong
	default:
		if strength < 0 {
			return 0
		}
		return c.WindFactorStrong
	}
}

func (c Config) TickDuration() time.Duration {
	return time.Duration(c.TickMS) * time.Millisecond
}

// LoadConfig reads sailing.cfg from path, falling back to defaults for missing keys.
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
	case "tick_ms":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.TickMS = n
		}
	case "hull_speed":
		if f, err := strconv.ParseFloat(val, 64); err == nil && f >= 0 {
			cfg.HullSpeed = f
		}
	case "wind_drift_ticks":
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			cfg.WindDriftTicks = n
		}
	case "wind_strength_min":
		if n, err := strconv.Atoi(val); err == nil && n >= 0 {
			cfg.WindStrengthMin = n
		}
	case "wind_strength_max":
		if n, err := strconv.Atoi(val); err == nil && n >= 0 {
			cfg.WindStrengthMax = n
		}
	case "npc_skip_percent":
		if n, err := strconv.Atoi(val); err == nil && n >= 0 && n <= 100 {
			cfg.NPCSkipPercent = n
		}
	case "wind_factor_light":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.WindFactorLight = f
		}
	case "wind_factor_fresh":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.WindFactorFresh = f
		}
	case "wind_factor_strong":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.WindFactorStrong = f
		}
	case "point_of_sail_run":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.PointOfSailRun = f
		}
	case "point_of_sail_broad":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.PointOfSailBroad = f
		}
	case "point_of_sail_beam":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.PointOfSailBeam = f
		}
	case "point_of_sail_close":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.PointOfSailClose = f
		}
	case "point_of_sail_irons":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.PointOfSailIrons = f
		}
	case "sail_full":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.SailFull = f
		}
	case "sail_half":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.SailHalf = f
		}
	case "sail_furled":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.SailFurled = f
		}
	}
}
