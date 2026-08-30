package common

import "testing"

func TestFacingToDeltaRoundTrip(t *testing.T) {
	for f := FacingN; f <= FacingNW; f++ {
		d := FacingToDelta(f)
		got, ok := FacingFromDelta(d.X, d.Y)
		if !ok || got != f {
			t.Fatalf("facing %v delta %+v round-trip = %v ok=%v", f, d, got, ok)
		}
	}
}

func TestFacingSeparation(t *testing.T) {
	if FacingSeparation(FacingN, FacingN) != 0 {
		t.Fatal("same facing should be 0 separation")
	}
	if FacingSeparation(FacingN, FacingS) != 4 {
		t.Fatal("opposite facings should be 4 separation")
	}
	if FacingSeparation(FacingN, FacingE) != 2 {
		t.Fatal("N and E should be beam (2)")
	}
}

func TestOppositeFacing(t *testing.T) {
	if OppositeFacing(FacingN) != FacingS {
		t.Fatal("opposite of N should be S")
	}
}
