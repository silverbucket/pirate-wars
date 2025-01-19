package common

import (
	"fmt"
	"math/rand"
	"pirate-wars/cmd/screen"
)

const (
	LogFile     = "pirate-wars.log"
	WorldWidth  = 600 // Y
	WorldHeight = 600 // X
	TotalTowns  = 20
	TotalNpcs   = 100
)

type ViewableArea struct {
	Top    int
	Left   int
	Bottom int
	Right  int
}

type Coordinates struct {
	X int // left right
	Y int // up down
}

// Directions to explore (up, down, left, right)
var Directions = []Coordinates{
	{-1, 0},  // up
	{-1, -1}, // up left
	{-1, 1},  // up right
	{1, 0},   // down
	{1, -1},  // down left
	{1, 1},   // down right
	{0, -1},  // left
	{0, 1},   // right
}

type AvatarReadOnly interface {
	GetPos() Coordinates
	Render() string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func GenID(pos Coordinates) string {
	b := letterRunes[rand.Intn(len(letterRunes))]
	return fmt.Sprintf("%v%03d%03d", string(b), pos.X, pos.Y)
}

func Inbounds(c Coordinates) bool {
	return c.X >= 0 && c.X < WorldHeight && c.Y >= 0 && c.Y < WorldWidth
}

func IsPositionAdjacent(p Coordinates, t Coordinates) bool {
	for _, dir := range Directions {
		n := AddDirection(p, dir)
		if t.X == n.X && t.Y == n.Y {
			return true
		}
	}
	return false
}

func IsPositionWithin(c Coordinates, v ViewableArea) bool {
	if (v.Left <= c.X && c.X <= v.Right) && (v.Top <= c.Y && c.Y <= v.Bottom) {
		return true
	}
	return false
}

func RandomPosition() Coordinates {
	return Coordinates{X: rand.Intn(WorldWidth - 1), Y: rand.Intn(WorldHeight - 1)}
}

func AddDirection(p Coordinates, d Coordinates) Coordinates {
	return Coordinates{p.X + d.X, p.Y + d.Y}
}

func ClosestTo(d Coordinates, p []Coordinates) Coordinates {
	closest := Coordinates{}
	val := 99999
	for _, o := range p {
		v := diff(d.X, o.X) + diff(d.Y, o.Y)
		if v < val {
			val = v
			closest = o
		}
	}
	return closest
}

func diff(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func GetViewableArea(pos Coordinates) ViewableArea {
	// center viewport on avatar
	left := pos.X - (screen.ViewWidth / 2)
	top := pos.Y - (screen.ViewHeight / 2)
	if left < 0 {
		left = 0
	}
	if top < 0 {
		top = 0
	}
	bottom := screen.ViewHeight + top - 1
	if bottom > WorldHeight-1 {
		bottom = WorldHeight - 1
	}
	right := screen.ViewWidth + left - 1
	if right > WorldWidth-1 {
		right = WorldWidth - 1
	}
	return ViewableArea{top, left, bottom, right}
}

func GetMiniMapScale(c Coordinates) Coordinates {
	return Coordinates{c.X / screen.MiniMapFactor, c.Y / screen.MiniMapFactor}
}

func CoordsMatch(c Coordinates, p Coordinates) bool {
	if c.X == p.X && c.Y == p.Y {
		return true
	}
	return false
}
