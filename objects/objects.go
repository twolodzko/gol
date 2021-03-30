package objects

import (
	"fmt"
	"strings"
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

func (l *List) Size() int {
	return len(l.Val)
}

func (l *List) Push(obj Object) {
	l.Val = append(l.Val, obj)
}

func (l List) String() string {
	var str []string
	for _, elem := range l.Val {
		str = append(str, elem.String())
	}
	return "(" + strings.Join(str, " ") + ")"
}

func NewList(objs ...Object) List {
	var l List
	for _, obj := range objs {
		l.Push(obj)
	}
	return l
}

func (l *List) Head() Object {
	return l.Val[0]
}

func (l *List) Tail() List {
	tail := l.Val[1:]
	if len(tail) > 0 {
		return List{tail}
	}
	return List{}
}
