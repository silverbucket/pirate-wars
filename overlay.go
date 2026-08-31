package main

import (
	"image"

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

func overlayRowHeight(r overlayRow) int {
	h := gfx.LineHeight
	if len(r.buttons) > 0 && buttonHeight > h {
		h = buttonHeight
	}
	return h + overlayRowGap
}

// overlayLayout returns the panel rect and the tap targets for the current screen.
func (gs *GameState) overlayLayout(s overlayScreen) (image.Rectangle, []button) {
	height := overlayPadding*2 + gfx.LineHeight + overlayRowGap
	for _, r := range s.rows {
		height += overlayRowHeight(r)
	}
	if max := viewportRect.Dy() - 40; height > max {
		height = max
	}

	panel := image.Rect(0, 0, overlayWidth, height).
		Add(image.Point{
			X: viewportRect.Min.X + (viewportRect.Dx()-overlayWidth)/2,
			Y: viewportRect.Min.Y + (viewportRect.Dy()-height)/2,
		})

	var buttons []button
	y := panel.Min.Y + overlayPadding + gfx.LineHeight + overlayRowGap
	for _, r := range s.rows {
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
			buttons = append(buttons, button{label: b.label, rect: rect, action: b.do})
			x -= w + buttonGap
		}
		y += overlayRowHeight(r)
	}
	return panel, buttons
}

func (gs *GameState) currentOverlayScreen() overlayScreen {
	switch ViewType {
	case world.ViewTypeDock:
		return gs.dockOverlayScreen()
	case world.ViewTypeHail:
		return gs.hailOverlayScreen()
	default:
		return overlayScreen{}
	}
}

func (gs *GameState) drawOverlay(screen *ebiten.Image) {
	s := gs.currentOverlayScreen()
	if s.title == "" && len(s.rows) == 0 {
		return
	}
	fillRect(screen, viewportRect, colorScrim)

	panel, buttons := gs.overlayLayout(s)
	fillRect(screen, panel, colorPanel)
	strokeRect(screen, panel, colorPanelEdge)

	drawText(screen, s.title, panel.Min.X+overlayPadding, panel.Min.Y+overlayPadding, colorHeading)

	y := panel.Min.Y + overlayPadding + gfx.LineHeight + overlayRowGap
	for _, r := range s.rows {
		if r.text != "" {
			c := colorText
			if r.dim {
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

func (gs *GameState) currentExamineEntity() entities.ViewableEntity {
	if ViewType == world.ViewTypeExamine {
		return ExamineData.GetFocusedEntity()
	}
	return entities.NewEmptyViewableEntity()
}

// wrapText breaks s into lines of at most width characters on word boundaries.
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
		case len(line)+1+len(word) <= width:
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
