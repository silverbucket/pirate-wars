package main

import (
	"pirate-wars/cmd/world"
)

// actionBarCaption is the legend line: the view name plus the sail presets and
// admin keys that have no button of their own. Every binding stays visible.
func (gs *GameState) actionBarCaption() string {
	title, keyMap := actionBarContext()
	caption := title
	if legend := keyMap.legend(); legend != "" {
		if caption != "" {
			caption += "   ·   "
		}
		caption += legend
	}
	return caption
}

// actionBarItems returns the commands that get their own button in the current
// view, one per action, each labelled with all of its keys.
//
// Unavailable commands stay in the list and render disabled. Dropping them
// reflowed every button to their right, so the row moved under the pointer as
// the player drifted past a town.
func (gs *GameState) actionBarItems() []keyItem {
	_, keyMap := actionBarContext()
	items := make([]keyItem, 0, len(keyMap))
	for _, k := range keyMap {
		if !k.isBarButton() {
			continue
		}
		items = append(items, k)
	}
	return items
}

// actionBarLabels is the bar's button text, in order. Tests read this.
func (gs *GameState) actionBarLabels() []string {
	items := gs.actionBarItems()
	labels := make([]string, 0, len(items))
	for _, k := range items {
		labels = append(labels, k.barLabel())
	}
	return labels
}

// isEnabled reports whether a command can run against the current state.
func (gs *GameState) isEnabled(k keyItem) bool {
	return k.barEnabled == nil || k.barEnabled(gs)
}

// buildButtons returns every tap target for this frame: the action bar plus any
// open overlay screen. Ebiten hit-tests these rects each frame, so a rebuild
// never leaves a stale widget under the pointer.
func (gs *GameState) buildButtons() []button {
	items := gs.actionBarItems()
	specs := make([]barSpec, 0, len(items))
	for _, k := range items {
		item := k
		specs = append(specs, barSpec{
			label:          item.barLabel(),
			enabled:        gs.isEnabled(item),
			disabledReason: item.disabledReason,
			action:         func() { item.exec(gs) },
		})
	}
	buttons := buttonRow(actionBarRect, specs)

	if hasOverlay() {
		_, overlayButtons := gs.overlayLayout(gs.currentOverlayScreen())
		buttons = append(buttons, overlayButtons...)
	}
	return buttons
}

// hasOverlay reports whether the current view draws a modal panel over the map.
func hasOverlay() bool {
	switch ViewType {
	case world.ViewTypeDock, world.ViewTypeHail, world.ViewTypeHelp, world.ViewTypeQuitConfirm:
		return true
	}
	return false
}
