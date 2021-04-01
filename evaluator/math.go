package evaluator

import (
	"fmt"

	. "github.com/twolodzko/goal/types"
)

func foldFnInt(fn func(Int, Int) Int) func([]Any) (Any, error) {
	return func(arr []Any) (Any, error) {
		if len(arr) < 2 {
			return nil, fmt.Errorf("wrong number of arguments (%d)", len(arr))
		}

		acc, ok := arr[0].(Int)
		if !ok {
			return 0, fmt.Errorf("invalid argument %v of type %T", acc, acc)
		}

		for _, x := range arr[1:] {
			x, ok := x.(Int)
			if !ok {
				return 0, fmt.Errorf("invalid argument %v of type %T", x, x)
			}
			acc = fn(acc, x)
		}
		return acc, nil
	}
}

func foldFnFloat(fn func(Float, Float) Float) func([]Any) (Any, error) {
	return func(arr []Any) (Any, error) {
		if len(arr) < 2 {
			return nil, fmt.Errorf("wrong number of arguments (%d)", len(arr))
		}

		acc, ok := arr[0].(Float)
		if !ok {
			return 0, fmt.Errorf("invalid argument %v of type %T", acc, acc)
		}

		for _, x := range arr[1:] {
			x, ok := x.(Float)
			if !ok {
				return 0, fmt.Errorf("invalid argument %v of type %T", x, x)
			}
			acc = fn(acc, x)
		}
		return acc, nil
	}
}
