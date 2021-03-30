package evaluator

import (
	"fmt"

	"github.com/twolodzko/goal/objects"
	"github.com/twolodzko/goal/parser"
)

func toString(expr []objects.Object) (objects.Object, error) {
	switch obj := expr[0].(type) {
	case objects.String:
		return obj, nil
	default:
		return objects.String{Val: obj.String()}, nil
	}
}

func floatToInt(f objects.Float) objects.Int {
	return objects.Int{Val: int(f.Val)}
}

func toInt(expr []objects.Object) (objects.Object, error) {
	switch obj := expr[0].(type) {
	case objects.Int:
		return obj, nil
	case objects.Float:
		return floatToInt(obj), nil
	case objects.String:
		switch {
		case parser.IsInt(obj.Val):
			return parser.ParseInt(obj.Val)
		case parser.IsFloat(obj.Val):
			f, err := parser.ParseFloat(obj.Val)
			if err != nil {
				return nil, err
			}
			return floatToInt(f), nil
		default:
			return nil, fmt.Errorf("cannot parse %v to int", obj.Val)
		}
	default:
		return nil, fmt.Errorf("cannot convert object of type %T to int", obj)
	}
}

func intToFloat(i objects.Int) objects.Float {
	return objects.Float{Val: float64(i.Val)}
}

func toFloat(expr []objects.Object) (objects.Object, error) {
	switch obj := expr[0].(type) {
	case objects.Float:
		return obj, nil
	case objects.Int:
		return intToFloat(obj), nil
	case objects.String:
		switch {
		case parser.IsFloat(obj.Val):
			return parser.ParseFloat(obj.Val)
		case parser.IsInt(obj.Val):
			i, err := parser.ParseInt(obj.Val)
			if err != nil {
				return nil, err
			}
			return intToFloat(i), nil
		default:
			return nil, fmt.Errorf("cannot parse %v to float", obj.Val)
		}
	default:
		return nil, fmt.Errorf("cannot convert object of type %T to float", obj)
	}
}
