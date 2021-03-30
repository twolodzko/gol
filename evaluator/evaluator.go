package evaluator

import (
	"fmt"

	"github.com/twolodzko/goal/objects"
)

func Eval(expr objects.Object) (objects.Object, error) {
	switch expr := expr.(type) {
	case objects.Int, objects.Float, objects.String:
		return expr, nil
	case objects.List:
		return evalFn(expr)
	default:
		return nil, fmt.Errorf("cannot evaluate object of type %T", expr)
	}
}

func evalFn(expr objects.List) (objects.Object, error) {
	if expr.Size() > 0 {
		switch name := expr.Head().(type) {
		case objects.Symbol:
			args, err := evalAll(expr.Tail())
			if err != nil {
				return nil, err
			}

			fn, ok := buildins[name.Val]
			if !ok {
				return nil, fmt.Errorf("undefined function: %s", name.Val)
			}

			return fn(args)
		default:
			return nil, fmt.Errorf("cannot evaluate list: %v", expr)
		}
	}
	return expr, nil
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
