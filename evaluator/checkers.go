package evaluator

import (
	. "github.com/twolodzko/goal/types"
)

func IsBool(obj Any) (Any, error) {
	_, ok := obj.(Bool)
	return Bool(ok), nil
}

func IsInt(obj Any) (Any, error) {
	_, ok := obj.(Int)
	return Bool(ok), nil
}

func IsFloat(obj Any) (Any, error) {
	_, ok := obj.(Float)
	return Bool(ok), nil
}

func IsString(obj Any) (Any, error) {
	_, ok := obj.(String)
	return Bool(ok), nil
}

func IsAtom(obj Any) (Any, error) {
	switch obj.(type) {
	case Bool, Int, Float, String:
		return Bool(true), nil
	default:
		return Bool(false), nil
	}
}

func IsList(obj Any) (Any, error) {
	_, ok := obj.(List)
	return Bool(ok), nil
}
