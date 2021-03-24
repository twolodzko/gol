package objects

import (
	"fmt"
)

// Object has a value and is printable
type Object interface {
	String() string
	Value() interface{}
}

// Symbol is a generic type for a named object
type Symbol struct {
	Name string
}

func (s Symbol) String() string {
	return s.Name
}

// Value of a Symbol is it's name
func (s Symbol) Value() string {
	return s.Name
}

// String is a custom string type
type String struct {
	Val string
}

func (s String) String() string {
	return fmt.Sprintf("\"%s\"", s.Val)
}

// Value of a String
func (s String) Value() string {
	return s.Val
}

// Int is a custom int type
type Int struct {
	Val int
}

func (i Int) String() string {
	return fmt.Sprintf("%v", i.Val)
}

// Value of an Int
func (i Int) Value() int {
	return i.Val
}

// Float is a custom float type
type Float struct {
	Val float64
}

func (f Float) String() string {
	return fmt.Sprintf("%v", f.Val)
}

// Value of a Float
func (f Float) Value() float64 {
	return f.Val
}

// List is a generic type for a list
type List struct {
	List []interface{}
}

// Push adds element to the list
func (l *List) Push(obj interface{}) {
	l.List = append(l.List, obj)
}

// Print the List in LISP style
func (l List) String() string {
	s := ""
	for i, elem := range l.List {
		s += fmt.Sprintf("%v", elem)
		if i < len(l.List)-1 {
			s += " "
		}
	}
	return "(" + s + ")"
}

// Value of a List is a slice of it's elements
func (l List) Value() []interface{} {
	return l.List
}

// NewList initialize a List object
func NewList(objs ...interface{}) List {
	var l List
	for _, obj := range objs {
		l.Push(obj)
	}
	return l
}
