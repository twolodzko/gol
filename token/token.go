package token

import "fmt"

const (
	INT    = "int"
	FLOAT  = "float"
	STRING = "str"
	SYMBOL = "sym"
	LPAREN = "("
	RPAREN = ")"
)

type Token struct {
	Literal string
	Type    string
}

func New(l string, t string) Token {
	return Token{Literal: l, Type: t}
}

func (t Token) String() string {
	switch t.Type {
	case LPAREN, RPAREN:
		return t.Type
	case STRING, INT, FLOAT, SYMBOL:
		return fmt.Sprintf("%q:%s", t.Literal, t.Type)
	default:
		return "<ERR>"
	}
}
