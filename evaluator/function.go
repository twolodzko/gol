package evaluator

import "github.com/twolodzko/gol/environment"

type function interface {
	Eval([]Any, *environment.Env) (Any, error)
}

type simpleFunction struct {
	fn func([]Any, *environment.Env) (Any, error)
}

func (f *simpleFunction) Eval(args []Any, env *environment.Env) (Any, error) {
	return f.fn(args, env)
}

type singleArgFunction struct {
	fn func(Any) (Any, error)
}

func (f *singleArgFunction) Eval(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 1 {
		return nil, &ErrNumArgs{len(args)}
	}
	obj, err := eval(args[0], env)
	if err != nil {
		return nil, err
	}
	return f.fn(obj)
}

type multiArgFunction struct {
	fn func([]Any) (Any, error)
}

func (f *multiArgFunction) Eval(args []Any, env *environment.Env) (Any, error) {
	objs, err := evalAll(args, env)
	if err != nil {
		return nil, err
	}
	return f.fn(objs)
}
