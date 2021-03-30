package evaluator

import (
	"fmt"

	"github.com/twolodzko/goal/objects"
)

func Eval(expr objects.Object) (objects.Object, error) {
	switch expr := expr.(type) {
	case objects.Int, objects.Float, objects.String:
		return expr, nil
	case objects.List:
		return evalList(expr)
	default:
		return nil, fmt.Errorf("cannot evaluate object of type %T", expr)
	}
}

func evalList(expr objects.List) (objects.Object, error) {
	if expr.Size() > 0 {
		switch name := expr.Head().(type) {
		case objects.Symbol:
			args := expr.Tail()
			return buildins[name.Val](args)
		default:
			return nil, fmt.Errorf("cannot evaluate list: %v", expr)
		}
	}
	return expr, nil
}
