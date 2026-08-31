package world

import (
	"image"
	"image/color"
	"math"
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/window"
)

// stubShip is a minimal AvatarReadOnly for exercising the glide smoother.
type stubShip struct {
	id     string
	pos    common.Coordinates
	prev   common.Coordinates
	facing common.Facing
	speed  float64
	moved  bool
}

func (s *stubShip) GetID() string                                  { return s.id }
func (s *stubShip) GetPos() common.Coordinates                     { return s.pos }
func (s *stubShip) GetPreviousPos() common.Coordinates             { return s.prev }
func (s *stubShip) GetFacing() common.Facing                       { return s.facing }
func (s *stubShip) GetTileImage() image.Image                      { return nil }
func (s *stubShip) GetTileImageFacing(f common.Facing) image.Image { return nil }
func (s *stubShip) GetViewableRange() window.Dimensions            { return window.Dimensions{} }
func (s *stubShip) GetLastSpeed() float64                          { return s.speed }
func (s *stubShip) IsHighlighted() bool                            { return false }
func (s *stubShip) GetColor() color.Color                          { return color.White }
func (s *stubShip) MovedThisTick() bool                            { return s.moved }

// TestGlideIsFrameContinuous drives the smoother through a realistic sailing
// pattern — speed below one cell per tick, so cell steps land irregularly —
// and asserts the on-screen position advances a little every frame instead of
// gliding one cell and freezing until the next step.
func TestGlideIsFrameContinuous(t *testing.T) {
	for _, speed := range []float64{0.3, 0.45, 0.6, 0.69} {
		testGlideAtSpeed(t, speed)
	}
}

func testGlideAtSpeed(t *testing.T, speed float64) {
	t.Helper()
	const (
		tickSecs  = 0.25
		frameSecs = 1.0 / 60
		simSecs   = 8.0
	)

	world := &MapView{}
	ship := &stubShip{id: "glide-test", pos: common.Coordinates{X: 10, Y: 10}, speed: speed}
	motion := Motion{Time: 0, TickSeconds: tickSecs}

	// Warm the smoother, then sail.
	world.visualPos(ship, motion)

	accum := 0.0
	nextTick := tickSecs
	lastX := 10.0
	var stalled, frames int
	maxStep := 0.0
	started := false

	for tm := frameSecs; tm < simSecs; tm += frameSecs {
		if tm >= nextTick {
			nextTick += tickSecs
			accum += speed
			if accum >= 1 {
				accum--
				ship.pos.X++
			}
		}
		motion.Time = tm
		v := world.visualPos(ship, motion)
		x, y := v.x, v.y
		if y != 10 {
			t.Fatalf("ship drifted off course: y = %v", y)
		}
		delta := x - lastX
		lastX = x
		if delta > 1e-9 {
			started = true
		}
		// Frames before the first step lands are a genuinely resting ship,
		// not jerkiness; count stalls only once under way.
		if started {
			frames++
			if delta < 1e-9 {
				stalled++
			}
		}
		if delta > maxStep {
			maxStep = delta
		}
		// The bound tolerates the momentary spike on the frame a whole-cell
		// step lands; sustained trail distance stays around 1.7 cells.
		if lag := float64(ship.pos.X) - x; lag > 2.6 {
			t.Fatalf("speed %.2f: visual lags logical by %.2f cells at t=%.2fs", speed, lag, tm)
		}
	}

	// The fastest possible trail segment is one cell in one tick, played at
	// up to the 1.35x catch-up rate; anything beyond that per-frame means a
	// teleport slipped through.
	if maxStep > frameSecs/tickSecs*1.35*1.1 {
		t.Fatalf("speed %.2f: largest per-frame step = %.3f cells; motion is jumping, not gliding", speed, maxStep)
	}
	// The old tick-lerp glided one cell then froze until the next step —
	// stalling well over a third of all frames at these speeds. The cursor
	// never dead-stops now, so under way the hull should move essentially
	// every frame.
	if ratio := float64(stalled) / float64(frames); ratio > 0.02 {
		t.Fatalf("speed %.2f: stationary on %.0f%% of frames under way; motion is stop-and-go", speed, ratio*100)
	}
	if math.Abs(lastX-float64(ship.pos.X)) > 2.0 {
		t.Fatalf("speed %.2f: visual x = %.2f never caught up to logical x = %d", speed, lastX, ship.pos.X)
	}
}

// TestGlideCoastsToRestWhenBecalmed: a ship that loses its wind mid-glide
// still finishes settling into its cell rather than freezing between squares.
func TestGlideCoastsToRestWhenBecalmed(t *testing.T) {
	world := &MapView{}
	ship := &stubShip{id: "becalmed", pos: common.Coordinates{X: 5, Y: 5}, speed: 0.5}
	motion := Motion{Time: 0, TickSeconds: 0.25}
	world.visualPos(ship, motion)

	ship.pos.X = 6 // the step lands...
	ship.speed = 0 // ...and the wind dies immediately after

	for tm := 0.016; tm < 3.0; tm += 0.016 {
		motion.Time = tm
		world.visualPos(ship, motion)
	}
	if x := world.visualPos(ship, Motion{Time: 3.1, TickSeconds: 0.25}).x; x != 6 {
		t.Fatalf("becalmed ship stuck between cells at x = %.2f", x)
	}
}

