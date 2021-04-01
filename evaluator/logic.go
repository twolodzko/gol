package evaluator

import (
	. "github.com/twolodzko/goal/types"
)

func isTrue(obj Any) bool {
	return obj != nil && obj != false
}

func isTrueFn(obj Any) (Any, error) {
	return Bool(isTrue(obj)), nil
}

func notTrueFn(obj Any) (Any, error) {
	return Bool(!isTrue(obj)), nil
}

func andFn(objs []Any) (Any, error) {
	if len(objs) == 0 {
		return Bool(false), nil
	}

	for _, x := range objs {
		if !isTrue(x) {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func orFn(objs []Any) (Any, error) {
	for _, x := range objs {
		if isTrue(x) {
			return Bool(true), nil
		}
	}
	return Bool(false), nil
}
