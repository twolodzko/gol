package evaluator

import (
	"fmt"

	. "github.com/twolodzko/goal/types"
)

type ErrNumArgs struct {
	num int
}

func (e *ErrNumArgs) Error() string {
	return fmt.Sprintf("wrong number of arguments (%d)", e.num)
}

type ErrWrongType struct {
	obj Any
}

func (e *ErrWrongType) Error() string {
	return fmt.Sprintf("invalid type %T for %v", e.obj, e.obj)
}

type ErrNaN struct {
	val Any
}

func (e *ErrNaN) Error() string {
	return fmt.Sprintf("not a number: %v", e.val)
}

type ErrNotAnInteger struct {
	val Any
}

func (e *ErrNotAnInteger) Error() string {
	return fmt.Sprintf("not an integer: %v", e.val)
}
