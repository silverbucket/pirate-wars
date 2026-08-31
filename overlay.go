package main

import (
	"image"
	"unicode/utf8"

	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/gfx"
	"pirate-wars/cmd/world"

	"github.com/hajimehoshi/ebiten/v2"
)

// overlayAction is one tap target inside an overlay screen.
type overlayAction struct {
	label string
	do    func()
}

// overlayRow is a line of overlay content: text, buttons, or both.
type overlayRow struct {
	text    string
	dim     bool
	heading bool
	buttons []overlayAction
}

type overlayScreen struct {
	title string
	rows  []overlayRow
}

const (
	overlayWidth   = 620
	overlayPadding = 16
	overlayRowGap  = 6
)

// overlayRowHeight is the vertical space one row occupies, gap included.
// overlayRegion is the area a modal panel may use: the map viewport, stopping
// above the action bar. The scrim is scoped to the same rectangle so the bar
// stays readable behind an open overlay — its keys are still live — and so does
// the side panel, which is where the gold and cargo a player is trading against
// are shown.
func overlayRegion() image.Rectangle {
	r := viewportRect
	if r.Max.Y > actionBarRect.Min.Y {
		r.Max.Y = actionBarRect.Min.Y
	}
	return r
}

func overlayRowHeight(r overlayRow) int {
	h := gfx.LineHeight
	if len(r.buttons) > 0 && buttonHeight > h {
		h = buttonHeight
	}
	return h + overlayRowGap
}

// overlayLayout returns the panel rect, the rows that fit inside it, and the tap
// targets for those rows.
//
// The panel height is clamped to the viewport. Laying out every row regardless
// put text and live tap targets below the panel, over the scrim and the action
// bar, for any screen taller than the clamp — a long hail, or the merchant list
// on a small window. Rows past the clamp are dropped from both the drawing and
// the hit testing, so what is clickable is always what is visible.
func (gs *GameState) overlayLayout(s overlayScreen) (image.Rectangle, []button) {
	panel, rows := gs.overlayPanel(s)

	var buttons []button
	y := panel.Min.Y + overlayPadding + gfx.LineHeight + overlayRowGap
	for _, r := range rows {
		x := panel.Max.X - overlayPadding
		// Buttons are laid out right to left so labels keep a stable left margin.
		for i := len(r.buttons) - 1; i >= 0; i-- {
			b := r.buttons[i]
			w := gfx.TextWidth(b.label) + buttonPadding*2
			if r.text == "" && w < 200 {
				w = 200
			}
			rect := image.Rect(x-w, y, x, y+buttonHeight)
			if r.text == "" {
				// Button-only rows are centred.
				rect = image.Rect(panel.Min.X+(panel.Dx()-w)/2, y, panel.Min.X+(panel.Dx()+w)/2, y+buttonHeight)
			}
			buttons = append(buttons, button{label: b.label, rect: rect, enabled: true, action: b.do})
			x -= w + buttonGap
		}
		y += overlayRowHeight(r)
	}
	return panel, buttons
}

// overlayPanel sizes the panel and returns the leading rows that fit in it.
func (gs *GameState) overlayPanel(s overlayScreen) (image.Rectangle, []overlayRow) {
	region := overlayRegion()
	maxHeight := region.Dy() - 40
	chrome := overlayPadding*2 + gfx.LineHeight + overlayRowGap

	height := chrome
	fitted := 0
	for _, r := range s.rows {
		if height+overlayRowHeight(r) > maxHeight {
			break
		}
		height += overlayRowHeight(r)
		fitted++
	}
	rows := s.rows[:fitted]

	// A truncated screen still has to offer its way out, so the final row —
	// always the dismiss button — replaces the last row that fit.
	if fitted < len(s.rows) && fitted > 0 {
		rows = append(append([]overlayRow{}, s.rows[:fitted-1]...), s.rows[len(s.rows)-1])
	}

	panel := image.Rect(0, 0, overlayWidth, height).
		Add(image.Point{
			X: region.Min.X + (region.Dx()-overlayWidth)/2,
			Y: region.Min.Y + (region.Dy()-height)/2,
		})
	return panel, rows
}

