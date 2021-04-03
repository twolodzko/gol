package evaluator

import (
	. "github.com/twolodzko/goal/types"
)

func isBool(obj Any) Any {
	_, ok := obj.(Bool)
	return Bool(ok)
}

func isInt(obj Any) Any {
	_, ok := obj.(Int)
	return Bool(ok)
}

func isFloat(obj Any) Any {
	_, ok := obj.(Float)
	return Bool(ok)
}

func isString(obj Any) Any {
	_, ok := obj.(String)
	return Bool(ok)
}

func isAtom(obj Any) Any {
	switch obj.(type) {
	case Bool, Int, Float, String:
		return Bool(true)
	default:
		return Bool(false)
	}
}

func isList(obj Any) Any {
	_, ok := obj.(List)
	return Bool(ok)
}
