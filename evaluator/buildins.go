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

	// core functions
	"def": &simpleFunction{
		defFn,
	},
	"fn": &simpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			return newLambda(args, env)
		},
	},
	"if": &tcoFunction{
		ifFn,
	},
	"cond": &tcoFunction{
		condFun,
	},
	"let": &tcoFunction{
		letFn,
	},
	"begin": &simpleFunction{
		beginFn,
	},
	"apply": &simpleFunction{
		applyFn,
	},
	"map": &simpleFunction{
		mapFn,
	},
	"set!": &simpleFunction{
		setFn,
	},

	// metaprogramming
	"quote": &simpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			return args[0], nil
		},
	},
	"quasiquote": &simpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			return quasiquote(args[0], env)
		},
	},
	"unquote": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return obj, nil
		},
	},
	"eval": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return eval(obj, env)
		},
	},

	// logical checks
	"=": &simpleFunction{
		equalFn,
	},
	"true?": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return Bool(isTrue(obj)), nil
		},
	},
	"not": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return Bool(!isTrue(obj)), nil
		},
	},
	"and": &simpleFunction{
		andFn,
	},
	"or": &simpleFunction{
		orFn,
	},

	// lists
	"list": &simpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			args, err := evalAll(args, env)
			return List(args), err
		},
	},
	"first": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, fmt.Errorf("%v is not a list", obj)
			}
			return l.Head(), nil
		},
	},
	"rest": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, fmt.Errorf("%v is not a list", obj)
			}
			return l.Tail(), nil
		},
	},
	"init": &singleArgFunction{
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
	"last": &singleArgFunction{
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
	"nth": &simpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 2 {
				return nil, &ErrNumArgs{len(args)}
			}
			args, err := evalAll(args, env)
			if err != nil {
				return nil, err
			}
			return nthFn(args)
		},
	},
	"append": &simpleFunction{
		appendFn,
	},
	"cons": &simpleFunction{
		prependFn,
	},
	"concat": &simpleFunction{
		concatFn,
	},
	"empty?": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			return Bool(len(l) == 0), nil
		},
	},

	// type checks
	"nil?": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return Bool(obj == nil), nil
		},
	},
	"int?": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			_, ok := obj.(Int)
			return Bool(ok), nil
		},
	},
	"float?": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			_, ok := obj.(Float)
			return Bool(ok), nil
		},
	},
	"str?": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			_, ok := obj.(String)
			return Bool(ok), nil
		},
	},
	"list?": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			_, ok := obj.(List)
			return Bool(ok), nil
		},
	},
	"atom?": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			switch obj.(type) {
			case Bool, Int, Float, String:
				return Bool(true), nil
			default:
				return Bool(false), nil
			}
		},
	},
	"fn?": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			_, ok := obj.(function)
			return Bool(ok), nil
		},
	},

	// type conversions
	"int": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return toInt(obj)
		},
	},
	"float": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return toFloat(obj)
		},
	},
	"str": &simpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			str, err := toString(args, env, "")
			return String(str), err
		},
	},

	// math
	">": &simpleFunction{
		gtFn,
	},
	"<": &simpleFunction{
		ltFn,
	},
	"+": &arithmeticFunction{
		func(x, y Int) Int { return x + y },
		func(x, y Float) Float { return x + y },
		0,
	},
	"-": &arithmeticFunction{
		func(x, y Int) Int { return x - y },
		func(x, y Float) Float { return x - y },
		0,
	},
	"*": &arithmeticFunction{
		func(x, y Int) Int { return x * y },
		func(x, y Float) Float { return x * y },
		1,
	},
	"/": &simpleFunction{
		floatDivFn,
	},
	"//": &simpleFunction{
		intDivFn,
	},
	"%": &arithmeticFunction{
		func(x, y Int) Int { return x % y },
		math.Mod,
		1,
	},
	"pow": &simpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			return applyFloatFn(args, math.Pow, 1)
		},
	},
	"rem": &simpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			return applyFloatFn(args, math.Remainder, 1)
		},
	},

	// other
	"println": &simpleFunction{
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
	"error": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			return nil, fmt.Errorf("%s", fmt.Sprintf("%v", obj))
		},
	},
	"time": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			start := time.Now()

			obj, err := eval(obj, env)
			if err != nil {
				return nil, err
			}

			end := time.Now()
			elapsed := end.Sub(start)
			fmt.Printf("%s\n", elapsed)

			return obj, nil
		},
	},
	"slurp": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			name, ok := obj.(String)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			lines, err := parser.ReadFile(string(name))
			return String(lines), err
		},
	},
	"read-string": &singleArgFunction{
		readStringFn,
	},
	"env": &simpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) > 0 {
				return nil, errors.New("env does not take any arguments")
			}
			printEnv(env, 0)
			return nil, nil
		},
	},
}
