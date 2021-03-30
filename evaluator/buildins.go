package evaluator

import (
	"github.com/twolodzko/goal/objects"
)

type buildin = func([]objects.Object) (objects.Object, error)

var (
	True  = objects.Symbol{Val: "true"}
	False = objects.Symbol{Val: "false"}
)

var buildins = map[string]objects.Object{
	"true":  True,
	"false": False,
	"str":   fixedArgsFunction{toString, 1},
	"int":   fixedArgsFunction{toInt, 1},
	"float": fixedArgsFunction{toFloat, 1},
	"true?": fixedArgsFunction{isTrue, 1},
	"not":   fixedArgsFunction{notTrue, 1},
	"list":  anyArgsFunction{list},
	"quote": anyArgsFunction{nil},
}

func isTrue(expr []objects.Object) (objects.Object, error) {
	// FIXME: need to evaluate the symbols!
	if expr[0] == False {
		return False, nil
	}
	return True, nil
}

func notTrue(expr []objects.Object) (objects.Object, error) {
	if expr[0] == False {
		return True, nil
	}
	return False, nil
}

func list(exprs []objects.Object) (objects.Object, error) {
	return objects.List{Val: exprs}, nil
}
