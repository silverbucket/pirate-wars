package common

type ViewableEntities []ViewableEntity

type ViewableEntity interface {
	GetPos() Coordinates
	GetID() string
	SetID(string)
	GetForegroundColor() string
	GetBackgroundColor() string
	SetBackgroundColor(string)
	Render() string
}

type EmptyViewableEntity struct{}

func (e *EmptyViewableEntity) GetPos() Coordinates {
	return Coordinates{-1, -1}
}
func (e *EmptyViewableEntity) GetID() string {
	return ""
}
func (e *EmptyViewableEntity) SetID(s string) {
}
func (e *EmptyViewableEntity) GetForegroundColor() string {
	return ""
}
func (e *EmptyViewableEntity) Render() string { return "" }
func (e *EmptyViewableEntity) GetBackgroundColor() string {
	return ""
}
func (e *EmptyViewableEntity) SetBackgroundColor(string) {}

func NewEmptyViewableEntity() *EmptyViewableEntity {
	return &EmptyViewableEntity{}
}
