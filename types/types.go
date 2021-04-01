package types

import (
	"fmt"
	"strings"
)

type (
	Symbol string
	String string
	List   []interface{}
	Bool   = bool
	Int    = int
	Float  = float64
	Any    = interface{}
)

func (s String) String() string {
	return fmt.Sprintf("\"%s\"", string(s))
}

func (l List) String() string {
	var str []string
	for _, elem := range l {
		str = append(str, fmt.Sprintf("%v", elem))
	}
	return "(" + strings.Join(str, " ") + ")"
}

func (l List) Head() Any {
	if len(l) > 0 {
		return l[0]
	}
	return nil
}

func (l List) Tail() List {
	if len(l) > 1 {
		return l[1:]
	}
	return List{}
}
