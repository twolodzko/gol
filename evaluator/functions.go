package evaluator

import (
	"errors"
	"fmt"

	"github.com/twolodzko/goal/objects"
)

type Fn = func([]objects.Object) (objects.Object, error)

type Function interface {
	Call([]objects.Object) (objects.Object, error)
}

type fixedArgsFunction struct {
	body    Fn
	numArgs int
}

func (fn *fixedArgsFunction) Call(args []objects.Object) (objects.Object, error) {
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

func (fn *anyArgsFunction) Call(args []objects.Object) (objects.Object, error) {
	return fn.body(args)
}

func (fn anyArgsFunction) String() string {
	return "<fn>"
}
