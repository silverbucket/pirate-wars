package common

import (
	"fmt"
	"math/rand"
)

const (
	LogFile       = "pirate-wars.log"
	WorldWidth    = 600 // Y
	WorldHeight   = 600 // X
	TotalTowns    = 20
	ViewWidth     = 75
	ViewHeight    = 50
	MiniMapFactor = 11
	TotalNpcs     = 100
)

// User Action Types
const (
	UserActionNone              = 0
	UserActionExamine           = 1
	UserActionInfo              = 2
	UserActionHelp              = 3
	UserActionMiniMap           = 4
	UserActionDebugHeatMap      = 5
	UserActionDebugViewableNpcs = 6
)

type ViewPort struct {
	width   int
	height  int
	topLeft int
}

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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func GenID(pos Coordinates, color int) string {
	b := letterRunes[rand.Intn(len(letterRunes))]
	return fmt.Sprintf("%v%03d%03d%03d", string(b), pos.X, pos.Y, color)
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
	if (v.Left < c.X && c.X < v.Right) && (v.Top < c.Y && c.Y < v.Bottom) {
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
	left := pos.X - (ViewWidth / 2)
	top := pos.Y - (ViewHeight / 2)
	if left < 0 {
		left = 0
	}
	if top < 0 {
		top = 0
	}
	return ViewableArea{top, left, ViewHeight + top - 1, ViewWidth + left - 1}
}
