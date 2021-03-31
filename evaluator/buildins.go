package evaluator

import (
	. "github.com/twolodzko/goal/types"
)

type buildin = func(List) (Any, error)

var buildins = map[string]Any{
	"true":  true,
	"false": false,
	"str":   fixedArgsFunction{toString, 1},
	"int":   fixedArgsFunction{toInt, 1},
	"float": fixedArgsFunction{toFloat, 1},
	// "true?": fixedArgsFunction{isTrue, 1},
	// "not":   fixedArgsFunction{notTrue, 1},
	"list":  anyArgsFunction{list},
	"quote": anyArgsFunction{nil},
}

// func isTrue(expr List) (Any, error) {
// 	return expr[0] != false, nil
// }

// func notTrue(expr List) (Any, error) {
// 	b, _ = !isTrue(expr[0])
// 	return !b.(bool), nil
// }

func list(exprs List) (Any, error) {
	return exprs, nil
}
