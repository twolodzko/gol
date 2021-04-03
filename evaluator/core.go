package evaluator

import (
	"fmt"

	"github.com/twolodzko/goal/environment"
	. "github.com/twolodzko/goal/types"
)

func isTrue(obj Any) bool {
	return obj != nil && obj != false
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

func defFn(args []Any, env *environment.Env) error {
	var err error

	if len(args) != 2 {
		return &ErrNumArgs{len(args)}
	}

	name, ok := args[0].(Symbol)
	if !ok {
		return &ErrWrongType{args[0]}
	}

	val, err := Eval(args[1], env)
	if err != nil {
		return err
	}

	err = env.Set(name, val)

	return err
}

func headFn(args []Any, env *environment.Env) (Any, error) {
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

func tailFn(args []Any, env *environment.Env) (Any, error) {
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

func apply(args []Any, fn func(Any) Any) Any {
	var out []Any
	for _, x := range args {
		out = append(out, fn(x))
	}

	switch len(out) {
	case 0:
		return nil
	case 1:
		return out[0]
	default:
		return List(out)
	}
}

func printFn(args []Any) {
	out := ""
	for _, o := range args {
		out += fmt.Sprintf("%v", o)
	}
	fmt.Printf("%s\n", out)
}
