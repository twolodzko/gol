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
