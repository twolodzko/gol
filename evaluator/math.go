package evaluator

import "github.com/twolodzko/gol/environment"

type ArithmeticFunction struct {
	intFn   func(x, y Int) Int
	floatFn func(x, y Float) Float
	start   Float
}

func (f *ArithmeticFunction) Call(args []Any, env *environment.Env) (Any, error) {
	args, err := EvalAll(args, env)
	if err != nil {
		return nil, err
	}
	return f.apply(args)
}

func (f *ArithmeticFunction) apply(arr []Any) (Any, error) {
	if len(arr) == 0 {
		return nil, &ErrNumArgs{len(arr)}
	}

	switch x := arr[0].(type) {
	case Int:
		if len(arr) == 1 {
			return f.intFn(Int(f.start), x), nil
		}
		acc := x
		for i, x := range arr[1:] {
			switch x := x.(type) {
			case Int:
				acc = f.intFn(acc, x)
			case Float:
				// fallback to floats
				return applyFloatFn(arr[1+i:], f.floatFn, Float(acc))
			default:
				return 0, &ErrNaN{x}
			}
		}
		return acc, nil

	case Float:
		if len(arr) == 1 {
			return f.floatFn(f.start, x), nil
		}
		return applyFloatFn(arr[1:], f.floatFn, x)

	default:
		return nil, &ErrNaN{x}
	}
}

func applyFloatFn(arr []Any, fn func(x, y Float) Float, start Float) (Float, error) {
	acc := start
	for _, x := range arr {
		switch x := x.(type) {
		case Float:
			acc = fn(acc, x)
		case Int:
			acc = fn(acc, Float(x))
		default:
			return 0, &ErrNaN{x}
		}
	}
	return acc, nil
}

func floatDivFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) == 0 {
		return nil, &ErrNumArgs{len(args)}
	}

	var start Float
	switch x := args[0].(type) {
	case Float:
		start = x
	case Int:
		start = Float(x)
	default:
		return nil, &ErrNaN{x}
	}

	if len(args) == 1 {
		return 1 / start, nil
	}

	return applyFloatFn(args[1:], func(x, y Float) Float { return x / y }, start)
}

func intDivFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) < 2 {
		return nil, &ErrNumArgs{len(args)}
	}

	var acc Int
	switch x := args[0].(type) {
	case Int:
		acc = x
	case Float:
		return nil, &ErrWrongType{x}
	default:
		return nil, &ErrNaN{x}
	}

	for _, x := range args[1:] {
		switch x := x.(type) {
		case Int:
			acc /= x
		case Float:
			return nil, &ErrWrongType{x}
		default:
			return nil, &ErrNaN{x}
		}
	}

	return acc, nil
}
