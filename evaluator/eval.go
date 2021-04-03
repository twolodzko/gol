package evaluator

import (
	"errors"
	"fmt"
	"math"

	"github.com/google/go-cmp/cmp"
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

	fnName, ok := expr.Head().(Symbol)
	if !ok {
		return nil, fmt.Errorf("%v is not callable", expr.Head())
	}
	args := expr.Tail()

	switch string(fnName) {

	case "if":
		return ifFn(args, env)

	case "def":
		return defFn(args, env)

	case "let":
		return letFn(args, env)

	case "head":
		return headFn(args, env)

	case "tail":
		return tailFn(args, env)

	case "quote":
		if len(args) == 1 {
			return args[0], nil
		}
		return List(args), nil

	}

	args, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}

	switch string(fnName) {

	case "list":
		return List(args), nil

	case "true?":
		return apply(args, func(x Any) Any { return Bool(isTrue(x)) }), nil

	case "not":
		return apply(args, func(x Any) Any { return Bool(!isTrue(x)) }), nil

	case "and":
		return andFn(args), nil

	case "or":
		return orFn(args), nil

	case "nil?":
		return apply(args, func(x Any) Any { return Bool(x == nil) }), nil

	case "bool?":
		return apply(args, isBool), nil

	case "int?":
		return apply(args, isInt), nil

	case "float?":
		return apply(args, isFloat), nil

	case "str?":
		return apply(args, isString), nil

	case "list?":
		return apply(args, isList), nil

	case "atom?":
		return apply(args, isAtom), nil

	case "eq?":
		if len(args) != 2 {
			return nil, &ErrNumArgs{len(args)}
		}
		return Bool(cmp.Equal(args[0], args[1])), nil

	case "error":
		if len(args) != 1 {
			return nil, &ErrNumArgs{len(args)}
		}
		return nil, fmt.Errorf("%s", fmt.Sprintf("%v", args[0]))

	case "print":
		printFn(args)
		return nil, nil

	case "+":
		return accumulate(args, func(x, y Float) Float { return x + y }, 0)

	case "-":
		return accumulate(args, func(x, y Float) Float { return x - y }, 0)

	case "*":
		return accumulate(args, func(x, y Float) Float { return x * y }, 1)

	case "/":
		return accumulate(args, func(x, y Float) Float { return x / y }, 1)

	case "%":
		return accumulate(args, math.Mod, 1)

	case "rem":
		return accumulate(args, math.Remainder, 1)

	case "pow":
		return accumulate(args, math.Pow, 1)

	default:
		return nil, errors.New("not implemented")

	}
}
