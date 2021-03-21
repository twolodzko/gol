package parser

import (
	"fmt"
)

// List is a generic type for a list of objects of any kind
type List struct {
	List []interface{}
}

// Push adds element to the list
func (l *List) Push(obj interface{}) {
	l.List = append(l.List, obj)
}

// Print list in LISP style
func (l List) String() string {
	s := "("
	for i, elem := range l.List {
		s += fmt.Sprintf("%v", elem)
		if i < len(l.List)-1 {
			s += " "
		}
	}
	s += ")"
	return s
}

func newList(objs ...interface{}) List {
	var list List
	for _, obj := range objs {
		list.Push(obj)
	}
	return list
}

func areSame(x List, y List) (bool, error) {
	if len(x.List) != len(y.List) {
		return false, fmt.Errorf("lengths of the lists differ %d vs %d", len(x.List), len(y.List))
	}

	for i := range x.List {
		xi, yi := x.List[i], y.List[i]
		if xi != yi {
			return false, fmt.Errorf(
				"elements at position %d differ: %v (%T) vs %v (%T)", i, xi, xi, yi, yi,
			)
		}
	}

	return true, nil
}
