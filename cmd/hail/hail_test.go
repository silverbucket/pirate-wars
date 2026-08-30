package hail

import (
	"strings"
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/town"
)

func TestHailPayload(t *testing.T) {
	cfg := economy.DefaultConfig()
	dest := town.NewTownForTest(common.Coordinates{X: 1, Y: 2}, cfg)

	n := &npc.Npc{}
	n.SetTestState("Black Bart", dest, economy.GoodRum, 7)

	p := PayloadFromNPC(n)
	if p.Name != "Black Bart" {
		t.Fatalf("name = %q", p.Name)
	}
	if p.Dest != dest.GetName() {
		t.Fatalf("dest = %q, want %q", p.Dest, dest.GetName())
	}
	if !strings.Contains(p.Cargo, "Rum") || !strings.Contains(p.Cargo, "7") {
		t.Fatalf("cargo = %q", p.Cargo)
	}
	text := p.Text()
	if !strings.Contains(text, "Hail:") || !strings.Contains(text, "Destination:") {
		t.Fatalf("text = %q", text)
	}
}
