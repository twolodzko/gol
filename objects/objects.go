package objects

import (
	"fmt"
)

type Object interface {
	String() string
}

type Symbol struct {
	Val string
}

func (s Symbol) String() string {
	return s.Val
}

type String struct {
	Val string
}

func (s String) String() string {
	return fmt.Sprintf("\"%s\"", s.Val)
}

type Int struct {
	Val int
}

func (i Int) String() string {
	return fmt.Sprintf("%v", i.Val)
}

type Float struct {
	Val float64
}

func (f Float) String() string {
	return fmt.Sprintf("%v", f.Val)
}

type List struct {
	Val []Object
}

func (l *List) Push(obj Object) {
	l.Val = append(l.Val, obj)
}

func (l List) String() string {
	s := ""
	for i, elem := range l.Val {
		s += fmt.Sprintf("%v", elem)
		if i < len(l.Val)-1 {
			s += " "
		}
	}
	return "(" + s + ")"
}

func NewList(objs ...Object) List {
	var l List
	for _, obj := range objs {
		l.Push(obj)
	}
	return l
}
