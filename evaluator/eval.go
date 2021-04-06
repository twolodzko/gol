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
	switch expr := expr.(type) {
	case nil, Bool, Int, Float, String, Function:
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
		args := expr.Tail()
		return fn.Call(args, env)
	default:
		return nil, fmt.Errorf("cannot evaluate %v of type %T", expr, expr)
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

func getFunction(obj Any, env *environment.Env) (Function, error) {
	switch obj := obj.(type) {
	case Function:
		return obj, nil
	case Symbol:
		o, err := env.Get(obj)
		if err != nil {
			return nil, err
		}
		fn, ok := o.(Function)
		if !ok {
			return nil, fmt.Errorf("%v (%T) is not callable", o, o)
		}
		return fn, nil
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