// currentOverlayScreen returns the modal panel for the current view.
func (gs *GameState) currentOverlayScreen() overlayScreen {
	switch ViewType {
	case world.ViewTypeDock:
		return gs.dockOverlayScreen()
	case world.ViewTypeHail:
		return gs.hailOverlayScreen()
	case world.ViewTypeHelp:
		return gs.helpOverlayScreen()
	case world.ViewTypeQuitConfirm:
		return gs.quitConfirmScreen()
	default:
		return overlayScreen{}
	}
}

// drawOverlay paints the scrim, the panel, and the rows that fit inside it.
func (gs *GameState) drawOverlay(screen *ebiten.Image) {
	s := gs.currentOverlayScreen()
	if s.title == "" && len(s.rows) == 0 {
		return
	}
	fillRect(screen, overlayRegion(), colorScrim)

	panel, rows := gs.overlayPanel(s)
	_, buttons := gs.overlayLayout(s)
	fillRect(screen, panel, colorPanel)
	strokeRect(screen, panel, colorPanelEdge)

	drawText(screen, s.title, panel.Min.X+overlayPadding, panel.Min.Y+overlayPadding, colorHeading)

	y := panel.Min.Y + overlayPadding + gfx.LineHeight + overlayRowGap
	for _, r := range rows {
		if r.text != "" {
			c := colorText
			switch {
			case r.heading:
				c = colorHeading
			case r.dim:
				c = colorTextDim
			}
			drawText(screen, r.text, panel.Min.X+overlayPadding, y+(overlayRowHeight(r)-overlayRowGap-gfx.LineHeight)/2, c)
		}
		y += overlayRowHeight(r)
	}

	for _, b := range buttons {
		drawButton(screen, b)
	}
}

// hailOverlayScreen shows what a hailed ship has to say.
func (gs *GameState) hailOverlayScreen() overlayScreen {
	rows := []overlayRow{}
	for _, line := range wrapText(gs.hailData.Text(), 70) {
		rows = append(rows, overlayRow{text: line})
	}
	rows = append(rows, overlayRow{buttons: []overlayAction{
		{label: barLabelFor(hailKeyMap(), "Dismiss"), do: gs.closeHail},
	}})
	return overlayScreen{title: "Hail", rows: rows}
}

// currentExamineEntity is the focused entity while examining, else an empty one.
func (gs *GameState) currentExamineEntity() entities.ViewableEntity {
	if ViewType == world.ViewTypeExamine {
		return ExamineData.GetFocusedEntity()
	}
	return entities.NewEmptyViewableEntity()
}

// wrapText breaks s into lines of at most width characters on word boundaries.
// Widths count runes, not bytes: measuring with len() wrapped any line holding a
// non-ASCII character early, and town and captain names are generated text.
func wrapText(s string, width int) []string {
	var lines []string
	line := ""
	for _, word := range splitWords(s) {
		switch {
		case word == "\n":
			lines = append(lines, line)
			line = ""
		case line == "":
			line = word
		case utf8.RuneCountInString(line)+1+utf8.RuneCountInString(word) <= width:
			line += " " + word
		default:
			lines = append(lines, line)
			line = word
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		lines = []string{""}
	}
	return lines
}

// splitWords breaks s on whitespace, keeping newlines as their own tokens.
func splitWords(s string) []string {
	var words []string
	cur := ""
	for _, r := range s {
		switch r {
		case '\n':
			if cur != "" {
				words = append(words, cur)
				cur = ""
			}
			words = append(words, "\n")
		case ' ', '\t':
			if cur != "" {
				words = append(words, cur)
				cur = ""
			}
		default:
			cur += string(r)
		}
	}
	if cur != "" {
		words = append(words, cur)
	}
	return words
}
