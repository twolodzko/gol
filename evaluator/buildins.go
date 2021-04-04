package evaluator

import (
	"errors"
	"fmt"
	"math"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/goal/environment"
	. "github.com/twolodzko/goal/types"
)

var buildins = map[Symbol]Any{
	"list": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			return List(args), nil
		},
	},
	"quote": &SpecialEvalFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) == 1 {
				return args[0], nil
			}
			return List(args), nil
		},
	},
	"if": &SpecialEvalFunction{
		ifFn,
	},
	"def": &SpecialEvalFunction{
		defFn,
	},
	"pop": &SpecialEvalFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			var out []Any
			for _, name := range args {
				switch name := name.(type) {
				case Symbol:
					val, _ := env.Get(name)
					delete(env.Objects, name)
					out = append(out, val)
				default:
					return nil, &ErrWrongType{name}
				}
			}
			return toAtomOrList(out), nil
		},
	},
	"let": &SpecialEvalFunction{
		letFn,
	},
	"do": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			return last(args), nil
		},
	},
	"first": &SpecialEvalFunction{
		firstFn,
	},
	"rest": &SpecialEvalFunction{
		restFn,
	},
	"nth": &SimpleFunction{
		nthFn,
	},
	"nil?": &VectorizableFunction{
		func(x Any) (Any, error) { return Bool(x == nil), nil },
	},
	"bool?": &VectorizableFunction{
		func(obj Any) (Any, error) {
			_, ok := obj.(Bool)
			return Bool(ok), nil
		},
	},
	"int?": &VectorizableFunction{
		func(obj Any) (Any, error) {
			_, ok := obj.(Int)
			return Bool(ok), nil
		},
	},
	"float?": &VectorizableFunction{
		func(obj Any) (Any, error) {
			_, ok := obj.(Float)
			return Bool(ok), nil
		},
	},
	"str?": &VectorizableFunction{
		func(obj Any) (Any, error) {
			_, ok := obj.(String)
			return Bool(ok), nil
		},
	},
	"int": &VectorizableFunction{
		toInt,
	},
	"float": &VectorizableFunction{
		toFloat,
	},
	"str": &VectorizableFunction{
		toString,
	},
	"list?": &VectorizableFunction{
		func(obj Any) (Any, error) {
			_, ok := obj.(List)
			return Bool(ok), nil
		},
	},
	"atom?": &VectorizableFunction{
		func(obj Any) (Any, error) {
			switch obj.(type) {
			case Bool, Int, Float, String:
				return Bool(true), nil
			default:
				return Bool(false), nil
			}
		},
	},
	"true?": &VectorizableFunction{
		func(x Any) (Any, error) { return Bool(isTrue(x)), nil },
	},
	"not": &VectorizableFunction{
		func(x Any) (Any, error) { return Bool(!isTrue(x)), nil },
	},
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
	"=": &SimpleFunction{
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
	"/": &SimpleFunction{
		floatDivFn,
	},
	"//": &SimpleFunction{
		intDivFn,
	},
	"%": &ArithmeticFunction{
		func(x, y Int) Int { return x % y },
		math.Mod,
		1,
	},
	"pow": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			return applyFloatFn(args, math.Pow, 1)
		},
	},
	"rem": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			return applyFloatFn(args, math.Remainder, 1)
		},
	},
	"env": &SpecialEvalFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) > 0 {
				return nil, errors.New("env does not take any arguments")
			}
			printEnv(env, 0)
			return nil, nil
		},
	},
}

type Function interface {
	Call([]Any, *environment.Env) (Any, error)
}

type SpecialEvalFunction struct {
	fn func([]Any, *environment.Env) (Any, error)
}

func (f *SpecialEvalFunction) Call(args []Any, env *environment.Env) (Any, error) {
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
	fn func(Any) (Any, error)
}

func (f *VectorizableFunction) Call(args []Any, env *environment.Env) (Any, error) {
	args, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}
	return f.apply(args)
}

func (f *VectorizableFunction) apply(args []Any) (Any, error) {
	var out []Any
	for _, x := range args {
		res, err := f.fn(x)
		if err != nil {
			return nil, err
		}
		out = append(out, res)
	}
	return toAtomOrList(out), nil
}

func toAtomOrList(args []Any) Any {
	switch len(args) {
	case 0:
		return nil
	case 1:
		return args[0]
	default:
		return List(args)
	}
}
