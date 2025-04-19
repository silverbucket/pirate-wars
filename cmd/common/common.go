package common

import (
	"fmt"
	"image/color"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	LogFile        = "pirate-wars.log"
	WorldCols  int = 800 // Y
	WorldRows  int = 800 // X
	TotalTowns     = 20
	TotalNpcs      = 100
)

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

func GenID(pos Coordinates) string {
	b := letterRunes[rand.Intn(len(letterRunes))]
	return fmt.Sprintf("%v%03d%03d", string(b), pos.X, pos.Y)
}

func Inbounds(c Coordinates) bool {
	return c.X >= 0 && c.X < WorldRows && c.Y >= 0 && c.Y < WorldCols
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

func RandomPosition() Coordinates {
	return Coordinates{X: rand.Intn(WorldCols - 1), Y: rand.Intn(WorldRows - 1)}
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

func CoordsMatch(c Coordinates, p Coordinates) bool {
	if c.X == p.X && c.Y == p.Y {
		return true
	}
	return false
}

func RenderContainer(r *canvas.Rectangle, t *canvas.Text) *fyne.Container {
	t.Alignment = fyne.TextAlignCenter
	return container.NewStack(r, t)
}

func CoordToKey(c Coordinates) int {
	return c.X*WorldCols + c.Y
}

// compare colors
func ColorEqual(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}
