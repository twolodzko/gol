package evaluator

import (
	"fmt"
	"strings"

	"github.com/twolodzko/goal/environment"
	"github.com/twolodzko/goal/parser"
	. "github.com/twolodzko/goal/types"
)

type Evaluator struct {
	env *environment.Env
}

func NewEvaluator() *Evaluator {
	baseEnv := environment.NewEnv()
	baseEnv.Objects = buildins
	workEnv := environment.NewEnclosedEnv(baseEnv)
	return &Evaluator{workEnv}
}

func (e *Evaluator) Eval(cmd string) ([]Any, error) {
	expr, err := parser.Parse(strings.NewReader(cmd))
	if err != nil {
		return nil, err
	}
	return EvalAll(expr, e.env)
}

func Eval(expr Any, env *environment.Env) (Any, error) {
	switch expr := expr.(type) {
	case nil, Bool, Int, Float, String:
		return expr, nil
	case Symbol:
		return env.Get(expr)
	case List:
		if len(expr) == 0 {
			return List{}, nil
		}

		name, ok := expr.Head().(Symbol)
		if !ok {
			return nil, fmt.Errorf("%v is not callable", expr.Head())
		}

		fn, err := env.Get(name)
		if err != nil {
			return nil, err
		}

		switch fn := fn.(type) {
		case Function:
			args := expr.Tail()
			return fn.Call(args, env)
		default:
			return nil, fmt.Errorf("%v is not callable", fn)
		}
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
