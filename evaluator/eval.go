package evaluator

import (
	"fmt"

	. "github.com/twolodzko/goal/types"
)

func EvalExpr(expr Any) (Any, error) {
	switch expr := expr.(type) {
	case nil, Bool, Int, Float, String:
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

func EvalAll(exprs []Any) ([]Any, error) {
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
		return nil, fmt.Errorf("%v is not callable", expr.Head())
	}
	args := expr.Tail()

	switch fnName {
	case Symbol("quote"):
		if len(args) == 1 {
			return args[0], nil
		}
		return args, nil
	case Symbol("if"):
		if len(args) != 3 {
			return nil, &errNumArgs{len(args)}
		}
		cond, err := EvalExpr(args[0])
		if err != nil {
			return nil, err
		}
		if isTrue(cond) {
			return EvalExpr(args[1])
		}
		return EvalExpr(args[2])
	default:
		obj, err := baseEnv.Get(fnName)
		if err != nil {
			return nil, err
		}

		fn, ok := obj.(buildIn)
		if !ok {
			return nil, fmt.Errorf("%q is not callable", fnName)
		}

		args, err = EvalAll(args)
		if err != nil {
			return nil, err
		}
		return fn(args)
	}
}
