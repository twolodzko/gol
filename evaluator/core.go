package evaluator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/gol/environment"
	"github.com/twolodzko/gol/parser"
)

func defFn(args []Any, env *environment.Env) (Any, error) {

	if len(args) != 2 {
		return nil, &ErrNumArgs{len(args)}
	}

	name, ok := args[0].(Symbol)
	if !ok {
		return nil, &ErrWrongType{args[0]}
	}

	val, err := eval(args[1], env)
	if err != nil {
		return nil, err
	}

	err = env.Set(name, val)

	return val, err
}

func setFn(args []Any, env *environment.Env) (Any, error) {

	if len(args) != 2 {
		return nil, &ErrNumArgs{len(args)}
	}

	name, ok := args[0].(Symbol)
	if !ok {
		return nil, &ErrWrongType{args[0]}
	}

	val, err := eval(args[1], env)
	if err != nil {
		return nil, err
	}

	foundEnv, err := env.Find(name)
	if err != nil {
		return nil, err
	}

	err = foundEnv.Set(name, val)

	return val, err
}

func beginFn(args []Any, env *environment.Env) (Any, error) {
	objs, err := evalAll(args, env)
	return last(objs), err
}

func quasiquote(arg Any, env *environment.Env) (Any, error) {
	switch obj := arg.(type) {
	case List:
		if len(obj) == 0 {
			return obj, nil
		}

		sym, ok := obj[0].(Symbol)

		// check for unquote's recursively
		if !ok || sym != "unquote" {
			var list List
			for _, o := range obj {
				val, err := quasiquote(o, env)
				if err != nil {
					return list, err
				}
				list = append(list, val)
			}
			return list, nil
		}

		// unquote
		if len(obj) != 2 {
			return nil, errors.New("nothing to unquote")
		}
		return eval(obj[1], env)
	default:
		return obj, nil
	}
}

func nthFn(args []Any) (Any, error) {
	switch l := args[0].(type) {
	case List:
		var n int
		switch arg := args[1].(type) {
		case Int:
			n = arg
		case Float:
			n = int(arg)
		default:
			return nil, &ErrWrongType{args[1]}
		}

		if n < len(l) {
			return l[n], nil
		}
		return nil, fmt.Errorf("arrempting to access %d element of %d", n, len(l))
	default:
		return nil, &ErrWrongType{args[0]}
	}
}

func appendFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) < 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	objs, err := evalAll(args, env)
	if err != nil {
		return nil, err
	}
	l, ok := objs[0].(List)
	if !ok {
		return nil, &ErrWrongType{args[0]}
	}
	return List(append(l, objs[1:]...)), nil
}

func prependFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) < 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	objs, err := evalAll(args, env)
	if err != nil {
		return nil, err
	}
	l, ok := objs[1].(List)
	if !ok {
		return nil, &ErrWrongType{args[0]}
	}
	return List(append([]Any{objs[0]}, l...)), nil
}

func concatFn(args []Any, env *environment.Env) (Any, error) {
	objs, err := evalAll(args, env)
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
}

func andFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) == 0 {
		return Bool(false), nil
	}
	for _, arg := range args {
		arg, err := eval(arg, env)
		if !isTrue(arg) {
			return Bool(false), err
		}
	}
	return Bool(true), nil
}

func orFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) == 0 {
		return Bool(false), nil
	}
	for _, arg := range args {
		obj, err := eval(arg, env)
		if isTrue(obj) {
			return Bool(true), err
		}
	}
	return Bool(false), nil
}

func equalFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) < 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	first, err := eval(args[0], env)
	if err != nil {
		return nil, err
	}

	for _, second := range args[1:] {
		second, err := eval(second, env)
		if err != nil {
			return nil, err
		}

		switch first := first.(type) {
		case Bool:
			second, ok := second.(Bool)
			if !ok || first != second {
				return Bool(false), nil
			}
		case Int:
			switch second := second.(type) {
			case Int:
				if first != second {
					return Bool(false), nil
				}
			case Float:
				if Float(first) != second {
					return Bool(false), nil
				}
			default:
				return Bool(false), nil
			}
		case Float:
			switch second := second.(type) {
			case Float:
				if first != second {
					return Bool(false), nil
				}
			case Int:
				if first != Float(second) {
					return Bool(false), nil
				}
			default:
				return Bool(false), nil
			}
		case String:
			second, ok := second.(String)
			if !ok || first != second {
				return Bool(false), nil
			}
		case function:
			second, ok := second.(function)
			if !ok || first != second {
				return Bool(false), nil
			}
		case tailCallOptimized:
			second, ok := second.(tailCallOptimized)
			if !ok || first != second {
				return Bool(false), nil
			}
		default:
			if !cmp.Equal(first, second) {
				return Bool(false), nil
			}
		}
	}

	return Bool(true), nil
}

func gtFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) < 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	first, err := eval(args[0], env)
	if err != nil {
		return nil, err
	}

	for _, second := range args[1:] {
		second, err := eval(second, env)
		if err != nil {
			return nil, err
		}

		switch first := first.(type) {
		case Int:
			switch second := second.(type) {
			case Int:
				if first <= second {
					return Bool(false), nil
				}
			case Float:
				if Float(first) <= second {
					return Bool(false), nil
				}
			default:
				return nil, &ErrWrongType{second}
			}
		case Float:
			switch second := second.(type) {
			case Float:
				if first <= second {
					return Bool(false), nil
				}
			case Int:
				if first <= Float(second) {
					return Bool(false), nil
				}
			default:
				return nil, &ErrWrongType{second}
			}
		default:
			return nil, &ErrWrongType{first}
		}

		first = second
	}

	return Bool(true), nil
}

func ltFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) < 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	first, err := eval(args[0], env)
	if err != nil {
		return nil, err
	}

	for _, second := range args[1:] {
		switch first := first.(type) {
		case Int:
			switch second := second.(type) {
			case Int:
				if first >= second {
					return Bool(false), nil
				}
			case Float:
				if Float(first) >= second {
					return Bool(false), nil
				}
			default:
				return nil, &ErrWrongType{second}
			}
		case Float:
			switch second := second.(type) {
			case Float:
				if first >= second {
					return Bool(false), nil
				}
			case Int:
				if first >= Float(second) {
					return Bool(false), nil
				}
			default:
				return nil, &ErrWrongType{second}
			}
		default:
			return nil, &ErrWrongType{first}
		}

		first = second
	}

	return Bool(true), nil
}

func readStringFn(obj Any, env *environment.Env) (Any, error) {
	code, ok := obj.(String)
	if !ok {
		return nil, &ErrWrongType{obj}
	}
	reader := strings.NewReader(string(code))
	expr, err := parser.Parse(reader)
	if err != nil {
		return nil, err
	}

	switch len(expr) {
	case 0:
		return nil, nil
	case 1:
		return expr[0], nil
	default:
		return List(expr), nil
	}
}

func toInt(obj Any) (Any, error) {
	switch obj := obj.(type) {
	case Int:
		return obj, nil
	case Float:
		return Int(obj), nil
	case String:
		switch {
		case parser.IsInt(string(obj)):
			return parser.ParseInt(string(obj))
		case parser.IsFloat(string(obj)):
			f, err := parser.ParseFloat(string(obj))
			if err != nil {
				return nil, err
			}
			return Int(f), nil
		default:
			return nil, fmt.Errorf("cannot convert %v to int", obj)
		}
	default:
		return nil, fmt.Errorf("cannot convert %v of type %T to int", obj, obj)
	}
}

func toFloat(obj Any) (Any, error) {
	switch obj := obj.(type) {
	case Float:
		return obj, nil
	case Int:
		return Float(obj), nil
	case String:
		switch {
		case parser.IsFloat(string(obj)):
			return parser.ParseFloat(string(obj))
		default:
			return nil, fmt.Errorf("cannot convert %v to float", obj)
		}
	default:
		return nil, fmt.Errorf("cannot convert %v of type %T to float", obj, obj)
	}
}

func printEnv(env *environment.Env, depth int) {
	if env.Parent != nil {
		printEnv(env.Parent, depth-1)
	}

	var out []string
	for key, val := range env.Objects {
		out = append(out, fmt.Sprintf("%v => %v", key, val))
	}

	fmt.Printf("%d: { %v }\n", depth, strings.Join(out, ", "))
}

func toString(args []Any, env *environment.Env, sep String) (string, error) {
	if len(args) == 0 {
		return "", &ErrNumArgs{len(args)}
	}
	objs, err := evalAll(args, env)
	if err != nil {
		return "", err
	}

	var str []string
	for _, obj := range objs {
		str = append(str, fmt.Sprintf("%v", obj))
	}
	return strings.Join(str, string(sep)), nil
}

func isTrue(obj Any) bool {
	switch obj := obj.(type) {
	case Bool:
		return obj
	default:
		return obj != nil
	}
}

func last(args []Any) Any {
	if len(args) == 0 {
		return nil
	}
	return args[len(args)-1]
}
