package sailing

import (
	"math/rand"
	"pirate-wars/cmd/common"
)

// Wind is map-wide wind state: direction the wind blows (downwind) and strength 0–3.
type Wind struct {
	Facing   common.Facing
	Strength int
	ticks    int
	cfg      Config
}

func NewWind(cfg Config) *Wind {
	w := &Wind{cfg: cfg}
	w.randomize()
	return w
}

// NewFixedWind returns wind locked to a facing and strength (deterministic sims/tests).
func NewFixedWind(cfg Config, facing common.Facing, strength int) *Wind {
	return &Wind{cfg: cfg, Facing: facing, Strength: strength}
}

func (w *Wind) randomize() {
	w.Facing = common.Facing(rand.Intn(8))
	min := w.cfg.WindStrengthMin
	max := w.cfg.WindStrengthMax
	if max < min {
		max = min
	}
	if max == min {
		w.Strength = min
	} else {
		w.Strength = min + rand.Intn(max-min+1)
	}
}

// Tick advances wind drift; call once per sailing tick.
func (w *Wind) Tick() {
	w.ticks++
	if w.ticks >= w.cfg.WindDriftTicks {
		w.ticks = 0
		w.drift()
	}
}

func (w *Wind) drift() {
	// Slowly rotate facing by ±1 step around the compass.
	delta := 1
	if rand.Intn(2) == 0 {
		delta = -1
	}
	w.Facing = common.RotateFacing(w.Facing, delta)

	min := w.cfg.WindStrengthMin
	max := w.cfg.WindStrengthMax
	if max < min {
		max = min
	}
	if max > min {
		change := rand.Intn(3) - 1 // -1, 0, or +1
		w.Strength += change
		if w.Strength < min {
			w.Strength = min
		}
		if w.Strength > max {
			w.Strength = max
		}
	}
}

// StrengthFactor returns the configured movement multiplier for current wind strength.
func (w *Wind) StrengthFactor() float64 {
	return w.cfg.WindStrengthFactor(w.Strength)
}

// Label returns a short HUD/debug label for wind direction.
func (w *Wind) Label() string {
	return common.FacingLabel(w.Facing)
}
