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
			args, err := EvalAll(args, env)
			if err != nil {
				return nil, err
			}
			return List(args), nil
		},
	},
	"quote": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			return args[0], nil
		},
	},
	"if": &SimpleFunction{
		ifFn,
	},
	"def": &SimpleFunction{
		defFn,
	},
	"del": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}

			name, ok := args[0].(Symbol)
			if !ok {
				return nil, &ErrWrongType{args[0]}
			}

			obj, _ := env.Get(name)
			delete(env.Objects, name)
			return obj, nil
		},
	},
	"let": &SimpleFunction{
		letFn,
	},
	"do": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			objs, err := EvalAll(args, env)
			if err != nil {
				return nil, err
			}
			return last(objs), nil
		},
	},
	"first": &SimpleFunction{
		firstFn,
	},
	"rest": &SimpleFunction{
		restFn,
	},
	"append": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) < 2 {
				return nil, &ErrNumArgs{len(args)}
			}
			objs, err := EvalAll(args, env)
			if err != nil {
				return nil, err
			}
			l, ok := objs[0].(List)
			if !ok {
				return nil, &ErrWrongType{args[0]}
			}
			return List(append(l, objs[1:]...)), nil
		},
	},
	"concat": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) < 2 {
				return nil, &ErrNumArgs{len(args)}
			}
			objs, err := EvalAll(args, env)
			if err != nil {
				return nil, err
			}
			var list List
			for _, obj := range objs {
				switch obj := obj.(type) {
				case List:
					list = append(list, obj...)
				default:
					return nil, &ErrWrongType{obj}
				}
			}
			return list, nil
		},
	},
	"nth": &SimpleFunction{
		nthFn,
	},
	"nil?": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			return Bool(obj == nil), nil
		},
	},
	"int?": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			_, ok := obj.(Int)
			return Bool(ok), nil
		},
	},
	"float?": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			_, ok := obj.(Float)
			return Bool(ok), nil
		},
	},
	"str?": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			_, ok := obj.(String)
			return Bool(ok), nil
		},
	},
	"fn?": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			_, ok := obj.(Function)
			return Bool(ok), nil
		},
	},
	"int": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			return toInt(obj)
		},
	},
	"float": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			return toFloat(obj)
		},
	},
	"str": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			return toString(obj)
		},
	},
	"list?": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			_, ok := obj.(List)
			return Bool(ok), nil
		},
	},
	"atom?": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			switch obj.(type) {
			case Bool, Int, Float, String:
				return Bool(true), nil
			default:
				return Bool(false), nil
			}
		},
	},
	"true?": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			return Bool(isTrue(obj)), nil
		},
	},
	"not": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			return Bool(!isTrue(obj)), nil
		},
	},
	"and": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			objs, err := EvalAll(args, env)
			if err != nil {
				return nil, err
			}
			return andFn(objs), nil
		},
	},
	"or": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			objs, err := EvalAll(args, env)
			if err != nil {
				return nil, err
			}
			return orFn(objs), nil
		},
	},
	"=": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 2 {
				return nil, &ErrNumArgs{len(args)}
			}
			objs, err := EvalAll(args, env)
			if err != nil {
				return nil, err
			}
			return Bool(cmp.Equal(objs[0], objs[1])), nil
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
			objs, err := EvalAll(args, env)
			if err != nil {
				return nil, err
			}

			out := ""
			for _, o := range objs {
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
	"env": &SimpleFunction{
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

type SimpleFunction struct {
	fn func([]Any, *environment.Env) (Any, error)
}

func (f *SimpleFunction) Call(args []Any, env *environment.Env) (Any, error) {
	return f.fn(args, env)
}
