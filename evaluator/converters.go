package evaluator

import (
	"fmt"

	"github.com/twolodzko/goal/parser"
	. "github.com/twolodzko/goal/types"
)

func toStringFn(obj Any) (Any, error) {
	switch obj := obj.(type) {
	case String:
		return obj, nil
	default:
		return String(fmt.Sprintf("%v", obj)), nil
	}
}

func toIntFn(obj Any) (Any, error) {
	switch obj := obj.(type) {
	case Int:
		return obj, nil
	case Float:
		return Int(obj), nil
	case String:
		switch {
		case parser.IsInt(string(obj)):
			return parser.ParseInt(string(obj))
		case parser.IsFloat(string(obj)):
			f, err := parser.ParseFloat(string(obj))
			if err != nil {
				return nil, err
			}
			return Int(f), nil
		default:
			return nil, fmt.Errorf("cannot convert %v to int", obj)
		}
	default:
		return nil, fmt.Errorf("cannot convert %v of type %T to int", obj, obj)
	}
}

func toFloatFn(obj Any) (Any, error) {
	switch obj := obj.(type) {
	case Float:
		return obj, nil
	case Int:
		return Float(obj), nil
	case String:
		switch {
		case parser.IsFloat(string(obj)):
			return parser.ParseFloat(string(obj))
		case parser.IsInt(string(obj)):
			i, err := parser.ParseInt(string(obj))
			if err != nil {
				return nil, err
			}
			return Float(i), nil
		default:
			return nil, fmt.Errorf("cannot convert %v to float", obj)
		}
	default:
		return nil, fmt.Errorf("cannot convert %v of type %T to float", obj, obj)
	}
}
