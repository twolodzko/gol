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
	"parse-string": &singleArgFunction{
		parseStringFn,
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

	// strings & i/o
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
			return nil, fmt.Errorf("%v", obj)
		},
	},
	"read-file": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			name, ok := obj.(String)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			lines, err := parser.ReadFile(string(name))
			return String(lines), err
		},
	},
	"chars": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			str, ok := obj.(String)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			var runes []Any
			for _, r := range str {
				runes = append(runes, String(string(r)))
			}
			return List(runes), nil
		},
	},
	"write-to-file": &simpleFunction{
		writeToFileFn,
	},

	// utils
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
	"env": &simpleFunction{
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) > 0 {
				return nil, errors.New("env does not take any arguments")
			}
			printEnv(env, 0)
			return nil, nil
		},
	},

	// math
	">": &simpleFunction{
		gtFn,
	},
	"<": &simpleFunction{
		ltFn,
	},
	"+": &multiArgFloatFunction{
		func(x, y Float) Float { return x + y },
		0,
	},
	"-": &multiArgFloatFunction{
		func(x, y Float) Float { return x - y },
		0,
	},
	"*": &multiArgFloatFunction{
		func(x, y Float) Float { return x * y },
		1,
	},
	"/": &multiArgFloatFunction{
		func(x, y Float) Float { return x / y },
		1,
	},
	"%": &multiArgFloatFunction{
		math.Mod,
		1,
	},
	"pow": &multiArgFloatFunction{
		math.Pow,
		1,
	},
	"rem": &multiArgFloatFunction{
		math.Remainder,
		1,
	},
	"inf?": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			switch obj := obj.(type) {
			case Float:
				return math.IsInf(obj, 0), nil
			case Int:
				return false, nil
			default:
				return nil, &ErrNaN{obj}
			}
		},
	},
	"nan?": &singleArgFunction{
		func(obj Any, env *environment.Env) (Any, error) {
			switch obj := obj.(type) {
			case Float:
				return math.IsNaN(obj), nil
			case Int:
				return false, nil
			default:
				return nil, &ErrNaN{obj}
			}
		},
	},
	"sqrt": &singleArgFloatFunction{
		math.Sqrt,
	},
	"cbrt": &singleArgFloatFunction{
		math.Cbrt,
	},
	"log": &singleArgFloatFunction{
		math.Log,
	},
	"log2": &singleArgFloatFunction{
		math.Log2,
	},
	"log10": &singleArgFloatFunction{
		math.Log10,
	},
	"exp": &singleArgFloatFunction{
		math.Exp,
	},
	"expm1": &singleArgFloatFunction{
		math.Expm1,
	},
	"floor": &singleArgFloatFunction{
		math.Floor,
	},
	"ceil": &singleArgFloatFunction{
		math.Ceil,
	},
	"sin": &singleArgFloatFunction{
		math.Sin,
	},
	"cos": &singleArgFloatFunction{
		math.Cos,
	},
	"tan": &singleArgFloatFunction{
		math.Tan,
	},
	"asin": &singleArgFloatFunction{
		math.Asin,
	},
	"acos": &singleArgFloatFunction{
		math.Acos,
	},
	"atan": &singleArgFloatFunction{
		math.Atan,
	},
	"sinh": &singleArgFloatFunction{
		math.Sinh,
	},
	"cosh": &singleArgFloatFunction{
		math.Cosh,
	},
	"tanh": &singleArgFloatFunction{
		math.Tanh,
	},
	"erf": &singleArgFloatFunction{
		math.Erf,
	},
	"erfc": &singleArgFloatFunction{
		math.Erfc,
	},
	"erfcinv": &singleArgFloatFunction{
		math.Erfcinv,
	},
	"erfinv": &singleArgFloatFunction{
		math.Erfinv,
	},
	"gamma": &singleArgFloatFunction{
		math.Gamma,
	},
	"int+": &multiArgIntFunction{
		func(x, y Int) Int { return x + y },
		0,
	},
	"int-": &multiArgIntFunction{
		func(x, y Int) Int { return x - y },
		0,
	},
	"int*": &multiArgIntFunction{
		func(x, y Int) Int { return x * y },
		1,
	},
	"int/": &multiArgIntFunction{
		func(x, y Int) Int { return x / y },
		1,
	},
	"int%": &multiArgIntFunction{
		func(x, y Int) Int { return x % y },
		1,
	},
}
