package evaluator

import (
	"fmt"

	"github.com/twolodzko/goal/parser"
	. "github.com/twolodzko/goal/types"
)

func toString(expr List) (Any, error) {
	switch obj := expr[0].(type) {
	case String:
		return obj, nil
	default:
		return String(fmt.Sprintf("%v", obj)), nil
	}
}

func floatToInt(f Float) Int {
	return Int(int(f))
}

func toInt(expr List) (Any, error) {
	switch obj := expr[0].(type) {
	case Int:
		return obj, nil
	case Float:
		return floatToInt(obj), nil
	case String:
		switch {
		case parser.IsInt(string(obj)):
			return parser.ParseInt(string(obj))
		case parser.IsFloat(string(obj)):
			f, err := parser.ParseFloat(string(obj))
			if err != nil {
				return nil, err
			}
			return floatToInt(f), nil
		default:
			return nil, fmt.Errorf("cannot parse %v to int", obj)
		}
	default:
		return nil, fmt.Errorf("cannot convert object of type %T to int", obj)
	}
}

func intToFloat(i Int) Float {
	return Float(float64(i))
}

func toFloat(expr List) (Any, error) {
	switch obj := expr[0].(type) {
	case Float:
		return obj, nil
	case Int:
		return intToFloat(obj), nil
	case String:
		switch {
		case parser.IsFloat(string(obj)):
			return parser.ParseFloat(string(obj))
		case parser.IsInt(string(obj)):
			i, err := parser.ParseInt(string(obj))
			if err != nil {
				return nil, err
			}
			return intToFloat(i), nil
		default:
			return nil, fmt.Errorf("cannot parse %v to float", obj)
		}
	default:
		return nil, fmt.Errorf("cannot convert object of type %T to float", obj)
	}
}
