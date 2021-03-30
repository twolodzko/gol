package evaluator

import (
	"errors"

	"github.com/twolodzko/goal/objects"
)

type buildin = func(objects.List) (objects.Object, error)

var (
	True  = objects.Symbol{Val: "true"}
	False = objects.Symbol{Val: "false"}
)

var buildins = map[string]buildin{
	"str":   fixedNumArgs(toString, 1),
	"int":   fixedNumArgs(toInt, 1),
	"float": fixedNumArgs(toFloat, 1),
	"true?": fixedNumArgs(isTrue, 1),
	"not":   fixedNumArgs(notTrue, 1),
	"list":  func(o objects.List) (objects.Object, error) { return o, nil },
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

func isTrue(o objects.Object) (objects.Object, error) {
	// FIXME: need to evaluate the symbols!
	if o == False {
		return False, nil
	}
	return True, nil
}

func notTrue(o objects.Object) (objects.Object, error) {
	if o == False {
		return True, nil
	}
	return False, nil
}
