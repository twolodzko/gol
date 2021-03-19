package parser

import (
	"fmt"
)

type Token struct {
	kind  string
	value string
}

func (token *Token) Repr() string {
	return fmt.Sprintf("%s:%s", token.kind, token.value)
}
