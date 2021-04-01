package evaluator

import (
	. "github.com/twolodzko/goal/types"
)

func isTrue(obj Any) (Any, error) {
	switch obj := obj.(type) {
	case List:
		return Bool(len(obj) > 0), nil
	default:
		return Bool(obj != false), nil
	}
}

func notTrue(obj Any) (Any, error) {
	b, _ := isTrue(obj)
	return !b.(Bool), nil
}

func andFn(objs []Any) (Any, error) {
	if len(objs) == 0 {
		return Bool(false), nil
	}

	for _, x := range objs {
		if ok, _ := notTrue(x); ok.(Bool) {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func orFn(objs []Any) (Any, error) {
	for _, x := range objs {
		if ok, _ := isTrue(x); ok.(Bool) {
			return Bool(true), nil
		}
	}
	return Bool(false), nil
}
