package evaluator

import (
	"fmt"

	. "github.com/twolodzko/goal/types"
)

type errNumArgs struct {
	num int
}

func (e *errNumArgs) Error() string {
	return fmt.Sprintf("wrong number of arguments (%d)", e.num)
}

type errWrongType struct {
	obj Any
}

func (e *errWrongType) Error() string {
	return fmt.Sprintf("invalid type %T for %v", e.obj, e.obj)
}
