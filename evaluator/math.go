package evaluator

import (
	. "github.com/twolodzko/goal/types"
)

func foldFnInt(fn func(Int, Int) Int) func([]Any) (Any, error) {
	return func(obj []Any) (Any, error) {
		if len(obj) < 2 {
			return nil, &errNumArgs{len(obj)}
		}

		acc, ok := obj[0].(Int)
		if !ok {
			return 0, &errWrongType{obj[0]}
		}

		for _, x := range obj[1:] {
			i, ok := x.(Int)
			if !ok {
				return 0, &errWrongType{x}
			}
			acc = fn(acc, i)
		}
		return acc, nil
	}
}

func foldFnFloat(fn func(Float, Float) Float) func([]Any) (Any, error) {
	return func(obj []Any) (Any, error) {
		if len(obj) < 2 {
			return nil, &errNumArgs{len(obj)}
		}

		acc, ok := obj[0].(Float)
		if !ok {
			return 0, &errWrongType{obj[0]}
		}

		for _, x := range obj[1:] {
			f, ok := x.(Float)
			if !ok {
				return 0, &errWrongType{x}
			}
			acc = fn(acc, f)
		}
		return acc, nil
	}
}
