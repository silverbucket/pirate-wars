package entities

import (
	"image/color"
	"pirate-wars/cmd/common"
	"testing"
)

func TestAvatarFacingFromMovement(t *testing.T) {
	start := common.Coordinates{X: 10, Y: 10}
	avatar := CreateAvatar(start, common.ShipWhite, color.White)

	avatar.SetPos(common.Coordinates{X: 11, Y: 10})
	if avatar.facing != common.FacingE {
		t.Fatalf("facing = %v, want %v", avatar.facing, common.FacingE)
	}

	avatar.SetPos(common.Coordinates{X: 11, Y: 11})
	if avatar.facing != common.FacingS {
		t.Fatalf("facing = %v, want %v", avatar.facing, common.FacingS)
	}

	// Idle should keep the last facing.
	before := avatar.image
	avatar.SetPos(common.Coordinates{X: 11, Y: 11})
	if avatar.facing != common.FacingS {
		t.Fatalf("idle facing = %v, want %v", avatar.facing, common.FacingS)
	}
	if avatar.image != before {
		t.Fatal("idle movement should not change tile image")
	}
}

func TestAvatarSetHeading(t *testing.T) {
	start := common.Coordinates{X: 10, Y: 10}
	avatar := CreateAvatar(start, common.ShipWhite, color.White)

	avatar.SetHeading(common.FacingW)
	if avatar.GetFacing() != common.FacingW {
		t.Fatalf("heading = %v, want %v", avatar.GetFacing(), common.FacingW)
	}

	before := avatar.GetPos()
	avatar.SetHeading(common.FacingW)
	if !common.CoordsMatch(avatar.GetPos(), before) {
		t.Fatal("set heading should not change position")
	}
}
func TestAvatarSpawnFacingNorth(t *testing.T) {
	avatar := CreateAvatar(common.Coordinates{X: 1, Y: 1}, common.ShipPirate, color.White)
	if avatar.facing != common.FacingN {
		t.Fatalf("spawn facing = %v, want %v", avatar.facing, common.FacingN)
	}
}
