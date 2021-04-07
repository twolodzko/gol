package evaluator

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/twolodzko/gol/environment"
	"github.com/twolodzko/gol/parser"
)

var buildins = map[Symbol]Any{
	"if": &TcoFunction{
		ifFn,
	},
	"cond": &TcoFunction{
		condFun,
	},
	"let": &TcoFunction{
		letFn,
	},
	"fn": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			return NewLambda(args, env)
		},
	},
	"def": &SimpleFunction{
		defFn,
	},
	"set!": &SimpleFunction{
		setFn,
	},
	"begin": &SimpleFunction{
		beginFn,
	},
	"list": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			args, err := EvalAll(args, env)
			return List(args), err
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
	"eval": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return Eval(obj, env)
		},
	},
	"first": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, fmt.Errorf("%v is not a list", obj)
			}
			return l.Head(), nil
		},
	},
	"rest": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, fmt.Errorf("%v is not a list", obj)
			}
			return l.Tail(), nil
		},
	},
	"init": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, fmt.Errorf("%v is not a list", obj)
			}
			if len(l) == 0 {
				return List{}, nil
			}
			return l[:len(l)-1], nil
		},
	},
	"last": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, fmt.Errorf("%v is not a list", obj)
			}
			if len(l) == 0 {
				return nil, nil
			}
			return l[len(l)-1], nil
		},
	},
	"nth": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 2 {
				return nil, &ErrNumArgs{len(args)}
			}
			args, err := EvalAll(args, env)
			if err != nil {
				return nil, err
			}
			return nthFn(args)
		},
	},
	"append": &SimpleFunction{
		appendFn,
	},
	"cons": &SimpleFunction{
		prependFn,
	},
	"concat": &SimpleFunction{
		concatFn,
	},
	"empty?": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			return Bool(len(l) == 0), nil
		},
	},
	"nil?": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return Bool(obj == nil), nil
		},
	},
	"int?": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			_, ok := obj.(Int)
			return Bool(ok), nil
		},
	},
	"float?": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			_, ok := obj.(Float)
			return Bool(ok), nil
		},
	},
	"str?": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			_, ok := obj.(String)
			return Bool(ok), nil
		},
	},
	"fn?": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			_, ok := obj.(Function)
			return Bool(ok), nil
		},
	},
	"int": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return toInt(obj)
		},
	},
	"float": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return toFloat(obj)
		},
	},
	"str": &SimpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			str, err := toString(args, env, "")
			return String(str), err
		},
	},
	"list?": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			_, ok := obj.(List)
			return Bool(ok), nil
		},
	},
	"atom?": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			switch obj.(type) {
			case Bool, Int, Float, String:
				return Bool(true), nil
			default:
				return Bool(false), nil
			}
		},
	},
	"true?": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return Bool(isTrue(obj)), nil
		},
	},
	"not": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
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
	"error": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return nil, fmt.Errorf("%s", fmt.Sprintf("%v", obj))
		},
	},
	"slurp": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			name, ok := obj.(String)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			lines, err := parser.ReadFile(string(name))
			return String(lines), err
		},
	},
	"read-string": &SingleArgFunction{
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
	"time": &SingleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			start := time.Now()

			obj, err := Eval(obj, env)
			if err != nil {
				return nil, err
			}

			end := time.Now()
			elapsed := end.Sub(start)
			fmt.Printf("%s\n", elapsed)

			return obj, nil
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
