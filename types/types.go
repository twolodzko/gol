package types

import (
	"fmt"
	"strings"
)

type (
	Any    interface{}
	Symbol string
	Bool   bool
	Int    int
	Float  float64
	String string
	List   []interface{}
)

func (s String) String() string {
	return fmt.Sprintf("%v", string(s))
}

func (l List) String() string {
	var str []string
	for _, elem := range l {
		str = append(str, fmt.Sprintf("%v", elem))
	}
	return "(" + strings.Join(str, " ") + ")"
}

func (l List) Head() Any {
	return l[0]
}

func (l List) Tail() List {
	tail := l[1:]
	if len(tail) > 0 {
		return tail
	}
	return nil
}
