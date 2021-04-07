package evaluator

import (
	"fmt"

	"github.com/twolodzko/gol/environment"
)

type TailCallOptimized interface {
	Call([]Any, *environment.Env) (Any, *environment.Env, error)
}

type Lambda struct {
	env  *environment.Env
	args []Symbol
	expr []Any
}

func (f *Lambda) Call(args []Any, env *environment.Env) (Any, *environment.Env, error) {
	var (
		err  error
		objs []Any
	)
	if len(args) != len(f.args) {
		return nil, env, &ErrNumArgs{len(args)}
	}
	objs, err = EvalAll(args, env)
	if err != nil {
		return nil, env, err
	}

	localEnv := environment.NewEnv(f.env)
	for i, val := range objs {
		localEnv.Set(f.args[i], val)
	}
	_, err = EvalAll(exceptLast(f.expr), localEnv)
	return last(f.expr), localEnv, err
}

func NewLambda(args []Any, env *environment.Env) (*Lambda, error) {
	if len(args) < 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	argList, ok := args[0].(List)
	if !ok {
		return nil, &ErrWrongType{args[0]}
	}
	argNames, err := areSymbols(argList)
	if err != nil {
		return nil, err
	}
	return &Lambda{env, argNames, args[1:]}, nil
}

func areSymbols(objs List) ([]Symbol, error) {
	var symbols []Symbol
	for _, obj := range objs {
		s, ok := obj.(Symbol)
		if !ok {
			return symbols, &ErrWrongType{obj}
		}
		symbols = append(symbols, s)
	}
	return symbols, nil
}

type TcoFunction struct {
	fn func([]Any, *environment.Env) (Any, *environment.Env, error)
}

func (f *TcoFunction) Call(args []Any, env *environment.Env) (Any, *environment.Env, error) {
	return f.fn(args, env)
}

func ifFn(args []Any, env *environment.Env) (Any, *environment.Env, error) {
	if len(args) != 3 {
		return nil, env, &ErrNumArgs{len(args)}
	}

	cond, err := Eval(args[0], env)
	if err != nil {
		return nil, env, err
	}

	if isTrue(cond) {
		return args[1], env, nil
	}
	return args[2], env, nil
}

func letFn(args []Any, env *environment.Env) (Any, *environment.Env, error) {
	if len(args) < 2 {
		return nil, env, &ErrNumArgs{len(args)}
	}

	localEnv := environment.NewEnv(env)

	bindings, ok := args[0].(List)
	if !ok {
		return nil, env, &ErrWrongType{args[0]}
	}

	n := len(bindings)
	if (n % 2) != 0 {
		return nil, env, fmt.Errorf("invalid variable bindings %v", bindings)
	}

	// odd entries are keys, even are the values
	// e.g. (let (x 1 y 2) (+ x y)) => 3
	for i := 0; i < n; i += 2 {
		name, ok := bindings[i].(Symbol)
		if !ok {
			return nil, env, &ErrWrongType{bindings[0]}
		}
		val, err := Eval(bindings[i+1], localEnv)
		if err != nil {
			return nil, env, err
		}
		localEnv.Set(name, val)
	}

	_, err := EvalAll(exceptLast(args[1:]), localEnv)
	return last(args), localEnv, err
}

func condFun(args []Any, env *environment.Env) (Any, *environment.Env, error) {
	for _, arg := range args {
		obj, ok := arg.(List)
		if !ok || len(obj) == 0 {
			return nil, env, fmt.Errorf("invalid condition: %v (%T)", arg, arg)
		}

		cond, err := Eval(obj[0], env)
		if err != nil {
			return nil, env, err
		}
		if isTrue(cond) {
			if len(obj) > 1 {
				_, err := EvalAll(exceptLast(obj[1:]), env)
				return last(obj), env, err
			}
			return nil, env, nil
		}
	}
	return nil, env, nil
}

func exceptLast(objs []Any) []Any {
	if len(objs) > 1 {
		return objs[:len(objs)-1]
	}
	return nil
}
