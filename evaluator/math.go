package evaluator

import (
	. "github.com/twolodzko/goal/types"
)

var (
	sumFn      = numSwitchFn(func(x, y Int) Int { return x + y }, func(x, y Float) Float { return x + y })
	difFn      = numSwitchFn(func(x, y Int) Int { return x - y }, func(x, y Float) Float { return x - y })
	mulFn      = numSwitchFn(func(x, y Int) Int { return x * y }, func(x, y Float) Float { return x * y })
	intDivFn   = intFn(func(x, y Int) Int { return x / y })
	floatDivFn = floatFn(func(x, y Float) Float { return x / y })
	modFn      = intFn(func(x, y Int) Int { return x % y })
)

func numSwitchFn(intFn func(Int, Int) Int, floatFn func(Float, Float) Float) buildIn {
	return func(obj []Any) (Any, error) {
		if len(obj) != 2 {
			return nil, &errNumArgs{len(obj)}
		}

		switch x := obj[0].(type) {
		case Int:
			switch y := obj[1].(type) {
			case Int:
				return intFn(x, y), nil
			case Float:
				return floatFn(Float(x), y), nil
			default:
				return nil, &errWrongType{y}
			}
		case Float:
			switch y := obj[1].(type) {
			case Float:
				return floatFn(x, y), nil
			case Int:
				return floatFn(x, Float(y)), nil
			default:
				return nil, &errWrongType{y}
			}
		default:
			return nil, &errWrongType{x}
		}
	}
}

func intFn(fn func(Int, Int) Int) buildIn {
	return func(obj []Any) (Any, error) {
		var x, y Int

		switch v := obj[0].(type) {
		case Int:
			x = v
		case Float:
			x = Int(v)
		default:
			return nil, &errWrongType{v}
		}

		switch v := obj[1].(type) {
		case Int:
			y = v
		case Float:
			y = Int(v)
		default:
			return nil, &errWrongType{v}
		}

		return fn(x, y), nil
	}
}

func floatFn(fn func(Float, Float) Float) buildIn {
	return func(obj []Any) (Any, error) {
		var x, y Float

		switch v := obj[0].(type) {
		case Float:
			x = v
		case Int:
			x = Float(v)
		default:
			return nil, &errWrongType{v}
		}

		switch v := obj[1].(type) {
		case Float:
			y = v
		case Int:
			y = Float(v)
		default:
			return nil, &errWrongType{v}
		}

		return fn(x, y), nil
	}
}
