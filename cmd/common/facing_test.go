package common

import "testing"

func TestFacingFromDelta(t *testing.T) {
	tests := []struct {
		dx, dy int
		want   Facing
	}{
		{0, -1, FacingN},
		{1, -1, FacingNE},
		{1, 0, FacingE},
		{1, 1, FacingSE},
		{0, 1, FacingS},
		{-1, 1, FacingSW},
		{-1, 0, FacingW},
		{-1, -1, FacingNW},
	}

	for _, tc := range tests {
		got, ok := FacingFromDelta(tc.dx, tc.dy)
		if !ok {
			t.Fatalf("FacingFromDelta(%d, %d) ok = false, want true", tc.dx, tc.dy)
		}
		if got != tc.want {
			t.Fatalf("FacingFromDelta(%d, %d) = %v, want %v", tc.dx, tc.dy, got, tc.want)
		}
	}

	if _, ok := FacingFromDelta(0, 0); ok {
		t.Fatal("FacingFromDelta(0, 0) should not match a facing")
	}
}

// TestTackOneOctantEachWay locks the octant ring used by the A/D steering keys:
// clockwise is starboard, anticlockwise is port, and both wrap.
func TestTackOneOctantEachWay(t *testing.T) {
	ring := []Facing{FacingN, FacingNE, FacingE, FacingSE, FacingS, FacingSW, FacingW, FacingNW}

	for i, f := range ring {
		wantStbd := ring[(i+1)%len(ring)]
		if got := TackStarboard(f); got != wantStbd {
			t.Fatalf("TackStarboard(%s) = %s, want %s",
				FacingLabel(f), FacingLabel(got), FacingLabel(wantStbd))
		}
		wantPort := ring[(i-1+len(ring))%len(ring)]
		if got := TackPort(f); got != wantPort {
			t.Fatalf("TackPort(%s) = %s, want %s",
				FacingLabel(f), FacingLabel(got), FacingLabel(wantPort))
		}
	}

	if got := TackStarboard(FacingNW); got != FacingN {
		t.Fatalf("starboard from NW should wrap to N, got %s", FacingLabel(got))
	}
	if got := TackPort(FacingN); got != FacingNW {
		t.Fatalf("port from N should wrap to NW, got %s", FacingLabel(got))
	}
}