// TestGlideTurnsAtTheCorner: when the helm turns, the logical facing flips
// immediately but the hull is still traversing its old segment. Drawing the
// new facing early makes the ship skid sideways — the visual heading must
// follow the trail and swing only where the path actually bends. And it must
// swing: a continuous sweep through the corner, never a snap between sprites.
func TestGlideTurnsAtTheCorner(t *testing.T) {
	const (
		tickSecs  = 0.25
		frameSecs = 1.0 / 60
	)
	world := &MapView{}
	ship := &stubShip{
		id:     "corner",
		pos:    common.Coordinates{X: 10, Y: 10},
		facing: common.FacingE,
		speed:  0.5,
	}
	motion := Motion{Time: 0, TickSeconds: tickSecs}
	world.visualPos(ship, motion)

	// Sail east for a few steps, then the helm turns north; two ticks later
	// the northward step lands, exactly as the tick loop would apply it.
	steps := []struct {
		at     float64
		pos    common.Coordinates
		facing common.Facing
	}{
		{0.50, common.Coordinates{X: 11, Y: 10}, common.FacingE},
		{1.00, common.Coordinates{X: 12, Y: 10}, common.FacingE},
		{1.50, common.Coordinates{X: 13, Y: 10}, common.FacingE},
		{1.75, common.Coordinates{X: 13, Y: 10}, common.FacingN}, // helm turns
		{2.00, common.Coordinates{X: 13, Y: 9}, common.FacingN},
		{2.50, common.Coordinates{X: 13, Y: 8}, common.FacingN},
		{3.00, common.Coordinates{X: 13, Y: 7}, common.FacingN},
	}
	east, north := facingAngle(common.FacingE), facingAngle(common.FacingN)
	i := 0
	prev := east
	maxSwing := 0.0
	turned := false
	for tm := frameSecs; tm < 7.0; tm += frameSecs {
		for i < len(steps) && tm >= steps[i].at {
			ship.pos = steps[i].pos
			ship.facing = steps[i].facing
			i++
		}
		motion.Time = tm
		v := world.visualPos(ship, motion)
		swing := math.Abs(arc(prev, v.heading))
		prev = v.heading
		if swing > maxSwing {
			maxSwing = swing
		}
		if !v.moving {
			continue
		}
		// Well short of the corner on an eastward segment: the hull must
		// not have begun to rotate, whatever the logical facing says.
		if v.y > 9.999 && v.x < 12.3 && math.Abs(arc(east, v.heading)) > 1e-9 {
			t.Fatalf("hull heading %.2f rad at (%.2f, %.2f) while still sailing east — skidding", v.heading, v.x, v.y)
		}
		// Well past the corner heading north: the sweep has completed.
		if v.y < 9.2 {
			if math.Abs(arc(north, v.heading)) > 0.02 {
				t.Fatalf("hull heading %.2f rad at (%.2f, %.2f) has not settled on north", v.heading, v.x, v.y)
			}
			turned = true
		}
	}
	if !turned {
		t.Fatal("ship never completed the turn")
	}
	// The sweep must be continuous: no frame may swing the hull faster than
	// the turn-rate bound. A sprite-to-sprite snap is a 45° jump — far beyond
	// what the bound allows in one frame.
	if maxSwing > turnRate*frameSecs*1.05 {
		t.Fatalf("hull swung %.3f rad in one frame (limit %.3f): the turn snaps instead of sweeping", maxSwing, turnRate*frameSecs)
	}
	if maxSwing == 0 {
		t.Fatal("hull never rotated")
	}
}

// TestGlideSwingsBowAtRest: a helm order while stationary swings the hull
// round smoothly on the spot, settling on the ordered heading.
func TestGlideSwingsBowAtRest(t *testing.T) {
	world := &MapView{}
	ship := &stubShip{id: "rest", pos: common.Coordinates{X: 4, Y: 4}, facing: common.FacingN}
	world.visualPos(ship, Motion{Time: 0, TickSeconds: 0.25})
	ship.facing = common.FacingW // an about-turn to port, three octants

	prev := facingAngle(common.FacingN)
	frames := 0
	for tm := 1.0 / 60; tm < 2.0; tm += 1.0 / 60 {
		v := world.visualPos(ship, Motion{Time: tm, TickSeconds: 0.25})
		if swing := math.Abs(arc(prev, v.heading)); swing > turnRate/60*1.05 {
			t.Fatalf("hull snapped %.3f rad in one frame", swing)
		} else if swing > 1e-9 {
			frames++
		}
		prev = v.heading
	}
	if math.Abs(arc(facingAngle(common.FacingW), prev)) > 0.01 {
		t.Fatalf("hull settled on %.2f rad, not west", prev)
	}
	// Three octants at the turn rate is ~0.6s: a visible manoeuvre, not a
	// blink.
	if frames < 30 {
		t.Fatalf("turn took only %d frames; it should be a visible swing", frames)
	}
}
