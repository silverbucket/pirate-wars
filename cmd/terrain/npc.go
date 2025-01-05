package terrain

import (
	"math/rand"
)

const NUM_NPCS = 50

type ColorScheme struct {
	Foreground string
	Background string
}

var ColorPossibilities = []ColorScheme{
	{"#AA0000", "#000000"},
	{"#00AA00", "#000000"},
	{"#00AA00", "#000000"},
	{"#00FF00", "#000000"},
	{"#0000FF", "#000000"},
	{"#FFFF00", "#000000"},
	{"#FF00FF", "#000000"},
	{"#00FFFF", "#000000"},
}

var NPCs []Avatar

func (t *Terrain) CreateNPC() {
	pos := t.RandomPositionDeepWater()
	t.Logger.Infof("Creating NPC at %d, %d", pos.X, pos.Y)
	npc := CreateAvatar(pos, '⏏', ColorPossibilities[rand.Intn(len(ColorPossibilities)-1)])
	NPCs = append(NPCs, npc)
}

func (t *Terrain) InitNPCs() {
	for i := 0; i < NUM_NPCS; i++ {
		t.CreateNPC()
	}
}
