package evaluator

import (
	"fmt"

	. "github.com/twolodzko/goal/types"
)

type buildin = func([]Any) (Any, error)

var buildins = map[string]Any{
	"list":  listFn,
	"quote": vectorize(quoteFn),
	"size":  vectorize(sizeFn),
	"head":  headFn,
	"tail":  tailFn,
	// type conversions
	"str":   vectorize(toString),
	"int":   vectorize(toInt),
	"float": vectorize(toFloat),
	// logic
	"true?": vectorize(isTrue),
	"not":   vectorize(notTrue),
	"and":   andFn,
	"or":    orFn,
	// math
	"int+":   foldFnInt(func(x, y Int) Int { return x + y }),
	"int-":   foldFnInt(func(x, y Int) Int { return x - y }),
	"int*":   foldFnInt(func(x, y Int) Int { return x * y }),
	"int/":   foldFnInt(func(x, y Int) Int { return x / y }),
	"float+": foldFnFloat(func(x, y Float) Float { return x + y }),
	"float-": foldFnFloat(func(x, y Float) Float { return x - y }),
	"float*": foldFnFloat(func(x, y Float) Float { return x * y }),
	"float/": foldFnFloat(func(x, y Float) Float { return x / y }),
}

func vectorize(fn func(Any) (Any, error)) buildin {
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

func quoteFn(obj Any) (Any, error) {
	return obj, nil
}

func sizeFn(obj Any) (Any, error) {
	switch obj := obj.(type) {
	case List:
		return len(obj), nil
	case String:
		return len(obj), nil
	default:
		return 0, nil
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
