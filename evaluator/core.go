package evaluator

import (
	"fmt"
	"strings"

	"github.com/twolodzko/goal/environment"
	. "github.com/twolodzko/goal/types"
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

func andFn(args []Any) Bool {
	if len(args) == 0 {
		return Bool(false)
	}

	for _, x := range args {
		if !isTrue(x) {
			return Bool(false)
		}
	}
	return Bool(true)
}

func orFn(args []Any) Bool {
	if len(args) == 0 {
		return Bool(false)
	}

	for _, x := range args {
		if isTrue(x) {
			return Bool(true)
		}
	}
	return Bool(false)
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
