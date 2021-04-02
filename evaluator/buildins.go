package evaluator

import (
	"errors"
	"fmt"

	. "github.com/twolodzko/goal/types"
)

type buildIn = func([]Any) (Any, error)

var buildIns = map[Symbol]Any{
	"list":    listFn,
	"size":    vectorize(sizeFn),
	"head":    headFn,
	"tail":    tailFn,
	"nil?":    vectorize(isNilFn),
	"error":   errorFn,
	"=":       equalFn,
	"print":   printFn,
	"println": printLnFn,
	// type conversions
	"str":   vectorize(toStringFn),
	"int":   vectorize(toIntFn),
	"float": vectorize(toFloatFn),
	// logic
	"true?": vectorize(isTrueFn),
	"not":   vectorize(notTrueFn),
	"and":   andFn,
	"or":    orFn,
	// math
	"+":   sumFn,
	"-":   difFn,
	"*":   mulFn,
	"/":   floatDivFn,
	"//":  intDivFn,
	"%":   intModFn,
	"pow": powFn,
	"rem": remFn,
}

func vectorize(fn func(Any) (Any, error)) buildIn {
	return func(objs []Any) (Any, error) {
		if len(objs) == 1 {
			return fn(objs[0])
		}

		var out List
		for _, x := range objs {
			result, err := fn(x)
			if err != nil {
				return out, err
			}
			out = append(out, result)
		}
		return out, nil
	}
}

func listFn(exprs []Any) (Any, error) {
	return List(exprs), nil
}

func sizeFn(obj Any) (Any, error) {
	switch obj := obj.(type) {
	case List:
		return len(obj), nil
	case String:
		return len(obj), nil
	default:
		return nil, nil
	}
}

func headFn(obj []Any) (Any, error) {
	if len(obj) != 1 {
		return nil, &errNumArgs{len(obj)}
	}

	l, ok := obj[0].(List)
	if !ok {
		return nil, fmt.Errorf("%v is not a list", obj[0])
	}

	return l.Head(), nil
}

func tailFn(obj []Any) (Any, error) {
	if len(obj) != 1 {
		return nil, &errNumArgs{len(obj)}
	}

	l, ok := obj[0].(List)
	if !ok {
		return nil, fmt.Errorf("%v is not a list", obj[0])
	}

	return l.Tail(), nil
}

func isNilFn(obj Any) (Any, error) {
	return Bool(obj == nil), nil
}

func errorFn(obj []Any) (Any, error) {
	if len(obj) != 1 {
		return nil, &errNumArgs{len(obj)}
	}
	msg, ok := obj[0].(String)
	if !ok {
		return nil, &errWrongType{obj[0]}
	}

	return nil, fmt.Errorf("%s", msg)
}

func equalFn(obj []Any) (Any, error) {
	if len(obj) != 2 {
		return nil, &errNumArgs{len(obj)}
	}
	if _, ok := obj[0].(List); ok {
		return nil, errors.New("= cannot be used to compare lists")
	}
	if _, ok := obj[1].(List); ok {
		return nil, errors.New("= cannot be used to compare lists")
	}
	return Bool(obj[0] == obj[1]), nil
}

func printFn(obj []Any) (Any, error) {
	out := ""
	for _, o := range obj {
		out += fmt.Sprintf("%v", o)
	}
	fmt.Print(out)
	return nil, nil
}

func printLnFn(obj []Any) (Any, error) {
	out := ""
	for _, o := range obj {
		out += fmt.Sprintf("%v", o)
	}
	fmt.Print(out + "\n")
	return nil, nil
}
