package evaluator

import (
	. "github.com/twolodzko/goal/types"
)

type buildin = func([]Any) (Any, error)

var buildins = map[string]Any{
	"true":  true,
	"false": false,
	"str":   vectorize(toString),
	"int":   vectorize(toInt),
	"float": vectorize(toFloat),
	"list":  listFn,
	"quote": vectorize(quoteFn),
	"size":  vectorize(sizeFn),
	"true?": vectorize(isTrue),
	"not":   vectorize(notTrue),
	"and":   andFn,
	"or":    orFn,
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
