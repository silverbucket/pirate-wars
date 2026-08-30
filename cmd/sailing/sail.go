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
