package evaluator

import "github.com/twolodzko/gol/environment"

type function interface {
	Call([]Any, *environment.Env) (Any, error)
}

type simpleFunction struct {
	fn func([]Any, *environment.Env) (Any, error)
}

func (f *simpleFunction) Call(args []Any, env *environment.Env) (Any, error) {
	return f.fn(args, env)
}

type singleArgFunction struct {
	fn func(Any, *environment.Env) (Any, error)
}

func (f *singleArgFunction) Call(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 1 {
		return nil, &ErrNumArgs{len(args)}
	}
	obj, err := eval(args[0], env)
	if err != nil {
		return nil, err
	}
	return f.fn(obj, env)
}
