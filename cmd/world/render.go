package world

import (
	"image"
	"math"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/gfx"
	"pirate-wars/cmd/resources"
	"pirate-wars/cmd/window"

	"github.com/hajimehoshi/ebiten/v2"
)

// Motion carries the clocks the renderer animates against.
type Motion struct {
	// Time is seconds since boot: the swell, wake churn and glide smoothing
	// all run on real time so they stay fluid however long a tick is.
	Time float64
	// TickSeconds is the sailing tick length, which converts a ship's speed
	// (cells per tick) into an on-screen rate (cells per second).
	TickSeconds float64
}

// Heading sweep tuning. A ship's on-screen heading is a continuous angle
// that carves through each turn rather than pivoting between the eight hull
// sprites at the cell boundary.
const (
	// octantSweepSecs is how long a one-octant (45°) turn takes, centred on
	// the corner; wider turns scale up proportionally.
	octantSweepSecs = 0.3
	// maxSweepSecs caps the widest turn (an about-face) so it still reads
	// as one decisive manoeuvre.
	maxSweepSecs = 1.0
	// turnRate bounds how fast the hull can swing when the desired heading
	// changes abruptly — a turn ordered while at rest, or a corner the trail
	// only learned about late. Radians per second.
	turnRate = 4.0
	// turnEase is the time constant of the ease-out onto the desired
	// heading, so a swing settles instead of stopping dead.
	turnEase = 0.05
)

// glideState is a ship's renderer-side motion trail: the cells it stepped
// through, stamped with when the renderer saw each step land.
type glideState struct {
	path    []crumb
	gap     float64 // smoothed seconds between recent steps
	lastTr  float64 // monotonic render-time cursor
	lastNow float64 // wall time of the previous frame, for rate limiting
	heading float64 // on-screen heading, radians clockwise from north
	// prevAngle is the heading of the segment most recently dropped off the
	// front of the trail, so the sweep through that corner can finish after
	// the cursor has passed it.
	prevAngle float64
	hasPrev   bool
}

type crumb struct {
	x, y float64
	t    float64
}

// shipView is one ship resolved for this frame: where to draw it, the heading
// its hull shows, and whether it is visibly under way. The heading follows the
// smoothed trail, not the logical facing — the logical facing flips the moment
// the helm turns, while the hull is still finishing its old segment, and
// drawing it rotated early looks like the ship is skidding sideways.
type shipView struct {
	ship    entities.AvatarReadOnly
	x, y    float64
	heading float64 // radians clockwise from north
	moving  bool
}

// octant is the nearest of the eight sprite facings to the heading.
func (v shipView) octant() common.Facing {
	n := int(math.Round(v.heading/(math.Pi/4))) % 8
	if n < 0 {
		n += 8
	}
	return common.Facing(n)
}

// tilt is the residual rotation from the octant sprite to the true heading,
// in [-π/8, π/8], applied on the GPU when the sprite is drawn.
func (v shipView) tilt() float64 {
	return arc(facingAngle(v.octant()), v.heading)
}

// direction is the unit vector the hull points along, in screen space.
func (v shipView) direction() (dx, dy float64) {
	return math.Sin(v.heading), -math.Cos(v.heading)
}

