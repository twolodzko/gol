package evaluator

import (
	. "github.com/twolodzko/goal/types"
)

// var (
// 	FloatSum = floatAccumulate(func(x, y Float) Float { return x + y }, 0)
// 	IntSum   = intAccumulate(func(x, y Int) Int { return x + y }, 0)
// 	FloatDif = floatAccumulate(func(x, y Float) Float { return x - y }, 0)
// 	IntDif   = intAccumulate(func(x, y Int) Int { return x - y }, 0)
// 	FloatMul = floatAccumulate(func(x, y Float) Float { return x * y }, 1)
// 	IntMul   = intAccumulate(func(x, y Int) Int { return x * y }, 1)
// 	FloatDiv = floatAccumulate(func(x, y Float) Float { return x / y }, 1)
// 	FloatMod = floatAccumulate(math.Mod, 1)
// 	IntMod   = intAccumulate(func(x, y Int) Int { return x % y }, 1)
// 	Pow      = floatAccumulate(math.Pow, 1)
// 	Rem      = floatAccumulate(math.Remainder, 1)
// )

func toFloat(x Any) (Float, bool) {
	switch x := x.(type) {
	case Float:
		return x, true
	case Int:
		return Float(x), true
	default:
		return 0, false
	}
}

func accumulate(obj []Any, fn func(Float, Float) Float, start Float) (Any, error) {
	var acc Float = start

	if len(obj) == 0 {
		return nil, &ErrNumArgs{len(obj)}
	}

	x, ok := toFloat(obj[0])
	if !ok {
		return 0, &ErrWrongType{obj[0]}
	}

	if len(obj) == 1 {
		return fn(acc, x), nil
	}

	acc = x
	for _, x := range obj[1:] {
		f, ok := toFloat(x)
		if !ok {
			return 0, &ErrWrongType{x}
		}
		acc = fn(acc, f)
	}
	return acc, nil
}

// func toInt(o Any) (Int, bool) {
// 	switch x := o.(type) {
// 	case Int:
// 		return x, true
// 	case Float:
// 		return Int(x), true
// 	default:
// 		return 0, false
// 	}
// }

// func intAccumulate(fn func(Int, Int) Int, start Int) Buildin {
// 	return func(obj []Any) (Any, error) {
// 		var acc Int = start

// 		if len(obj) == 0 {
// 			return nil, &ErrNumArgs{len(obj)}
// 		}

// 		x, ok := toInt(obj[0])
// 		if !ok {
// 			return 0, &ErrWrongType{obj[0]}
// 		}

// 		if len(obj) == 1 {
// 			return fn(acc, x), nil
// 		}

// 		acc = x
// 		for _, x := range obj[1:] {
// 			f, ok := toInt(x)
// 			if !ok {
// 				return 0, &ErrWrongType{x}
// 			}
// 			acc = fn(acc, f)
// 		}
// 		return acc, nil
// 	}
// }

// func IntDiv(obj []Any) (Any, error) {
// 	if len(obj) < 2 {
// 		return nil, &ErrNumArgs{len(obj)}
// 	}

// 	acc, ok := toInt(obj[0])
// 	if !ok {
// 		return 0, &ErrWrongType{obj[0]}
// 	}

// 	for _, x := range obj[1:] {
// 		i, ok := toInt(x)
// 		if !ok {
// 			return 0, &ErrWrongType{x}
// 		}
// 		if x == 0 {
// 			return 0, errors.New("integer divide by zero")
// 		}

// 		acc /= i
// 	}
// 	return acc, nil
// }
