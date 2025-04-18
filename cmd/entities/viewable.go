package entities

import (
	"image/color"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/window"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type ViewableEntities []ViewableEntity

type ViewableEntity interface {
	GetFlag() string
	GetType() string
	GetName() string
	GetPos() common.Coordinates
	GetID() string
	GetForegroundColor() color.Color
	Highlight()
	Render() *fyne.Container
	GetViewableRange() window.Dimensions
}

type EmptyViewableEntity struct{}

func (e *EmptyViewableEntity) GetPos() common.Coordinates {
	return common.Coordinates{X: -1, Y: -1}
}
func (e *EmptyViewableEntity) GetID() string {
	return ""
}
func (e *EmptyViewableEntity) GetName() string {
	return ""
}
func (e *EmptyViewableEntity) GetFlag() string {
	return ""
}
func (e *EmptyViewableEntity) GetType() string {
	return ""
}
func (e *EmptyViewableEntity) GetForegroundColor() color.Color {
	return color.White
}
func (e *EmptyViewableEntity) Render() *fyne.Container {
	return container.NewWithoutLayout()
}
func (e *EmptyViewableEntity) Highlight() {}
func (e *EmptyViewableEntity) GetViewableRange() window.Dimensions {
	return window.Dimensions{Width: 20, Height: 20}
}

func NewEmptyViewableEntity() *EmptyViewableEntity {
	return &EmptyViewableEntity{}
}
