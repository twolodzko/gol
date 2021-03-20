package parser

import (
	"fmt"
	"reflect"
)

type List struct {
	list []interface{}
}

func (l *List) Push(obj interface{}) {
	l.list = append(l.list, obj)
}

func NewList(objs ...interface{}) List {
	var list List
	for _, obj := range objs {
		list.Push(obj)
	}
	return list
}

func areListsSame(x List, y List) (bool, error) {
	if len(x.list) != len(y.list) {
		return false, fmt.Errorf("lengths of the lists differ %d vs %d", len(x.list), len(y.list))
	}

	for i := range x.list {
		xi, yi := x.list[i], y.list[i]
		if xi != yi {
			return false, fmt.Errorf(
				"elements at position %d differ: %v (%s) vs %v (%s)",
				i, xi, reflect.TypeOf(xi).Kind(), yi, reflect.TypeOf(yi).Kind(),
			)
		}
	}

	return true, nil
}
