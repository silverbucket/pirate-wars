package common

type ViewableEntities []ViewableEntity

type ViewableEntity interface {
	GetFlag() string
	GetType() string
	GetName() string
	GetPos() Coordinates
	GetID() string
	GetForegroundColor() string
	Highlight()
	Render() string
}

type EmptyViewableEntity struct{}

func (e *EmptyViewableEntity) GetPos() Coordinates {
	return Coordinates{-1, -1}
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
func (e *EmptyViewableEntity) GetForegroundColor() string {
	return ""
}
func (e *EmptyViewableEntity) Render() string { return "" }
func (e *EmptyViewableEntity) Highlight()     {}

func NewEmptyViewableEntity() *EmptyViewableEntity {
	return &EmptyViewableEntity{}
}
