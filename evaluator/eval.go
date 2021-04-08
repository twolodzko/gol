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
	return evalAll(expr, e.env)
}

func eval(expr Any, env *environment.Env) (Any, error) {
	var (
		newExpr Any
		newEnv  *environment.Env
	)

	for {
		switch expr := expr.(type) {
		case nil, Bool, Int, Float, String, function:
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
			case tailCallOptimized:
				args := expr.Tail()
				newExpr, newEnv, err = fn.PartialEval(args, env)
				if err != nil {
					return nil, err
				}
			case function:
				args := expr.Tail()
				return fn.Eval(args, env)
			}

		default:
			return nil, fmt.Errorf("cannot evaluate %v (%T)", expr, expr)
		}

		env = newEnv
		expr = newExpr
	}
}

func evalAll(exprs []Any, env *environment.Env) ([]Any, error) {
	var evaluated []Any
	for _, expr := range exprs {
		val, err := eval(expr, env)
		if err != nil {
			return nil, err
		}
		evaluated = append(evaluated, val)
	}
	return evaluated, nil
}

func getFunction(obj Any, env *environment.Env) (Any, error) {
	switch obj := obj.(type) {
	case function:
		return obj, nil
	case Symbol:
		o, err := env.Get(obj)
		if err != nil {
			return nil, err
		}
		switch fn := o.(type) {
		case function:
			return fn, nil
		default:
			return nil, &ErrNotCallable{o}
		}
	case List:
		val, err := eval(obj, env)
		if err != nil {
			return nil, err
		}
		return getFunction(val, env)
	default:
		return nil, &ErrNotCallable{obj}
	}
}
