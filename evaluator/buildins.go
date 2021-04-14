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
		// (def <name> <expr>)
		// (def (<name> <arg>...) <expr>...)
		defFn,
	},
	"fn": &simpleFunction{
		// (fn (<arg>...) <expr>...)
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) < 2 {
				return nil, &ErrNumArgs{len(args)}
			}
			return newLambda(args[0], args[1:], env)
		},
	},
	"if": &tcoFunction{
		// (if <condition> <expr> <expr>)
		ifFn,
	},
	"cond": &tcoFunction{
		// (cond (<condition> <expr>...)...)
		condFun,
	},
	"let": &tcoFunction{
		// (let (<name> <expr>...) <expr>...)
		letFn,
	},
	"begin": &multiArgFunction{
		// (begin <expr>...)
		func(objs []Any) (Any, error) {
			return last(objs), nil
		},
	},
	"apply": &simpleFunction{
		// (apply <expr> <list>)
		applyFn,
	},
	"map": &simpleFunction{
		// (map <expr> <list>)
		mapFn,
	},
	"set!": &simpleFunction{
		// (set! <name> <expr>)
		setFn,
	},

	// metaprogramming
	"quote": &simpleFunction{
		// (quote <expr>)
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			return args[0], nil
		},
	},
	"quasiquote": &simpleFunction{
		// (quasiquote <expr>)
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			return quasiquote(args[0], env)
		},
	},
	"unquote": &singleArgFunction{
		// (unquote <expr>)
		func(obj Any) (Any, error) {
			return obj, nil
		},
	},
	"eval": &simpleFunction{
		// (eval <expr>)
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}
			obj, err := eval(args[0], env)
			if err != nil {
				return nil, err
			}
			return eval(obj, env)
		},
	},
	"parse-string": &singleArgFunction{
		// (parse-string <expr>)
		parseStringFn,
	},

	// logical checks
	"=": &simpleFunction{
		// (= <expr> <expr>)
		equalFn,
	},
	"true?": &singleArgFunction{
		// (true? <expr>)
		func(obj Any) (Any, error) {
			return Bool(isTrue(obj)), nil
		},
	},
	"not": &singleArgFunction{
		// (not <expr>)
		func(obj Any) (Any, error) {
			return Bool(!isTrue(obj)), nil
		},
	},
	"and": &simpleFunction{
		// (and <expr>...)
		andFn,
	},
	"or": &simpleFunction{
		// (or <expr>...)
		orFn,
	},

	// lists
	"list": &multiArgFunction{
		// (list <expr>...)
		func(objs []Any) (Any, error) {
			return List(objs), nil
		},
	},
	"first": &singleArgFunction{
		// (first <list>)
		func(obj Any) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, fmt.Errorf("%v is not a list", obj)
			}
			return l.Head(), nil
		},
	},
	"rest": &singleArgFunction{
		// (rest <list>)
		func(obj Any) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, fmt.Errorf("%v is not a list", obj)
			}
			return l.Tail(), nil
		},
	},
	"init": &singleArgFunction{
		// (init <list>)
		func(obj Any) (Any, error) {
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
		// (last <list>)
		func(obj Any) (Any, error) {
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
		// (nth <list> <int>)
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
	"conj": &simpleFunction{
		// (conj <list> <expr>...)
		appendFn,
	},
	"cons": &simpleFunction{
		// (cons <expr> <list>)
		prependFn,
	},
	"concat": &multiArgFunction{
		// (concat <list>...)
		concatFn,
	},
	"reverse": &singleArgFunction{
		// (reverse <list>)
		func(obj Any) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			return List(reverse(l)), nil
		},
	},
	"empty?": &singleArgFunction{
		// (empty? <list>)
		func(obj Any) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			return Bool(len(l) == 0), nil
		},
	},
	"count": &singleArgFunction{
		// (count <list>)
		func(obj Any) (Any, error) {
			l, ok := obj.(List)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			return Int(len(l)), nil
		},
	},

	// type checks
	"nil?": &singleArgFunction{
		// (nil? <expr>)
		func(obj Any) (Any, error) {
			return Bool(obj == nil), nil
		},
	},
	"int?": &singleArgFunction{
		// (int? <expr>)
		func(obj Any) (Any, error) {
			_, ok := obj.(Int)
			return Bool(ok), nil
		},
	},
	"float?": &singleArgFunction{
		// (float? <expr>)
		func(obj Any) (Any, error) {
			_, ok := obj.(Float)
			return Bool(ok), nil
		},
	},
	"str?": &singleArgFunction{
		// (str? <expr>)
		func(obj Any) (Any, error) {
			_, ok := obj.(String)
			return Bool(ok), nil
		},
	},
	"list?": &singleArgFunction{
		// (list? <expr>)
		func(obj Any) (Any, error) {
			_, ok := obj.(List)
			return Bool(ok), nil
		},
	},
	"atom?": &singleArgFunction{
		// (atom? <expr>)
		func(obj Any) (Any, error) {
			switch obj.(type) {
			case Bool, Int, Float, String:
				return Bool(true), nil
			default:
				return Bool(false), nil
			}
		},
	},
	"fn?": &singleArgFunction{
		// (fn? <expr>)
		func(obj Any) (Any, error) {
			_, ok := obj.(function)
			return Bool(ok), nil
		},
	},

	// type conversions
	"int": &singleArgFunction{
		// (int <expr>)
		func(obj Any) (Any, error) {
			return toInt(obj)
		},
	},
	"float": &singleArgFunction{
		// (float <expr>)
		func(obj Any) (Any, error) {
			return toFloat(obj)
		},
	},
	"str": &multiArgFunction{
		// (str <expr>...)
		func(objs []Any) (Any, error) {
			str, err := toString(objs, "")
			return String(str), err
		},
	},

	// strings
	"print": &multiArgFunction{
		// (print <expr>...)
		func(objs []Any) (Any, error) {
			if len(objs) == 0 {
				return nil, nil
			}
			str, err := toString(objs, " ")
			if err != nil {
				return nil, err
			}
			fmt.Print(str)
			return nil, nil
		},
	},
	"println": &multiArgFunction{
		// (println <expr>...)
		func(objs []Any) (Any, error) {
			if len(objs) == 0 {
				fmt.Println()
				return nil, nil
			}
			str, err := toString(objs, " ")
			if err != nil {
				return nil, err
			}
			fmt.Println(str)
			return nil, nil
		},
	},
	"chars": &singleArgFunction{
		// (chars <expr>)
		func(obj Any) (Any, error) {
			str, ok := obj.(String)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			var chars []Any
			for _, r := range str {
				ch := String(string(r)).Quote()
				chars = append(chars, ch)
			}
			return List(chars), nil
		},
	},
	"pretty-str": &singleArgFunction{
		// (pretty-str <expr>)
		func(obj Any) (Any, error) {
			str, ok := obj.(String)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			return str.Unquote()
		},
	},
	"escaped-str": &singleArgFunction{
		// (escaped-str <expr>)
		func(obj Any) (Any, error) {
			str, ok := obj.(String)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			return str.Quote(), nil
		},
	},

	// I/O
	"read-file": &singleArgFunction{
		// (read-file <filename>)
		func(obj Any) (Any, error) {
			name, ok := obj.(String)
			if !ok {
				return nil, &ErrWrongType{obj}
			}
			lines, err := parser.ReadFile(string(name))
			return String(lines), err
		},
	},
	"write-to-file": &multiArgFunction{
		// (write-to-file <filename> <expr>)
		writeToFileFn,
	},

	// utils
	"error": &multiArgFunction{
		// (error <expr>...)
		func(objs []Any) (Any, error) {
			str, err := toString(objs, "")
			if err != nil {
				return nil, fmt.Errorf("failed parsing error message: %s", err)
			}
			return nil, fmt.Errorf(str)
		},
	},
	"time": &simpleFunction{
		// (time <expr>...)
		func(args []Any, env *environment.Env) (Any, error) {
			if len(args) != 1 {
				return nil, &ErrNumArgs{len(args)}
			}

			start := time.Now()

			objs, err := evalAll(args, env)
			if err != nil {
				return nil, err
			}

			end := time.Now()
			elapsed := end.Sub(start)
			fmt.Printf("%s\n", elapsed)

			return last(objs), nil
		},
	},
	"env": &simpleFunction{
		// (env)
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
		// (> <expr>...)
		gtFn,
	},
	"<": &simpleFunction{
		// (< <expr>...)
		ltFn,
	},
	"+": &multiArgFloatFunction{
		// (+ <expr>...)
		func(x, y Float) Float { return x + y },
		0,
	},
	"-": &multiArgFloatFunction{
		// (- <expr>...)
		func(x, y Float) Float { return x - y },
		0,
	},
	"*": &multiArgFloatFunction{
		// (* <expr>...)
		func(x, y Float) Float { return x * y },
		1,
	},
	"/": &multiArgFloatFunction{
		// (/ <expr>...)
		func(x, y Float) Float { return x / y },
		1,
	},
	"%": &multiArgFloatFunction{
		// (% <expr>...)
		math.Mod,
		1,
	},
	"pow": &multiArgFloatFunction{
		// (pow <expr>...)
		math.Pow,
		1,
	},
	"rem": &multiArgFloatFunction{
		// (rem <expr>...)
		math.Remainder,
		1,
	},
	"inf?": &singleArgFunction{
		// (inf? <expr>)
		func(obj Any) (Any, error) {
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
		// (nan? <expr>)
		func(obj Any) (Any, error) {
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
		// (sqrt <expr>)
		math.Sqrt,
	},
	"cbrt": &singleArgFloatFunction{
		// (cbrt <expr>)
		math.Cbrt,
	},
	"log": &singleArgFloatFunction{
		// (log <expr>)
		math.Log,
	},
	"log2": &singleArgFloatFunction{
		// (log2 <expr>)
		math.Log2,
	},
	"log10": &singleArgFloatFunction{
		// (log10 <expr>)
		math.Log10,
	},
	"exp": &singleArgFloatFunction{
		// (exp <expr>)
		math.Exp,
	},
	"expm1": &singleArgFloatFunction{
		// (expm1 <expr>)
		math.Expm1,
	},
	"floor": &singleArgFloatFunction{
		// (floor <expr>)
		math.Floor,
	},
	"ceil": &singleArgFloatFunction{
		// (ceil <expr>)
		math.Ceil,
	},
	"sin": &singleArgFloatFunction{
		// (sin <expr>)
		math.Sin,
	},
	"cos": &singleArgFloatFunction{
		// (cos <expr>)
		math.Cos,
	},
	"tan": &singleArgFloatFunction{
		// (tan <expr>)
		math.Tan,
	},
	"asin": &singleArgFloatFunction{
		// (asin <expr>)
		math.Asin,
	},
	"acos": &singleArgFloatFunction{
		// (acos <expr>)
		math.Acos,
	},
	"atan": &singleArgFloatFunction{
		// (atan <expr>)
		math.Atan,
	},
	"sinh": &singleArgFloatFunction{
		// (sinh <expr>)
		math.Sinh,
	},
	"cosh": &singleArgFloatFunction{
		// (cosh <expr>)
		math.Cosh,
	},
	"tanh": &singleArgFloatFunction{
		// (tanh <expr>)
		math.Tanh,
	},
	"erf": &singleArgFloatFunction{
		// (erf <expr>)
		math.Erf,
	},
	"erfc": &singleArgFloatFunction{
		// (erfc <expr>)
		math.Erfc,
	},
	"erfcinv": &singleArgFloatFunction{
		// (erfcinv <expr>)
		math.Erfcinv,
	},
	"erfinv": &singleArgFloatFunction{
		// (erfinv <expr>)
		math.Erfinv,
	},
	"gamma": &singleArgFloatFunction{
		// (gamma <expr>)
		math.Gamma,
	},
	"int+": &multiArgIntFunction{
		// (int+ <expr>...)
		func(x, y Int) Int { return x + y },
		0,
	},
	"int-": &multiArgIntFunction{
		// (int- <expr>...)
		func(x, y Int) Int { return x - y },
		0,
	},
	"int*": &multiArgIntFunction{
		// (int* <expr>...)
		func(x, y Int) Int { return x * y },
		1,
	},
	"int/": &multiArgIntFunction{
		// (int/ <expr>...)
		func(x, y Int) Int { return x / y },
		1,
	},
	"int%": &multiArgIntFunction{
		// (int% <expr>...)
		func(x, y Int) Int { return x % y },
		1,
	},
}
