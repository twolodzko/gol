package evaluator

import (
	"errors"
	"fmt"

	"github.com/twolodzko/goal/objects"
)

func Eval(expr objects.Object) (objects.Object, error) {
	switch expr := expr.(type) {
	case objects.Int, objects.Float, objects.String:
		return expr, nil
	case objects.Symbol:
		val, err := baseEnv.Get(expr.Val)
		if err != nil {
			return nil, err
		}
		return val, nil
	case objects.List:
		return evalList(expr)
	default:
		return nil, fmt.Errorf("cannot evaluate object of type %T", expr)
	}
}

func evalList(expr objects.List) (objects.Object, error) {
	if expr.Size() > 0 {
		switch name := expr.Head().(type) {
		case objects.Symbol:
			if name.Val == "quote" {
				return quote(expr.Tail())
			}

			fn, err := baseEnv.Get(name.Val)
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
	return expr, nil
}

func evalFn(fn Function, exprs []objects.Object) (objects.Object, error) {
	return fn.Call(exprs)
}

func evalAll(exprs []objects.Object) ([]objects.Object, error) {
	for i, expr := range exprs {
		val, err := Eval(expr)
		if err != nil {
			return nil, err
		}
		exprs[i] = val
	}
	return exprs, nil
}

func quote(expr []objects.Object) (objects.Object, error) {
	if len(expr) != 1 {
		return nil, errors.New("wrong number of arguments")
	}
	return expr[0], nil
}
