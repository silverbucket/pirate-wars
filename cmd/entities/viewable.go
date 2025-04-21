package entities

import (
	"image"
	"image/color"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/window"
)

type ViewableEntities []ViewableEntity

type ViewableEntity interface {
	GetFlag() string
	GetType() string
	GetName() string
	GetPos() common.Coordinates
	GetPreviousPos() common.Coordinates
	GetID() string
	GetColor() color.Color
	Highlight(b bool)
	IsHighlighted() bool
	GetTileImage() image.Image
	GetViewableRange() window.Dimensions
}

type EmptyViewableEntity struct{}

func (e *EmptyViewableEntity) GetPos() common.Coordinates {
	return common.Coordinates{X: -1, Y: -1}
}
func (e *EmptyViewableEntity) GetPreviousPos() common.Coordinates {
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
func (e *EmptyViewableEntity) GetColor() color.Color {
	return color.Transparent
}
func (e *EmptyViewableEntity) IsHighlighted() bool {
	return false
}
func (e *EmptyViewableEntity) GetTileImage() image.Image {
	return nil
}
func (e *EmptyViewableEntity) Highlight(b bool) {}
func (e *EmptyViewableEntity) GetViewableRange() window.Dimensions {
	return window.Dimensions{Width: 20, Height: 20}
}

func NewEmptyViewableEntity() *EmptyViewableEntity {
	return &EmptyViewableEntity{}
}
