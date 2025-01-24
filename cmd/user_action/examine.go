package user_action

import "pirate-wars/cmd/common"

type UserActionExamine struct {
	id   int
	list common.ViewableEntities
}

func Examine() UserActionExamine {
	return UserActionExamine{id: UserActionIdExamine, list: []common.ViewableEntity{}}
}

func (e *UserActionExamine) GetID() int {
	return e.id
}

func (e *UserActionExamine) GetFocusedEntity() common.ViewableEntity {
	if len(e.list) == 0 {
		return common.NewEmptyViewableEntity()
	} else {
		return e.list[0]
	}
}

func (e *UserActionExamine) FocusLeft() {
	size := len(e.list)
	if size > 1 {
		e.list = append(common.ViewableEntities{e.list[size-1]}, e.list[:size-1]...)
	}
}

func (e *UserActionExamine) FocusRight() {
	if len(e.list) > 1 {
		e.list = append(e.list[1:], e.list[0])
	}
}

func (e *UserActionExamine) AddItem(i common.ViewableEntity) {
	e.list = append(e.list, i)
}
