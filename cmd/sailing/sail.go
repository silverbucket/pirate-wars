package sailing

// SailSetting controls how much canvas is set.
type SailSetting int

const (
	SailFull SailSetting = iota
	SailHalf
	SailFurled
)

func (s SailSetting) Label() string {
	switch s {
	case SailFull:
		return "Full"
	case SailHalf:
		return "Half"
	case SailFurled:
		return "Furled"
	default:
		return "?"
	}
}

func (s SailSetting) Next() SailSetting {
	switch s {
	case SailFull:
		return SailHalf
	case SailHalf:
		return SailFurled
	default:
		return SailFull
	}
}

// More sets one step more canvas: furled → half → full. Full is the maximum, so
// asking for more at full canvas is a no-op rather than a wrap to furled.
func (s SailSetting) More() SailSetting {
	switch s {
	case SailFurled:
		return SailHalf
	default:
		return SailFull
	}
}

// Less sets one step less canvas: full → half → furled. Furled is the minimum.
func (s SailSetting) Less() SailSetting {
	switch s {
	case SailFull:
		return SailHalf
	default:
		return SailFurled
	}
}

func (cfg Config) SailMultiplier(s SailSetting) float64 {
	switch s {
	case SailFull:
		return cfg.SailFull
	case SailHalf:
		return cfg.SailHalf
	case SailFurled:
		return cfg.SailFurled
	default:
		return cfg.SailFull
	}
}
