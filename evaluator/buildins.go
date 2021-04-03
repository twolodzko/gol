package evaluator

import (
	"fmt"
	"math"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/goal/environment"
	. "github.com/twolodzko/goal/types"
)

func InitBuildin() *environment.Env {
	env := environment.NewEnv()
	env.Objects = map[Symbol]Any{
		"list": &SimpleFunction{
			func(args []Any, env *environment.Env) (Any, error) {
				return List(args), nil
			},
		},
		"quote": &SpecialFunction{
			func(args []Any, env *environment.Env) (Any, error) {
				if len(args) == 1 {
					return args[0], nil
				}
				return List(args), nil
			},
		},
		"if":     &SpecialFunction{ifFn},
		"def":    &SpecialFunction{defFn},
		"let":    &SpecialFunction{letFn},
		"head":   &SpecialFunction{headFn},
		"tail":   &SpecialFunction{tailFn},
		"nil?":   &VectorizableFunction{func(x Any) Any { return Bool(x == nil) }},
		"bool?":  &VectorizableFunction{isBool},
		"int?":   &VectorizableFunction{isInt},
		"float?": &VectorizableFunction{isFloat},
		"str?":   &VectorizableFunction{isString},
		"list?":  &VectorizableFunction{isList},
		"atom?":  &VectorizableFunction{isAtom},
		"true?":  &VectorizableFunction{func(x Any) Any { return Bool(isTrue(x)) }},
		"not":    &VectorizableFunction{func(x Any) Any { return Bool(!isTrue(x)) }},
		"and": &SimpleFunction{
			func(args []Any, env *environment.Env) (Any, error) {
				return andFn(args), nil
			},
		},
		"or": &SimpleFunction{
			func(args []Any, env *environment.Env) (Any, error) {
				return orFn(args), nil
			},
		},
		"eq?": &SimpleFunction{
			func(args []Any, env *environment.Env) (Any, error) {
				if len(args) != 2 {
					return nil, &ErrNumArgs{len(args)}
				}
				return Bool(cmp.Equal(args[0], args[1])), nil
			},
		},
		"error": &SimpleFunction{
			func(args []Any, env *environment.Env) (Any, error) {
				if len(args) != 1 {
					return nil, &ErrNumArgs{len(args)}
				}
				return nil, fmt.Errorf("%s", fmt.Sprintf("%v", args[0]))
			},
		},
		"print": &SimpleFunction{
			func(args []Any, env *environment.Env) (Any, error) {
				out := ""
				for _, o := range args {
					out += fmt.Sprintf("%v", o)
				}
				fmt.Printf("%s\n", out)
				return nil, nil
			},
		},
		"+": &ArithmeticFunction{
			func(x, y Int) Int { return x + y },
			func(x, y Float) Float { return x + y },
			0,
		},
		"-": &ArithmeticFunction{
			func(x, y Int) Int { return x - y },
			func(x, y Float) Float { return x - y },
			0,
		},
		"*": &ArithmeticFunction{
			func(x, y Int) Int { return x * y },
			func(x, y Float) Float { return x * y },
			1,
		},
		"/": &ArithmeticFunction{
			func(x, y Int) Int { return x / y },
			func(x, y Float) Float { return x / y },
			1,
		},
		"%": &ArithmeticFunction{
			func(x, y Int) Int { return x % y },
			func(x, y Float) Float { return math.Mod(x, y) },
			1,
		},
		"pow": &SimpleFunction{
			func(args []Any, env *environment.Env) (Any, error) {
				return accumulate(args, math.Pow, 1)
			},
		},
		"rem": &SimpleFunction{
			func(args []Any, env *environment.Env) (Any, error) {
				return accumulate(args, math.Remainder, 1)
			},
		},
		"env": &SpecialFunction{
			func(args []Any, env *environment.Env) (Any, error) {
				fmt.Printf("%v\n", env.Objects)
				return nil, nil
			},
		},
	}

	return env
}

type Function interface {
	Call([]Any, *environment.Env) (Any, error)
}

type SpecialFunction struct {
	fn func([]Any, *environment.Env) (Any, error)
}

func (f *SpecialFunction) Call(args []Any, env *environment.Env) (Any, error) {
	return f.fn(args, env)
}

type SimpleFunction struct {
	fn func([]Any, *environment.Env) (Any, error)
}

func (f *SimpleFunction) Call(args []Any, env *environment.Env) (Any, error) {
	args, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}
	return f.fn(args, env)
}

type VectorizableFunction struct {
	fn func(Any) Any
}

func (f *VectorizableFunction) Call(args []Any, env *environment.Env) (Any, error) {
	args, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}
	return apply(args, f.fn), nil
}