// Draw renders the tilemap viewport centred on the player's smoothed position
// into dst. The camera works in world pixels, not cells, so it glides with
// the player instead of stepping a cell at a time.
func (world *MapView) Draw(
	dst *ebiten.Image,
	avatar entities.AvatarReadOnly,
	npcs []entities.AvatarReadOnly,
	highlight entities.ViewableEntity,
	motion Motion,
) {
	cell := float64(window.CellSize)
	vw := float64(dst.Bounds().Dx())
	vh := float64(dst.Bounds().Dy())

	views := make([]shipView, 0, len(npcs)+1)
	for _, ship := range append([]entities.AvatarReadOnly{avatar}, npcs...) {
		views = append(views, world.visualPos(ship, motion))
	}
	player := views[0]

	// Camera top-left in world pixels, clamped to the world.
	camX := clamp((player.x+0.5)*cell-vw/2, 0, float64(common.WorldCols)*cell-vw)
	camY := clamp((player.y+0.5)*cell-vh/2, 0, float64(common.WorldRows)*cell-vh)

	// Terrain pass: every cell overlapping the viewport.
	col0 := int(camX / cell)
	row0 := int(camY / cell)
	for x := col0; float64(x)*cell < camX+vw; x++ {
		for y := row0; float64(y)*cell < camY+vh; y++ {
			if x < 0 || x >= common.WorldCols || y < 0 || y >= common.WorldRows {
				continue
			}
			pos := common.Coordinates{X: x, Y: y}
			blit(dst, world.terrainTileAt(pos), float64(x)*cell-camX, float64(y)*cell-camY)
		}
	}

	highlightVisible := false
	h := highlight.GetPos()
	if h.X >= 0 {
		if !highlight.IsHighlighted() {
			highlight.Highlight(true)
		}
		_, _, _, a := highlight.GetColor().RGBA()
		highlightVisible = a > 0
	}

	// Wakes go under every hull so a ship can sail over a neighbour's trail.
	// A wake shows whenever the smoothed hull is actually in motion, which at
	// steady sail is continuous — not only on the ticks a cell step landed.
	// It trails the true heading, so it swings through a turn with the hull.
	for _, v := range views {
		if !v.moving {
			continue
		}
		dx, dy := v.direction()
		wx := (v.x - dx*0.85) * cell
		wy := (v.y - dy*0.85) * cell
		frame := int(motion.Time*3) % resources.DirectionalWakeFrames
		blitRotated(dst, resources.GetDirectionalWake(v.octant(), frame), wx-camX, wy-camY, v.tilt())
	}

	for _, v := range views {
		if v.ship.IsHighlighted() && !highlightVisible {
			continue
		}
		roll, bob := swell(v.ship, motion.Time)
		blitRotated(dst, v.ship.GetTileImageFacing(v.octant()), v.x*cell-camX, v.y*cell-camY+bob, roll+v.tilt())
	}

	if h.X >= 0 && highlight.IsHighlighted() && highlightVisible {
		hx, hy := float64(h.X)*cell, float64(h.Y)*cell
		// If the highlight sits on a ship, the ring rides its glide.
		for _, v := range views {
			if common.CoordsMatch(v.ship.GetPos(), h) {
				hx, hy = v.x*cell, v.y*cell
				break
			}
		}
		blit(dst, resources.GetExamineRingOverlay(window.CellSize), hx-camX, hy-camY)
	}

	// Player chrome rides the player's smoothed position, drawn last so
	// nothing can cover the heading readout. The bow marker shows the
	// ordered heading immediately — instant feedback that the helm answered,
	// while the hull itself carves round a beat later.
	blit(dst, resources.GetPlayerMarkerOverlay(window.CellSize), player.x*cell-camX, player.y*cell-camY)
	blit(dst, resources.GetBowMarkerOverlay(window.CellSize, avatar.GetFacing()), player.x*cell-camX, player.y*cell-camY)
}

