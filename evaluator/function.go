package evaluator

import (
	"github.com/twolodzko/gol/environment"
	. "github.com/twolodzko/gol/types"
)

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
	for i, val := range args {
		f.env.Set(f.args[i], val)
	}
	objs, err := EvalAll(f.expr, f.env)
	return last(objs), err
}

func NewLambda(args []Any, env *environment.Env) (*Lambda, error) {
	if len(args) < 2 {
		return &Lambda{}, &ErrNumArgs{len(args)}
	}
	localEnv := environment.NewEnv(env)
	argList, ok := args[0].(List)
	if !ok {
		return &Lambda{}, &ErrWrongType{args[0]}
	}
	argNames, err := toSymbols(argList)
	if err != nil {
		return &Lambda{}, err
	}
	return &Lambda{localEnv, argNames, args[1:]}, nil
}

func toSymbols(objs List) ([]Symbol, error) {
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
