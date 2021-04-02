package evaluator

import (
	"errors"
	"math"

	. "github.com/twolodzko/goal/types"
)

var (
	floatSumFn = floatAccumulate(func(x, y Float) Float { return x + y }, 0)
	intSumFn   = intAccumulate(func(x, y Int) Int { return x + y }, 0)
	floatDifFn = floatAccumulate(func(x, y Float) Float { return x - y }, 0)
	intDifFn   = intAccumulate(func(x, y Int) Int { return x - y }, 0)
	floatMulFn = floatAccumulate(func(x, y Float) Float { return x * y }, 1)
	intMulFn   = intAccumulate(func(x, y Int) Int { return x * y }, 1)
	floatDivFn = floatAccumulate(func(x, y Float) Float { return x / y }, 1)
	floatModFn = floatAccumulate(math.Mod, 1)
	intModFn   = intAccumulate(func(x, y Int) Int { return x % y }, 1)
	powFn      = floatAccumulate(math.Pow, 1)
	remFn      = floatAccumulate(math.Remainder, 1)
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

func floatAccumulate(fn func(Float, Float) Float, start Float) BuildIn {
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

func toInt(o Any) (Int, bool) {
	switch x := o.(type) {
	case Int:
		return x, true
	case Float:
		return Int(x), true
	default:
		return 0, false
	}
}

func intAccumulate(fn func(Int, Int) Int, start Int) BuildIn {
	return func(obj []Any) (Any, error) {
		var acc Int = start

		if len(obj) == 0 {
			return nil, &errNumArgs{len(obj)}
		}

		x, ok := toInt(obj[0])
		if !ok {
			return 0, &errWrongType{obj[0]}
		}

		if len(obj) == 1 {
			return fn(acc, x), nil
		}

		acc = x
		for _, x := range obj[1:] {
			f, ok := toInt(x)
			if !ok {
				return 0, &errWrongType{x}
			}
			acc = fn(acc, f)
		}
		return acc, nil
	}
}

func intDivFn(obj []Any) (Any, error) {
	if len(obj) < 2 {
		return nil, &errNumArgs{len(obj)}
	}

	acc, ok := toInt(obj[0])
	if !ok {
		return 0, &errWrongType{obj[0]}
	}

	for _, x := range obj[1:] {
		i, ok := toInt(x)
		if !ok {
			return 0, &errWrongType{x}
		}
		if x == 0 {
			return 0, errors.New("integer divide by zero")
		}

		acc /= i
	}
	return acc, nil
}
