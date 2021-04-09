package evaluator

import (
	"fmt"

	"github.com/twolodzko/gol/environment"
)

func getFloat(obj Any, env *environment.Env) (Float, error) {
	obj, err := eval(obj, env)
	if err != nil {
		return 0, err
	}
	switch obj := obj.(type) {
	case Float:
		return obj, nil
	case Int:
		return Float(obj), nil
	default:
		return 0, &ErrNaN{obj}
	}
}

type singleArgFloatFunction struct {
	fn func(Float) Float
}

func (f *singleArgFloatFunction) Eval(args []Any, env *environment.Env) (Any, error) {
	if len(args) != 1 {
		return nil, &ErrNumArgs{len(args)}
	}
	num, err := getFloat(args[0], env)
	if err != nil {
		return nil, err
	}
	return f.fn(num), nil
}

type multiArgFloatFunction struct {
	fn    func(x, y Float) Float
	start Float
}

func (f *multiArgFloatFunction) Eval(args []Any, env *environment.Env) (Any, error) {
	if len(args) == 0 {
		return f.start, nil
	}

	num, err := getFloat(args[0], env)
	if err != nil {
		return nil, err
	}

	if len(args) == 1 {
		return f.fn(f.start, num), nil
	}

	res := num
	for _, obj := range args[1:] {
		num, err := getFloat(obj, env)
		if err != nil {
			return nil, err
		}
		res = f.fn(res, num)
	}
	return res, nil
}

func getInt(obj Any, env *environment.Env) (Int, error) {
	obj, err := eval(obj, env)
	if err != nil {
		return 0, err
	}
	switch obj := obj.(type) {
	case Int:
		return obj, nil
	default:
		return 0, fmt.Errorf("%v (%T) is not an int", obj, obj)
	}
}

type multiArgIntFunction struct {
	fn    func(x, y Int) Int
	start Int
}

func (f *multiArgIntFunction) Eval(args []Any, env *environment.Env) (Any, error) {
	if len(args) == 0 {
		return f.start, nil
	}

	num, err := getInt(args[0], env)
	if err != nil {
		return nil, err
	}

	if len(args) == 1 {
		return f.fn(f.start, num), nil
	}

	res := num
	for _, obj := range args[1:] {
		num, err := getInt(obj, env)
		if err != nil {
			return nil, err
		}
		res = f.fn(res, num)
	}
	return res, nil
}

func gtFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) < 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	first, err := eval(args[0], env)
	if err != nil {
		return nil, err
	}

	for _, second := range args[1:] {
		second, err := eval(second, env)
		if err != nil {
			return nil, err
		}

		switch first := first.(type) {
		case Int:
			switch second := second.(type) {
			case Int:
				if first <= second {
					return Bool(false), nil
				}
			case Float:
				if Float(first) <= second {
					return Bool(false), nil
				}
			default:
				return nil, &ErrWrongType{second}
			}
		case Float:
			switch second := second.(type) {
			case Float:
				if first <= second {
					return Bool(false), nil
				}
			case Int:
				if first <= Float(second) {
					return Bool(false), nil
				}
			default:
				return nil, &ErrWrongType{second}
			}
		default:
			return nil, &ErrWrongType{first}
		}

		first = second
	}

	return Bool(true), nil
}

func ltFn(args []Any, env *environment.Env) (Any, error) {
	if len(args) < 2 {
		return nil, &ErrNumArgs{len(args)}
	}
	first, err := eval(args[0], env)
	if err != nil {
		return nil, err
	}

	for _, second := range args[1:] {
		switch first := first.(type) {
		case Int:
			switch second := second.(type) {
			case Int:
				if first >= second {
					return Bool(false), nil
				}
			case Float:
				if Float(first) >= second {
					return Bool(false), nil
				}
			default:
				return nil, &ErrWrongType{second}
			}
		case Float:
			switch second := second.(type) {
			case Float:
				if first >= second {
					return Bool(false), nil
				}
			case Int:
				if first >= Float(second) {
					return Bool(false), nil
				}
			default:
				return nil, &ErrWrongType{second}
			}
		default:
			return nil, &ErrWrongType{first}
		}

		first = second
	}

	return Bool(true), nil
}
