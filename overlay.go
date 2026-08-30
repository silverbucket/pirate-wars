package main

import (
	"image/color"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/window"
	"pirate-wars/cmd/world"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (gs *GameState) buildOverlayShell() *fyne.Container {
	bg := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 200})
	bg.Resize(fyne.NewSize(float32(window.ViewPort.Dimensions.Width), float32(window.ViewPort.Dimensions.Height)))
	panel := container.NewVBox()
	gs.overlayPanel = panel
	gs.overlayRoot = container.NewStack(bg, container.NewCenter(panel))
	gs.overlayRoot.Hide()
	return gs.overlayRoot
}

func (gs *GameState) showOverlay() {
	if gs.overlayRoot == nil {
		return
	}
	gs.overlayPanel.Objects = []fyne.CanvasObject{gs.currentOverlayContent()}
	gs.overlayRoot.Show()
	gs.overlayRoot.Refresh()
}

func (gs *GameState) hideOverlay() {
	if gs.overlayRoot != nil {
		gs.overlayRoot.Hide()
		gs.overlayRoot.Refresh()
	}
}

func (gs *GameState) refreshOverlay() {
	if gs.overlayRoot == nil || gs.overlayRoot.Hidden {
		return
	}
	gs.overlayPanel.Objects = []fyne.CanvasObject{gs.currentOverlayContent()}
	gs.overlayPanel.Refresh()
	gs.updatePanels(gs.currentExamineEntity())
}

func (gs *GameState) currentOverlayContent() fyne.CanvasObject {
	switch ViewType {
	case world.ViewTypeDock:
		return gs.dockOverlayContent()
	case world.ViewTypeHail:
		return gs.hailOverlayContent()
	default:
		return widget.NewLabel("")
	}
}

func (gs *GameState) hailOverlayContent() fyne.CanvasObject {
	text := widget.NewLabel(gs.hailData.Text())
	text.Wrapping = fyne.TextWrapWord
	closeBtn := widget.NewButton("Dismiss", func() {
		gs.closeHail()
	})
	return container.NewVBox(widget.NewLabel("Hail"), text, closeBtn)
}

func (gs *GameState) currentExamineEntity() entities.ViewableEntity {
	if ViewType == world.ViewTypeExamine {
		return ExamineData.GetFocusedEntity()
	}
	return entities.NewEmptyViewableEntity()
}
