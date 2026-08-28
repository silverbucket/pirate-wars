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
