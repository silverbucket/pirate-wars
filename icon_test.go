package main

import "testing"

func TestWindowIconsDecodeEverySize(t *testing.T) {
	icons := windowIcons()
	if len(icons) != 3 {
		t.Fatalf("decoded %d icons, want 3", len(icons))
	}
	for i, want := range []int{32, 64, 256} {
		b := icons[i].Bounds()
		if b.Dx() != want || b.Dy() != want {
			t.Errorf("icon %d is %dx%d, want %dx%d", i, b.Dx(), b.Dy(), want, want)
		}
	}
}