// visualPos resolves a ship's on-screen position in fractional cells, the
// heading its hull should show, and whether it is visibly under way.
//
// Logical movement is quantised: a whole cell lands each time the speed
// accumulator fills, at uneven tick gaps. Rendering that directly is
// teleports; lerping only across the landing tick is glide-stop-spurt. So
// the renderer keeps a breadcrumb trail of the ship's recent steps and draws
// the hull a fixed beat behind real time, interpolating along the trail —
// entity interpolation, as netcode does it. Within each segment velocity is
// constant, so motion is fluid at any speed and any step cadence.
//
// The heading follows the trail while under way — the hull turns where the
// path turns, sweeping through each corner rather than pivoting on it — and
// the logical facing at rest, so a ship can still swing its bow around while
// stationary.
func (world *MapView) visualPos(ship entities.AvatarReadOnly, motion Motion) shipView {
	now := motion.Time
	tickSecs := motion.TickSeconds
	if tickSecs <= 0 {
		tickSecs = 1
	}
	p := ship.GetPos()
	px, py := float64(p.X), float64(p.Y)

	if world.glides == nil {
		world.glides = map[string]*glideState{}
	}
	st, ok := world.glides[ship.GetID()]
	if !ok {
		st = &glideState{path: []crumb{{px, py, now}}, lastTr: now, lastNow: now, heading: facingAngle(ship.GetFacing())}
		world.glides[ship.GetID()] = st
	}

	dt := now - st.lastNow
	if dt < 0 {
		dt = 0
	}
	st.lastNow = now

	// Expected seconds per cell: observed step gaps once we have them, the
	// ship's speed as the estimate before that.
	speed := clamp(ship.GetLastSpeed(), 0.15, 2)
	gap := st.gap
	if gap == 0 {
		gap = tickSecs / speed
	}
	gap = clamp(gap, tickSecs, 3)

	last := &st.path[len(st.path)-1]
	if px != last.x || py != last.y {
		if math.Hypot(px-last.x, py-last.y) > 2.5 {
			// A respawn or desync, not sailing: snap.
			st.path = []crumb{{px, py, now}}
			st.gap = 0
			st.lastTr = now
			st.hasPrev = false
			st.heading = facingAngle(ship.GetFacing())
			return shipView{ship: ship, x: px, y: py, heading: st.heading}
		}
		observed := now - last.t
		if observed > 2.2*gap {
			// Setting off from rest: restart the trail close behind the
			// step so the cast-off starts promptly instead of waiting out
			// the full interpolation delay. The cursor jumps with it —
			// invisible, since the hull sits on that crumb either way.
			last.t = now - 0.7*gap
			st.lastTr = last.t
			st.gap = 0
			st.hasPrev = false
		} else if st.gap == 0 {
			st.gap = observed
		} else {
			// Smooth the cadence estimate heavily: gaps alternate between
			// whole tick counts (2,1,2,1...), and chasing each swing makes
			// the delay dip below the longer gaps — cursor underruns, which
			// read as micro-stutters.
			st.gap = 0.7*st.gap + 0.3*observed
		}
		st.path = append(st.path, crumb{px, py, now})
	}

	// Render one gap plus a generous tick behind real time. Step gaps are
	// whole numbers of ticks, so their jitter is about a tick regardless of
	// speed — that fixed slack is what absorbs it. The cursor chases the
	// delayed clock at a smoothly bounded rate — up to 1.35x when behind,
	// easing to 0.6x when ahead, never a dead stop: even a few frozen frames
	// per step read as a stutter at low speeds, where each cadence-estimate
	// correction used to halt the cursor outright.
	target := now - gap - 1.2*tickSecs
	rate := clamp(1+(target-st.lastTr)*1.5, 0.6, 1.35)
	tr := st.lastTr + rate*dt
	st.lastTr = tr

	for len(st.path) >= 2 && st.path[1].t <= tr {
		if ang, ok := segmentAngle(st.path[0], st.path[1]); ok {
			st.prevAngle = ang
			st.hasPrev = true
		}
		st.path = st.path[1:]
	}
	a := st.path[0]
	if len(st.path) == 1 {
		// Nothing to interpolate toward: hold the cursor at the trail's
		// end, so the next crumb's segment starts from f=0 — otherwise the
		// cursor drifts past a.t while the ship rests and the hull hops
		// mid-segment the moment a step lands.
		if tr > a.t {
			tr = a.t
			st.lastTr = tr
		}
		return st.resolve(ship, a.x, a.y, facingAngle(ship.GetFacing()), false, dt)
	}
	if tr <= a.t {
		return st.resolve(ship, a.x, a.y, facingAngle(ship.GetFacing()), false, dt)
	}
	b := st.path[1]
	f := (tr - a.t) / (b.t - a.t)
	x := a.x + (b.x-a.x)*f
	y := a.y + (b.y-a.y)*f

	// The desired heading is this segment's, blended through the corners at
	// either end: finishing the swing out of the previous segment, and
	// anticipating the swing into the next one once the trail knows it.
	// Each sweep is centred on its corner so the hull carves the bend.
	desired, ok := segmentAngle(a, b)
	if !ok {
		desired = st.heading
	}
	if st.hasPrev {
		desired = sweep(st.prevAngle, desired, a.t, tr)
	}
	if len(st.path) >= 3 {
		if next, ok := segmentAngle(b, st.path[2]); ok {
			desired = sweep(desired, next, b.t, tr)
		}
	}
	return st.resolve(ship, x, y, desired, true, dt)
}

