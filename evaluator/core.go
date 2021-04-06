package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/gol/environment"
	"github.com/twolodzko/gol/parser"
	. "github.com/twolodzko/gol/types"
)

func isTrue(obj Any) bool {
	switch obj := obj.(type) {
	case Bool:
		return obj
	default:
		return obj != nil
	}
}

func ifFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 3 {
		return nil, &ErrNumArgs{len(args)}
	}

	cond, err := Eval(args[0], env)
	if err != nil {
		return nil, err
	}

	if isTrue(cond) {
		return Eval(args[1], env)
	}
	return Eval(args[2], env)
}

func defFn(args []Any, env *environment.Env) (Any, error) {
	var err error

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

func letFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) < 2 {
		return nil, &ErrNumArgs{len(args)}
	}

	vars, ok := args[0].(List)
	if !ok {
		return nil, &ErrWrongType{args[0]}
	}
	if len(vars) != 2 {
		return nil, fmt.Errorf("invalid variable assignment %v", vars)
	}

	localEnv := environment.NewEnclosedEnv(env)

	name, ok := vars[0].(Symbol)
	if !ok {
		return nil, &ErrWrongType{vars[0]}
	}

	localEnv.Set(name, vars[1])

	res, err := EvalAll(args[1:], localEnv)

	if len(res) == 0 {
		return nil, err
	}

	return last(res), err
}

func firstFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 1 {
		return nil, &ErrNumArgs{len(args)}
	}
	obj, err := Eval(args[0], env)
	if err != nil {
		return nil, err
	}

	l, ok := obj.(List)
	if !ok {
		return nil, fmt.Errorf("%v is not a list", obj)
	}
	return l.Head(), nil
}

func restFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 1 {
		return nil, &ErrNumArgs{len(args)}
	}
	obj, err := Eval(args[0], env)
	if err != nil {
		return nil, err
	}

	l, ok := obj.(List)
	if !ok {
		return nil, fmt.Errorf("%v is not a list", obj)
	}
	return l.Tail(), nil
}

func initFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 1 {
		return nil, &ErrNumArgs{len(args)}
	}
	obj, err := Eval(args[0], env)
	if err != nil {
		return nil, err
	}

	l, ok := obj.(List)
	if !ok {
		return nil, fmt.Errorf("%v is not a list", obj)
	}
	if len(l) == 0 {
		return nil, nil
	}
	return l[:len(l)-1], nil
}

func lastFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 1 {
		return nil, &ErrNumArgs{len(args)}
	}
	obj, err := Eval(args[0], env)
	if err != nil {
		return nil, err
	}

	l, ok := obj.(List)
	if !ok {
		return nil, fmt.Errorf("%v is not a list", obj)
	}
	if len(l) == 0 {
		return nil, nil
	}
	return l[len(l)-1], nil
}

func nthFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	args, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}

	switch l := args[0].(type) {
	case List:
		var n int
		switch val := args[1].(type) {
		case Int:
			n = val
		case Float:
			n = int(val)
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
	objs, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}
	if len(objs) == 0 {
		return Bool(false), nil
	}

	for _, x := range objs {
		if !isTrue(x) {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func orFn(args []Any, env *environment.Env) (Any, error) {
	objs, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}
	if len(objs) == 0 {
		return Bool(false), nil
	}

	for _, x := range objs {
		if isTrue(x) {
			return Bool(true), nil
		}
	}
	return Bool(false), nil
}

func equalFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	objs, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}

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

func gtFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	objs, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}

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

func ltFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	objs, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}

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

func readStringFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 1 {
		return nil, &ErrNumArgs{len(args)}
	}
	obj, err := Eval(args[0], env)
	if err != nil {
		return nil, err
	}

	code, ok := obj.(String)
	if !ok {
		return nil, &ErrWrongType{obj}
	}
	expr, err := parser.Parse(strings.NewReader(string(code)))
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

func readFile(name String) (String, error) {
	var lines []string

	file, err := os.Open(string(name))
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return String(strings.Join(lines, "\n")), nil
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
