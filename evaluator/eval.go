package evaluator

import (
	"fmt"
	"strings"

	"github.com/twolodzko/gol/environment"
	"github.com/twolodzko/gol/parser"
)

type Evaluator struct {
	env *environment.Env
}

func NewEvaluator() *Evaluator {
	baseEnv := environment.NewEnv(nil)
	baseEnv.Objects = buildins

	// so that we shadow rather than overwrite the buildins
	workEnv := environment.NewEnv(baseEnv)
	return &Evaluator{workEnv}
}

func (e *Evaluator) EvalString(code string) ([]Any, error) {
	expr, err := parser.Parse(strings.NewReader(code))
	if err != nil {
		return nil, err
	}
	return EvalAll(expr, e.env)
}

func Eval(expr Any, env *environment.Env) (Any, error) {
	var (
		newExpr Any
		newEnv  *environment.Env
	)

	for {
		switch expr := expr.(type) {
		case nil, Bool, Int, Float, String, Function, TailCallOptimized:
			return expr, nil
		case Symbol:
			return env.Get(expr)
		case List:
			if len(expr) == 0 {
				return List{}, nil
			}
			fn, err := getFunction(expr.Head(), env)
			if err != nil {
				return nil, err
			}

			switch fn := fn.(type) {
			case Function:
				args := expr.Tail()
				return fn.Call(args, env)
			case TailCallOptimized:
				args := expr.Tail()
				newExpr, newEnv, err = fn.Call(args, env)
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, fmt.Errorf("cannot evaluate %v of type %T", expr, expr)
		}

		env = newEnv
		expr = newExpr
	}
}

func EvalAll(exprs []Any, env *environment.Env) ([]Any, error) {
	var evaluated []Any
	for _, expr := range exprs {
		val, err := Eval(expr, env)
		if err != nil {
			return nil, err
		}
		evaluated = append(evaluated, val)
	}
	return evaluated, nil
}

func getFunction(obj Any, env *environment.Env) (Any, error) {
	switch obj := obj.(type) {
	case Function, TailCallOptimized:
		return obj, nil
	case Symbol:
		o, err := env.Get(obj)
		if err != nil {
			return nil, err
		}
		switch fn := o.(type) {
		case Function, TailCallOptimized:
			return fn, nil
		default:
			return nil, fmt.Errorf("%v (%T) is not callable", o, o)
		}
	case List:
		val, err := Eval(obj, env)
		if err != nil {
			return nil, err
		}
		return getFunction(val, env)
	default:
		return nil, fmt.Errorf("%v (%T) is not callable", obj, obj)
	}
}
