package evaluator

import "github.com/twolodzko/gol/environment"

type Function interface {
	Call([]Any, *environment.Env) (Any, error)
}

type SimpleFunction struct {
	fn func([]Any, *environment.Env) (Any, error)
}

func (f *SimpleFunction) Call(args []Any, env *environment.Env) (Any, error) {
	return f.fn(args, env)
}

type SingleArgFunction struct {
	fn func(Any, *environment.Env) (Any, error)
}

func (f *SingleArgFunction) Call(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 1 {
		return nil, &ErrNumArgs{len(args)}
	}
	obj, err := Eval(args[0], env)
	if err != nil {
		return nil, err
	}
	return f.fn(obj, env)
}

type Lambda struct {
	env  *environment.Env
	args []Symbol
	expr []Any
}

func (f *Lambda) Call(args []Any, env *environment.Env) (Any, error) {
	if len(args) != len(f.args) {
		return nil, &ErrNumArgs{len(args)}
	}
	args, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}

	localEnv := environment.NewEnv(f.env)
	for i, val := range args {
		localEnv.Set(f.args[i], val)
	}
	objs, err := EvalAll(f.expr, localEnv)
	return last(objs), err
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
