package evaluator

import (
	"errors"
	"fmt"

	. "github.com/twolodzko/goal/types"
)

type Fn = func(List) (Any, error)

type Function interface {
	Call(List) (Any, error)
}

type fixedArgsFunction struct {
	body    Fn
	numArgs int
}

func (fn *fixedArgsFunction) Call(args List) (Any, error) {
	if fn.numArgs != len(args) {
		return nil, errors.New("wrong number of arguments")
	}
	return fn.body(args[:fn.numArgs])
}

func (fn fixedArgsFunction) String() string {
	return fmt.Sprintf("<fn (%d)>", fn.numArgs)
}

type anyArgsFunction struct {
	body Fn
}

func (fn *anyArgsFunction) Call(args List) (Any, error) {
	return fn.body(args)
}

func (fn anyArgsFunction) String() string {
	return "<fn>"
}
