package evaluator

import (
	"errors"
	"fmt"

	"github.com/twolodzko/goal/objects"
	"github.com/twolodzko/goal/parser"
)

type buildin = func(objects.List) (objects.Object, error)

var buildins = map[string]buildin{
	"str":   fixedNumArgs(toString, 1),
	"int":   fixedNumArgs(toInt, 1),
	"float": fixedNumArgs(toFloat, 1),
}

func fixedNumArgs(fn func(objects.Object) (objects.Object, error), numArgs int) func(objects.List) (objects.Object, error) {
	return func(o objects.List) (objects.Object, error) {
		if numArgs == o.Size() {
			return fn(o.Head())
		} else {
			return nil, errors.New("wrong number of arguments")
		}
	}
}

func toString(o objects.Object) (objects.Object, error) {
	switch o := o.(type) {
	case objects.String:
		return o, nil
	default:
		return objects.String{Val: o.String()}, nil
	}
}

func floatToInt(f objects.Float) objects.Int {
	return objects.Int{Val: int(f.Val)}
}

func toInt(o objects.Object) (objects.Object, error) {
	switch o := o.(type) {
	case objects.Int:
		return o, nil
	case objects.Float:
		return floatToInt(o), nil
	case objects.String:
		switch {
		case parser.IsInt(o.Val):
			return parser.ParseInt(o.Val)
		case parser.IsFloat(o.Val):
			f, err := parser.ParseFloat(o.Val)
			if err != nil {
				return nil, err
			}
			return floatToInt(f), nil
		default:
			return nil, fmt.Errorf("cannot parse %v to int", o.Val)
		}
	default:
		return nil, fmt.Errorf("cannot convert object of type %T to int", o)
	}
}

func intToFloat(i objects.Int) objects.Float {
	return objects.Float{Val: float64(i.Val)}
}

func toFloat(o objects.Object) (objects.Object, error) {
	switch o := o.(type) {
	case objects.Float:
		return o, nil
	case objects.Int:
		return intToFloat(o), nil
	case objects.String:
		switch {
		case parser.IsFloat(o.Val):
			return parser.ParseFloat(o.Val)
		case parser.IsInt(o.Val):
			i, err := parser.ParseInt(o.Val)
			if err != nil {
				return nil, err
			}
			return intToFloat(i), nil
		default:
			return nil, fmt.Errorf("cannot parse %v to float", o.Val)
		}
	default:
		return nil, fmt.Errorf("cannot convert object of type %T to float", o)
	}
}