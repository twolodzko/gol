package evaluator

import (
	"errors"
	"fmt"
	"math"

	"github.com/twolodzko/gol/environment"
	"github.com/twolodzko/gol/parser"
	. "github.com/twolodzko/gol/types"
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
	"eval": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			return Eval(obj, env)
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
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, ok := args[0].(List)
			if !ok {
				return nil, &ErrWrongType{args[0]}
			}
			objs, err := EvalAll(obj, env)
			if err != nil {
				return nil, err
			}
			return last(objs), nil
		},
	},
	"fn": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			return NewLambda(args, env)
		},
	},
	"first": &SimpleFunction{
		firstFn,
	},
	"rest": &SimpleFunction{
		restFn,
	},
	"init": &SimpleFunction{
		initFn,
	},
	"last": &SimpleFunction{
		lastFn,
	},
	"append": &SimpleFunction{
		appendFn,
	},
	"concat": &SimpleFunction{
		concatFn,
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
			str, err := toString(args, env, "")
			if err != nil {
				return nil, err
			}
			return String(str), nil
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
		andFn,
	},
	"or": &SimpleFunction{
		orFn,
	},
	"=": &SimpleFunction{
		equalFn,
	},
	">": &SimpleFunction{
		gtFn,
	},
	"<": &SimpleFunction{
		ltFn,
	},
	"error": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			return nil, fmt.Errorf("%s", fmt.Sprintf("%v", args[0]))
		},
	},
	"slurp": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			name, ok := obj.(String)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			lines, err := parser.ReadFile(string(name))
			return String(lines), err
		},
	},
	"read-string": &SimpleFunction{
		readStringFn,
	},
	"println": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) == 0 {
				fmt.Println()
				return nil, nil
			}
			str, err := toString(args, env, " ")
			if err != nil {
				return nil, err
			}
			fmt.Printf("%s\n", str)
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
