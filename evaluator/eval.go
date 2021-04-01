package evaluator

import (
	"fmt"

	. "github.com/twolodzko/goal/types"
)

func EvalExpr(expr Any, env *Env) (Any, error) {
	switch expr := expr.(type) {
	case nil, Bool, Int, Float, String:
		return expr, nil
	case Symbol:
		val, err := env.Get(expr)
		if err != nil {
			return nil, err
		}
		return EvalExpr(val, env)
	case List:
		return evalList(expr, env)
	default:
		return nil, fmt.Errorf("cannot evaluate %v of type %T", expr, expr)
	}
}

func EvalAll(exprs []Any, env *Env) ([]Any, error) {
	var out []Any
	for _, expr := range exprs {
		val, err := EvalExpr(expr, env)
		if err != nil {
			return nil, err
		}
		out = append(out, val)
	}
	return out, nil
}

func evalList(expr List, env *Env) (Any, error) {
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
		cond, err := EvalExpr(args[0], env)
		if err != nil {
			return nil, err
		}
		if isTrue(cond) {
			return EvalExpr(args[1], env)
		}
		return EvalExpr(args[2], env)
	case Symbol("let"):
		if len(args) != 2 {
			return nil, &errNumArgs{len(args)}
		}

		name, ok := args[0].(Symbol)
		if !ok {
			return nil, &errWrongType{args[0]}
		}

		err := env.Set(name, args[1])
		if err != nil {
			return nil, err
		}

		return nil, nil
	default:
		obj, err := env.Get(fnName)
		if err != nil {
			return nil, err
		}

		fn, ok := obj.(buildIn)
		if !ok {
			return nil, fmt.Errorf("%q is not callable", fnName)
		}

		args, err = EvalAll(args, env)
		if err != nil {
			return nil, err
		}
		return fn(args)
	}
}
