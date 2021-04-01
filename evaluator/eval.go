package evaluator

import (
	"fmt"

	. "github.com/twolodzko/goal/types"
)

func EvalExpr(expr Any) (Any, error) {
	switch expr := expr.(type) {
	case Bool, Int, Float, String:
		return expr, nil
	case Symbol:
		val, err := baseEnv.Get(expr)
		if err != nil {
			return nil, err
		}
		return EvalExpr(val)
	case List:
		return evalList(expr)
	default:
		return nil, fmt.Errorf("cannot evaluate %v of type %T", expr, expr)
	}
}

func Eval(exprs []Any) ([]Any, error) {
	var out []Any
	for _, expr := range exprs {
		val, err := EvalExpr(expr)
		if err != nil {
			return nil, err
		}
		out = append(out, val)
	}
	return out, nil
}

func evalList(expr List) (Any, error) {
	if len(expr) == 0 {
		return List{}, nil
	}

	fnName, ok := expr.Head().(Symbol)
	if !ok {
		return nil, fmt.Errorf("%q is not callable", fnName)
	}
	args := expr.Tail()

	obj, err := baseEnv.Get(fnName)
	if err != nil {
		return nil, err
	}

	fn, ok := obj.(buildin)
	if !ok {
		return nil, fmt.Errorf("%q is not callable", fnName)
	}

	if fnName != "quote" {
		args, err = Eval(args)
		if err != nil {
			return nil, err
		}
	}

	return fn(args)
}
