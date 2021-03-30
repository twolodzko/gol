package evaluator

import (
	"errors"

	"github.com/twolodzko/goal/objects"
)

type buildin = func([]objects.Object) (objects.Object, error)

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
	"list":  func(exprs []objects.Object) (objects.Object, error) { return objects.List{Val: exprs}, nil },
}

func fixedNumArgs(fn func(objects.Object) (objects.Object, error), numArgs int) buildin {
	return func(exprs []objects.Object) (objects.Object, error) {
		if numArgs == len(exprs) {
			return fn(exprs[0])
		} else {
			return nil, errors.New("wrong number of arguments")
		}
	}
}

func isTrue(expr objects.Object) (objects.Object, error) {
	// FIXME: need to evaluate the symbols!
	if expr == False {
		return False, nil
	}
	return True, nil
}

func notTrue(expr objects.Object) (objects.Object, error) {
	if expr == False {
		return True, nil
	}
	return False, nil
}
