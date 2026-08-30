package sailing

import (
	"testing"

	"pirate-wars/cmd/common"
)

func TestOccupantAtForHailBump(t *testing.T) {
	playerPos := common.Coordinates{X: 1, Y: 1}
	npcPos := common.Coordinates{X: 2, Y: 1}
	o := NewOccupancy(map[string]common.Coordinates{
		"p1": playerPos,
		"n1": npcPos,
	})
	id, ok := o.OccupantAt(npcPos, "p1")
	if !ok || id != "n1" {
		t.Fatalf("occupant = %q ok=%v, want n1", id, ok)
	}
	_, ok = o.OccupantAt(playerPos, "p1")
	if ok {
		t.Fatal("should not report self as occupant")
	}
}
