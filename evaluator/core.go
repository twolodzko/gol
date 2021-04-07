package evaluator

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/gol/environment"
	"github.com/twolodzko/gol/parser"
)

func isTrue(obj Any) bool {
	switch obj := obj.(type) {
	case Bool:
		return obj
	default:
		return obj != nil
	}
}

func defFn(args []Any, env *environment.Env) (Any, error) {

	if len(args) != 2 {
		return nil, &ErrNumArgs{len(args)}
	}

	name, ok := args[0].(Symbol)
	if !ok {
		return nil, &ErrWrongType{args[0]}
	}

	val, err := Eval(args[1], env)
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

	val, err := Eval(args[1], env)
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
	objs, err := EvalAll(args, env)
	return last(objs), err
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
	objs, err := EvalAll(args, env)
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
	objs, err := EvalAll(args, env)
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
}

func andFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) == 0 {
		return Bool(false), nil
	}
	for _, arg := range args {
		arg, err := Eval(arg, env)
		if err != nil {
			return nil, err
		}
		if !isTrue(arg) {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func orFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) == 0 {
		return Bool(false), nil
	}
	for _, arg := range args {
		obj, err := Eval(arg, env)
		if err != nil {
			return nil, err
		}
		if isTrue(obj) {
			return Bool(true), nil
		}
	}
	return Bool(false), nil
}

func equalFn(objs []Any) (Any, error) {
	switch first := objs[0].(type) {
	case Bool:
		second, ok := objs[1].(Bool)
		if !ok {
			return Bool(false), nil
		}
		return Bool(first == second), nil
	case Int:
		switch second := objs[1].(type) {
		case Int:
			return Bool(first == second), nil
		case Float:
			return Bool(Float(first) == second), nil
		default:
			return Bool(false), nil
		}
	case Float:
		switch second := objs[1].(type) {
		case Float:
			return Bool(first == second), nil
		case Int:
			return Bool(first == Float(second)), nil
		default:
			return Bool(false), nil
		}
	case String:
		second, ok := objs[1].(String)
		if !ok {
			return Bool(false), nil
		}
		return Bool(first == second), nil
	case Function:
		second, ok := objs[1].(Function)
		if !ok {
			return Bool(false), nil
		}
		return Bool(first == second), nil
	default:
		return Bool(cmp.Equal(objs[0], objs[1])), nil
	}
}

func gtFn(objs []Any) (Any, error) {
	switch first := objs[0].(type) {
	case Int:
		switch second := objs[1].(type) {
		case Int:
			return Bool(first > second), nil
		case Float:
			return Bool(Float(first) > second), nil
		default:
			return nil, &ErrWrongType{objs[1]}
		}
	case Float:
		switch second := objs[1].(type) {
		case Float:
			return Bool(first > second), nil
		case Int:
			return Bool(first > Float(second)), nil
		default:
			return nil, &ErrWrongType{objs[1]}
		}
	default:
		return nil, &ErrWrongType{objs[0]}
	}
}

func ltFn(objs []Any) (Any, error) {
	switch first := objs[0].(type) {
	case Int:
		switch second := objs[1].(type) {
		case Int:
			return Bool(first < second), nil
		case Float:
			return Bool(Float(first) < second), nil
		default:
			return nil, &ErrWrongType{objs[1]}
		}
	case Float:
		switch second := objs[1].(type) {
		case Float:
			return Bool(first < second), nil
		case Int:
			return Bool(first < Float(second)), nil
		default:
			return nil, &ErrWrongType{objs[1]}
		}
	default:
		return nil, &ErrWrongType{objs[0]}
	}
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

func toString(args []Any, env *environment.Env, sep String) (string, error) {
	if len(args) == 0 {
		return "", &ErrNumArgs{len(args)}
	}
	objs, err := EvalAll(args, env)
	if err != nil {
		return "", err
	}

	var str []string
	for _, obj := range objs {
		str = append(str, fmt.Sprintf("%v", obj))
	}
	return strings.Join(str, string(sep)), nil
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

func last(args []Any) Any {
	if len(args) == 0 {
		return nil
	}
	return args[len(args)-1]
}
