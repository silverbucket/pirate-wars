package main

import (
	"pirate-wars/cmd/world"
)

// actionBarCaption is the legend line: the view name plus the heading and admin
// keys that have no button of their own. Every binding stays visible.
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
func (gs *GameState) actionBarItems() []keyItem {
	_, keyMap := actionBarContext()
	items := make([]keyItem, 0, len(keyMap))
	for _, k := range keyMap {
		if !k.isBarButton() {
			continue
		}
		if k.barVisible != nil && !k.barVisible(gs) {
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

// buildButtons returns every tap target for this frame: the action bar plus any
// open overlay screen. Ebiten hit-tests these rects each frame, so a rebuild
// never leaves a stale widget under the pointer.
func (gs *GameState) buildButtons() []button {
	items := gs.actionBarItems()
	labels := make([]string, 0, len(items))
	actions := make([]func(), 0, len(items))
	for _, k := range items {
		item := k
		labels = append(labels, item.barLabel())
		actions = append(actions, func() { item.exec(gs) })
	}
	buttons := buttonRow(actionBarRect, labels, actions)

	if ViewType == world.ViewTypeDock || ViewType == world.ViewTypeHail {
		_, overlayButtons := gs.overlayLayout(gs.currentOverlayScreen())
		buttons = append(buttons, overlayButtons...)
	}
	return buttons
}
