package evaluator

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/gol/environment"
	"github.com/twolodzko/gol/parser"
)

func defFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) < 2 {
		return nil, &ErrNumArgs{len(args)}
	}

	switch first := args[0].(type) {
	case Symbol:
		// simple assignments
		if len(args) != 2 {
			return nil, &ErrNumArgs{len(args)}
		}
		val, err := eval(args[1], env)
		if err != nil {
			return nil, err
		}
		env.Set(first, val)
		return val, err
	case List:
		// named functions shorthand notation
		if len(first) < 1 {
			return nil, &ErrNumArgs{len(first)}
		}
		name, ok := first.Head().(Symbol)
		if !ok {
			return nil, &ErrWrongType{first.Head()}
		}
		fn, err := newLambda(first.Tail(), args[1:], env)
		if err != nil {
			return nil, err
		}
		env.Set(name, fn)
		return fn, nil
	default:
		return nil, &ErrWrongType{first}
	}
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

func applyFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 2 {
		return nil, &ErrNumArgs{len(args)}
	}

	obj, err := eval(args[1], env)
	if err != nil {
		return nil, err
	}
	l, ok := obj.(List)
	if !ok {
		return nil, &ErrWrongType{args[1]}
	}

	expr := append(List{args[0]}, l...)
	return eval(expr, env)
}

func mapFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 2 {
		return nil, &ErrNumArgs{len(args)}
	}

	obj, err := eval(args[1], env)
	if err != nil {
		return nil, err
	}
	l, ok := obj.(List)
	if !ok {
		return nil, &ErrWrongType{args[1]}
	}

	expr := List{args[0]}
	var out List
	for _, arg := range l {
		res, err := eval(append(expr, arg), env)
		if err != nil {
			return nil, err
		}
		out = append(out, res)
	}

	return out, nil
}

func quasiquote(arg Any, env *environment.Env) (Any, error) {
	var list List

	switch obj := arg.(type) {
	case List:
		if len(obj) == 0 {
			return obj, nil
		}
		switch obj[0] {
		case Symbol("unquote"):
			if len(obj) != 2 {
				return nil, errors.New("nothing to unquote")
			}
			return eval(obj[1], env)
		case Symbol("quasiquote"):
			return obj, nil
		default:
			for _, o := range obj {
				val, err := quasiquote(o, env)
				if err != nil {
					return list, err
				}
				list = append(list, val)
			}
			return list, nil
		}
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
	if len(args) != 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	objs, err := evalAll(args, env)
	if err != nil {
		return nil, err
	}
	l, ok := objs[1].(List)
	if !ok {
		return nil, &ErrWrongType{args[1]}
	}
	return List(append([]Any{objs[0]}, l...)), nil
}

func concatFn(objs []Any) (Any, error) {
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
		default:
			if !cmp.Equal(first, second) {
				return Bool(false), nil
			}
		}
	}

	return Bool(true), nil
}

func parseStringFn(obj Any) (Any, error) {
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

func toString(objs []Any, sep String) (string, error) {
	if len(objs) == 0 {
		return "", &ErrNumArgs{len(objs)}
	}

	var str []string
	for _, obj := range objs {
		switch obj := obj.(type) {
		case String:
			str = append(str, obj.Raw())
		default:
			str = append(str, fmt.Sprintf("%v", obj))
		}
	}
	return strings.Join(str, string(sep)), nil
}

func writeToFileFn(objs []Any) (Any, error) {
	if len(objs) != 2 {
		return nil, &ErrNumArgs{len(objs)}
	}
	fileName, ok := objs[0].(String)
	if !ok {
		return nil, &ErrWrongType{objs[0]}
	}

	file, err := os.OpenFile(string(fileName), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var str string
	switch obj := objs[1].(type) {
	case String:
		str = obj.Raw()
	default:
		str = fmt.Sprintf("%v", obj)
	}
	_, err = fmt.Fprintln(file, str)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func reverse(arr []Any) []Any {
	n := len(arr)
	out := make([]Any, n)
	copy(out, arr)
	for i := 0; i < n/2; i++ {
		j := n - i - 1
		out[i], out[j] = out[j], out[i]
	}
	return out
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
