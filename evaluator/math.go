package evaluator

import (
	"math"

	. "github.com/twolodzko/goal/types"
)

var (
	sumFn      = accumulate(func(x, y Float) Float { return x + y }, 0)
	difFn      = accumulate(func(x, y Float) Float { return x - y }, 0)
	mulFn      = accumulate(func(x, y Float) Float { return x * y }, 1)
	intDivFn   = intFn(func(x, y Int) Int { return x / y })
	floatDivFn = accumulate(func(x, y Float) Float { return x / y }, 1)
	intModFn   = intFn(func(x, y Int) Int { return x % y })
	powFn      = accumulate(math.Pow, 1)
	remFn      = accumulate(math.Remainder, 1)
)

func toFloat(o Any) (Float, bool) {
	switch x := o.(type) {
	case Float:
		return x, true
	case Int:
		return Float(x), true
	default:
		return 0, false
	}
}

func accumulate(fn func(Float, Float) Float, start Float) buildIn {
	return func(obj []Any) (Any, error) {
		var acc Float = start

		if len(obj) == 0 {
			return nil, &errNumArgs{len(obj)}
		}

		x, ok := toFloat(obj[0])
		if !ok {
			return 0, &errWrongType{obj[0]}
		}

		if len(obj) == 1 {
			return fn(acc, x), nil
		}

		acc = x
		for _, x := range obj[1:] {
			f, ok := toFloat(x)
			if !ok {
				return 0, &errWrongType{x}
			}
			acc = fn(acc, f)
		}
		return acc, nil
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
