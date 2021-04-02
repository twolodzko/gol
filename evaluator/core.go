package evaluator

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	. "github.com/twolodzko/goal/types"
)

func Size(obj Any) (Any, error) {
	switch obj := obj.(type) {
	case List:
		return len(obj), nil
	case String:
		return len(obj), nil
	default:
		return nil, nil
	}
}

func Head(obj []Any) (Any, error) {
	if len(obj) != 1 {
		return nil, &ErrNumArgs{len(obj)}
	}

	l, ok := obj[0].(List)
	if !ok {
		return nil, fmt.Errorf("%v is not a list", obj[0])
	}

	return l.Head(), nil
}

func Tail(obj []Any) (Any, error) {
	if len(obj) != 1 {
		return nil, &ErrNumArgs{len(obj)}
	}

	l, ok := obj[0].(List)
	if !ok {
		return nil, fmt.Errorf("%v is not a list", obj[0])
	}

	return l.Tail(), nil
}

func IsNil(obj Any) (Any, error) {
	return Bool(obj == nil), nil
}

func Error(obj []Any) (Any, error) {
	if len(obj) != 1 {
		return nil, &ErrNumArgs{len(obj)}
	}
	msg, ok := obj[0].(String)
	if !ok {
		return nil, &ErrWrongType{obj[0]}
	}

	return nil, fmt.Errorf("%s", msg)
}

func AreSame(obj []Any) (Any, error) {
	if len(obj) != 2 {
		return nil, &ErrNumArgs{len(obj)}
	}
	return Bool(cmp.Equal(obj[0], obj[1])), nil
}

func Print(obj []Any) (Any, error) {
	out := ""
	for _, o := range obj {
		out += fmt.Sprintf("%v", o)
	}
	fmt.Print(out)
	return nil, nil
}

func PrintLn(obj []Any) (Any, error) {
	Print(obj)
	fmt.Println()
	return nil, nil
}

func ToList(exprs []Any) (Any, error) {
	return List(exprs), nil
}
