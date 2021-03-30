package evaluator

import (
	"errors"

	"github.com/twolodzko/goal/objects"
)

type buildin = func(objects.List) (objects.Object, error)

var buildins = map[string]buildin{
	"str":   fixedNumArgs(toString, 1),
	"int":   fixedNumArgs(toInt, 1),
	"float": fixedNumArgs(toFloat, 1),
}

func fixedNumArgs(fn func(objects.Object) (objects.Object, error), numArgs int) buildin {
	return func(o objects.List) (objects.Object, error) {
		if numArgs == o.Size() {
			return fn(o.Head())
		} else {
			return nil, errors.New("wrong number of arguments")
		}
	}
}
