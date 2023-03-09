package utils

import (
	"container/list"

	m "github.com/murphy214/mercantile"
)

type Stack struct {
	list *list.List
}

func NewStack() *Stack {
	list := list.New()
	return &Stack{list}
}

func (stack *Stack) Push(value m.TileID) {
	stack.list.PushBack(value)
}

func (stack *Stack) Pop() m.TileID {
	e := stack.list.Back()
	if e != nil {
		stack.list.Remove(e)
		return e.Value.(m.TileID)
	}
	return m.TileID{}
}

func (stack *Stack) Peak() m.TileID {
	e := stack.list.Back()
	if e != nil {
		return e.Value.(m.TileID)
	}

	return m.TileID{}
}

func (stack *Stack) Len() int {
	return stack.list.Len()
}

func (stack *Stack) Empty() bool {
	return stack.list.Len() == 0
}
