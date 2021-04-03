package evaluator

import (
	"fmt"

	"github.com/twolodzko/goal/environment"
	. "github.com/twolodzko/goal/types"
)

func Eval(expr Any, env *environment.Env) (Any, error) {
	switch expr := expr.(type) {
	case nil, Bool, Int, Float, String:
		return expr, nil
	case Symbol:
		val, err := env.Get(expr)
		if err != nil {
			return nil, err
		}
		return val, nil
	case List:
		return evalList(expr, env)
	default:
		return nil, fmt.Errorf("cannot evaluate %v of type %T", expr, expr)
	}
}

func EvalAll(exprs []Any, env *environment.Env) ([]Any, error) {
	var out []Any
	for _, expr := range exprs {
		val, err := Eval(expr, env)
		if err != nil {
			return nil, err
		}
		out = append(out, val)
	}
	return out, nil
}

func evalList(expr List, env *environment.Env) (Any, error) {

	if len(expr) == 0 {
		return List{}, nil
	}

	args := expr.Tail()

	fnName, ok := expr.Head().(Symbol)
	if !ok {
		return nil, fmt.Errorf("%v is not callable", expr.Head())
	}

	fn, err := env.Get(fnName)
	if err != nil {
		return nil, err
	}

	switch fn := fn.(type) {
	case Function:
		return fn.Call(args, env)
	default:
		return nil, fmt.Errorf("%v is not callable", fn)
	}
}
