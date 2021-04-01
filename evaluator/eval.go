package evaluator

import (
	"fmt"

	. "github.com/twolodzko/goal/types"
)

func Eval(expr Any) (Any, error) {
	switch expr := expr.(type) {
	case Bool, Int, Float, String:
		return expr, nil
	case Symbol:
		val, err := baseEnv.Get(expr)
		if err != nil {
			return nil, err
		}
		return Eval(val)
	case List:
		return evalList(expr)
	default:
		return nil, fmt.Errorf("cannot evaluate %v of type %T", expr, expr)
	}
}

func EvalAll(exprs []Any) (List, error) {
	var out []Any
	for _, expr := range exprs {
		val, err := Eval(expr)
		if err != nil {
			return nil, err
		}
		out = append(out, val)
	}
	return out, nil
}

func evalList(expr List) (Any, error) {
	var args []Any
	fnName := expr.Head()

	if len(expr) == 0 {
		return List{}, nil
	}

	switch fnName := fnName.(type) {
	case Symbol:
		fn, err := baseEnv.Get(fnName)
		if err != nil {
			return nil, err
		}

		switch fn := fn.(type) {
		case buildin:
			args = expr.Tail()

			if fnName != "quote" {
				args, err = EvalAll(args)
				if err != nil {
					return nil, err
				}
			}

			return fn(args)
		default:
			return nil, fmt.Errorf("%v is not a function", fn)
		}
	default:
		return nil, fmt.Errorf("cannot evaluate list: %v", expr)
	}
}
