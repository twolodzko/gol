package evaluator

import "fmt"

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
	return fmt.Sprintf("invalid type for %v (%T)", e.obj, e.obj)
}

type ErrNaN struct {
	val Any
}

func (e *ErrNaN) Error() string {
	return fmt.Sprintf("%v (%T) is not a number", e.val, e.val)
}

type ErrNotCallable struct {
	val Any
}

func (e *ErrNotCallable) Error() string {
	return fmt.Sprintf("%v (%T) is not callable", e.val, e.val)
}

type ErrTrace struct {
	callStack []Any
	err       error
}

func (e *ErrTrace) Error() string {
	msg := "\n"
	for i := len(e.callStack) - 1; i >= 0; i-- {
		msg += fmt.Sprintf("%d: %s\n", len(e.callStack)-i, e.callStack[i])
	}
	msg += fmt.Sprintf("\nraised: %s", e.err)
	return msg
}

func Trace(err error, context Any) *ErrTrace {
	switch err := err.(type) {
	case *ErrTrace:
		err.callStack = append(err.callStack, context)
		return err
	default:
		return &ErrTrace{[]Any{context}, err}
	}
}
