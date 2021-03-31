package evaluator

import (
	"errors"
	"fmt"

	. "github.com/twolodzko/goal/types"
)

func Eval(expr Any) (Any, error) {
	switch expr := expr.(type) {
	case Int, Float, String:
		return expr, nil
	case Symbol:
		val, err := baseEnv.Get(string(expr))
		if err != nil {
			return nil, err
		}
		return val, nil
	case List:
		return evalList(expr)
	default:
		return nil, fmt.Errorf("cannot evaluate object of type %T", expr)
	}
}

func evalList(expr List) (Any, error) {
	if len(expr) > 0 {
		switch name := expr.Head().(type) {
		case Symbol:
			if name == "quote" {
				return quote(expr.Tail())
			}

			fn, err := baseEnv.Get(string(name))
			if err != nil {
				return nil, err
			}

			switch fn := fn.(type) {
			case fixedArgsFunction:
				return fn.Call(expr.Tail())
			case anyArgsFunction:
				return fn.Call(expr.Tail())
			}
		default:
			return nil, fmt.Errorf("cannot evaluate list: %v", expr)
		}
	}
	return List{}, nil
}

func evalFn(fn Function, exprs List) (Any, error) {
	return fn.Call(exprs)
}

func evalAll(exprs List) (List, error) {
	for i, expr := range exprs {
		val, err := Eval(expr)
		if err != nil {
			return nil, err
		}
		exprs[i] = val
	}
	return exprs, nil
}

func quote(expr List) (Any, error) {
	if len(expr) != 1 {
		return nil, errors.New("wrong number of arguments")
	}
	return expr[0], nil
}
