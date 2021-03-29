package objects

import (
	"fmt"
)

// Object has a value and is printable
type Object interface {
	String() string
}

// Symbol is a generic type for a named object
type Symbol struct {
	Val string
}

func (s Symbol) String() string {
	return s.Val
}

// String is a custom string type
type String struct {
	Val string
}

func (s String) String() string {
	return fmt.Sprintf("\"%s\"", s.Val)
}

// Int is a custom int type
type Int struct {
	Val int
}

func (i Int) String() string {
	return fmt.Sprintf("%v", i.Val)
}

// Float is a custom float type
type Float struct {
	Val float64
}

func (f Float) String() string {
	return fmt.Sprintf("%v", f.Val)
}

// List is a generic type for a list
type List struct {
	Val []Object
}

// Push adds element to the list
func (l *List) Push(obj Object) {
	l.Val = append(l.Val, obj)
}

// Print the List in LISP style
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

// NewList initialize a List object
func NewList(objs ...Object) List {
	var l List
	for _, obj := range objs {
		l.Push(obj)
	}
	return l
}
