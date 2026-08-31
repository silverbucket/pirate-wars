package main

import (
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/harbor"
	"pirate-wars/cmd/world"
)

// combinedWorld routes passability to the painted harbor mask inside the harbor rect.
type combinedWorld struct {
	main   *world.MapView
	harbor *harbor.World
}

func newCombinedWorld(main *world.MapView, hw *harbor.World) *combinedWorld {
	return &combinedWorld{main: main, harbor: hw}
}

func (c *combinedWorld) IsPassableByBoat(pos common.Coordinates) bool {
	if c == nil {
		return false
	}
	if harbor.InRegion(pos) && c.harbor != nil {
		return c.harbor.IsPassableByBoat(pos)
	}
	if c.main != nil {
		return c.main.IsPassableByBoat(pos)
	}
	return false
}

func (c *combinedWorld) IsDock(pos common.Coordinates) bool {
	if c == nil || c.harbor == nil {
		return false
	}
	return c.harbor.IsDock(pos)
}

func (c *combinedWorld) InHarbor(pos common.Coordinates) bool {
	return harbor.InRegion(pos)
}
