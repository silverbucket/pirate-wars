package common

// User Action Types
const (
	UserActionIdNone              = 0
	UserActionIdExamine           = 1
	UserActionIdInfo              = 2
	UserActionIdHelp              = 3
	UserActionIdMiniMap           = 4
	UserActionIdDebugHeatMap      = 5
	UserActionIdDebugViewableNpcs = 6
)

type UserActionExamine struct {
	id   int
	list ViewableEntities
}

func (e *UserActionExamine) GetID() int {
	return e.id
}

func (e *UserActionExamine) GetFocusedEntity() ViewableEntity {
	if len(e.list) == 0 {
		return NewEmptyViewableEntity()
	} else {
		return e.list[0]
	}
}

func (e *UserActionExamine) FocusLeft() {
	size := len(e.list)
	if size > 1 {
		e.list = append(e.list[:size-1], e.list[size-1])
	}
}

func (e *UserActionExamine) FocusRight() {
	if len(e.list) > 1 {
		e.list = append(e.list[1:], e.list[0])
	}
}

func NewUserActionExamine() UserActionExamine {
	return UserActionExamine{id: UserActionIdExamine, list: []ViewableEntity{}}
}

func (e *UserActionExamine) AddItem(i ViewableEntity) {
	e.list = append(e.list, i)
}
