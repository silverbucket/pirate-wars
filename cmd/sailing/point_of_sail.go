package sailing

import "pirate-wars/cmd/common"

// PointOfSail names the angle between ship heading and wind direction.
type PointOfSail int

const (
	PointRun PointOfSail = iota
	PointBroad
	PointBeam
	PointClose
	PointIrons
)

// ClassifyPointOfSail returns the point-of-sail for a heading relative to wind blowing downwind.
func ClassifyPointOfSail(heading, windFacing common.Facing) PointOfSail {
	sep := common.FacingSeparation(heading, windFacing)
	switch sep {
	case 0:
		return PointRun
	case 1:
		return PointBroad
	case 2:
		return PointBeam
	case 3:
		return PointClose
	default:
		return PointIrons
	}
}

func (cfg Config) PointOfSailMultiplier(pos PointOfSail) float64 {
	switch pos {
	case PointRun:
		return cfg.PointOfSailRun
	case PointBroad:
		return cfg.PointOfSailBroad
	case PointBeam:
		return cfg.PointOfSailBeam
	case PointClose:
		return cfg.PointOfSailClose
	case PointIrons:
		return cfg.PointOfSailIrons
	default:
		return cfg.PointOfSailIrons
	}
}

// EffectiveSpeed computes hull × sail × point-of-sail × wind strength for HUD and movement rolls.
func (cfg Config) EffectiveSpeed(heading common.Facing, sail SailSetting, wind *Wind) float64 {
	if wind == nil {
		return 0
	}
	pos := ClassifyPointOfSail(heading, wind.Facing)
	return cfg.HullSpeed *
		cfg.SailMultiplier(sail) *
		cfg.PointOfSailMultiplier(pos) *
		wind.StrengthFactor()
}
