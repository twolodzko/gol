package evaluator

import (
	"fmt"

	"github.com/twolodzko/gol/environment"
)

type tailCallOptimized interface {
	function
	PartialEval([]Any, *environment.Env) (Any, *environment.Env, error)
}

type lambda struct {
	env  *environment.Env
	args []Symbol
	expr []Any
}

func (f *lambda) Eval(args []Any, env *environment.Env) (Any, error) {
	expr, env, err := f.PartialEval(args, env)
	if err != nil {
		return nil, err
	}
	return eval(expr, env)
}

func (f *lambda) PartialEval(args []Any, env *environment.Env) (Any, *environment.Env, error) {
	var (
		err  error
		objs []Any
	)
	if len(args) != len(f.args) {
		return nil, env, &ErrNumArgs{len(args)}
	}
	objs, err = evalAll(args, env)
	if err != nil {
		return nil, env, err
	}

	localEnv := environment.NewEnv(f.env)
	for i, val := range objs {
		localEnv.Set(f.args[i], val)
	}
	_, err = evalAll(exceptLast(f.expr), localEnv)
	return last(f.expr), localEnv, err
}

func newLambda(args Any, body []Any, env *environment.Env) (*lambda, error) {
	argList, ok := args.(List)
	if !ok {
		return nil, &ErrWrongType{args}
	}
	argNames, err := areSymbols(argList)
	if err != nil {
		return nil, err
	}
	return &lambda{env, argNames, body}, nil
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

type tcoFunction struct {
	fn func([]Any, *environment.Env) (Any, *environment.Env, error)
}

func (f *tcoFunction) Eval(args []Any, env *environment.Env) (Any, error) {
	expr, env, err := f.PartialEval(args, env)
	if err != nil {
		return nil, err
	}
	return eval(expr, env)
}

func (f *tcoFunction) PartialEval(args []Any, env *environment.Env) (Any, *environment.Env, error) {
	return f.fn(args, env)
}

func ifFn(args []Any, env *environment.Env) (Any, *environment.Env, error) {
	if len(args) != 3 {
		return nil, env, &ErrNumArgs{len(args)}
	}

	cond, err := eval(args[0], env)
	if err != nil {
		return nil, env, err
	}

	if isTrue(cond) {
		return args[1], env, nil
	}
	return args[2], env, nil
}

func letFn(args []Any, env *environment.Env) (Any, *environment.Env, error) {
	var err error

	if len(args) < 2 {
		return nil, env, &ErrNumArgs{len(args)}
	}

	localEnv := environment.NewEnv(env)

	bindings, ok := args[0].(List)
	if !ok {
		return nil, env, &ErrWrongType{args[0]}
	}
	_, err = setVariables(bindings, localEnv)
	if err != nil {
		return nil, localEnv, err
	}

	_, err = evalAll(exceptLast(args[1:]), localEnv)
	return last(args), localEnv, err
}

func condFun(args []Any, env *environment.Env) (Any, *environment.Env, error) {
	for _, arg := range args {
		obj, ok := arg.(List)
		if !ok || len(obj) == 0 {
			return nil, env, fmt.Errorf("invalid condition: %v (%T)", arg, arg)
		}

		cond, err := eval(obj[0], env)
		if err != nil {
			return nil, env, err
		}
		if isTrue(cond) {
			if len(obj) > 1 {
				_, err := evalAll(exceptLast(obj[1:]), env)
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

// odd entries are keys, even are the values
func setVariables(args []Any, env *environment.Env) (Any, error) {
	var val Any

	n := len(args)
	if (n % 2) != 0 {
		return nil, fmt.Errorf("invalid variable bindings %v", args)
	}

	for i := 0; i < n; i += 2 {
		name, ok := args[i].(Symbol)
		if !ok {
			return nil, &ErrWrongType{args[i]}
		}
		val, err := eval(args[i+1], env)
		if err != nil {
			return nil, err
		}
		env.Set(name, val)
	}

	return val, nil
}
