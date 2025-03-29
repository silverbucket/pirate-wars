package user_action

import (
	"pirate-wars/cmd/entities"
)

type UserActionExamine struct {
	id   int
	list entities.ViewableEntities
}

func Examine() UserActionExamine {
	return UserActionExamine{id: UserActionIdExamine, list: []entities.ViewableEntity{}}
}

func (e *UserActionExamine) GetID() int {
	return e.id
}

func (e *UserActionExamine) GetFocusedEntity() entities.ViewableEntity {
	if len(e.list) == 0 {
		return entities.NewEmptyViewableEntity()
	} else {
		return e.list[0]
	}
}

func (e *UserActionExamine) FocusLeft() {
	size := len(e.list)
	if size > 1 {
		e.list = append(entities.ViewableEntities{e.list[size-1]}, e.list[:size-1]...)
	}
}

func (e *UserActionExamine) FocusRight() {
	if len(e.list) > 1 {
		e.list = append(e.list[1:], e.list[0])
	}
}

func (e *UserActionExamine) AddItem(i entities.ViewableEntity) {
	e.list = append(e.list, i)
}