// resolve eases the on-screen heading toward the desired one and builds the
// frame's view. The ease absorbs whatever the sweep could not plan for — a
// helm order at rest, or a corner the trail learned about late — bounded to
// turnRate so it is always a swing, never a snap.
func (st *glideState) resolve(ship entities.AvatarReadOnly, x, y, desired float64, moving bool, dt float64) shipView {
	d := arc(st.heading, desired)
	step := d * (1 - math.Exp(-dt/turnEase))
	if lim := turnRate * dt; math.Abs(step) > lim {
		step = math.Copysign(lim, d)
	}
	st.heading = math.Remainder(st.heading+step, 2*math.Pi)
	return shipView{ship: ship, x: x, y: y, heading: st.heading, moving: moving}
}

// sweep blends from one heading to the next across a window centred on the
// corner time, easing in and out. Before the window it is the old heading,
// after it the new; between, the turn is under way.
func sweep(from, to, corner, tr float64) float64 {
	turn := arc(from, to)
	if turn == 0 {
		return to
	}
	w := clamp(octantSweepSecs*math.Abs(turn)/(math.Pi/4), octantSweepSecs, maxSweepSecs)
	u := clamp((tr-(corner-w/2))/w, 0, 1)
	return from + turn*u*u*(3-2*u)
}

// segmentAngle is the heading of the step from a to b.
func segmentAngle(a, b crumb) (float64, bool) {
	f, ok := common.FacingFromDelta(int(math.Round(b.x-a.x)), int(math.Round(b.y-a.y)))
	if !ok {
		return 0, false
	}
	return facingAngle(f), true
}

// facingAngle maps a sprite facing to radians clockwise from north.
func facingAngle(f common.Facing) float64 {
	return float64(f) * math.Pi / 4
}

// arc is the signed shortest rotation from one heading to another.
func arc(from, to float64) float64 {
	return math.Remainder(to-from, 2*math.Pi)
}

// swell returns the hull's roll (radians) and bob (pixels) at time t. Each
// ship gets a phase from its ID so the fleet rides the sea independently.
func swell(ship entities.AvatarReadOnly, t float64) (roll, bob float64) {
	var phase float64
	for _, c := range ship.GetID() {
		phase += float64(c)
	}
	roll = 0.045 * math.Sin(t*1.4+phase)
	bob = 1.6 * math.Sin(t*2.2+phase*1.7)
	return roll, bob
}

func clamp(v, lo, hi float64) float64 {
	if hi < lo {
		return lo
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// blit draws a cached tile image at (x, y) through the Ebiten texture cache.
func blit(dst *ebiten.Image, src image.Image, x, y float64) {
	tex := gfx.Texture(src)
	if tex == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	dst.DrawImage(tex, op)
}

// blitRotated draws a sprite rotated about its own centre — the swell roll
// and the residual heading between octant sprites.
func blitRotated(dst *ebiten.Image, src image.Image, x, y, rad float64) {
	tex := gfx.Texture(src)
	if tex == nil {
		return
	}
	b := tex.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
	op.GeoM.Rotate(rad)
	op.GeoM.Translate(x+float64(b.Dx())/2, y+float64(b.Dy())/2)
	dst.DrawImage(tex, op)
}
