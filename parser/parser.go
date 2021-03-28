package parser

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/twolodzko/goal/objects"
	"github.com/twolodzko/goal/token"
)

func IsReaderError(err error) bool {
	return err != nil && err != io.EOF
}

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(t []token.Token) *Parser {
	return &Parser{t, 0}
}

func (p *Parser) token() (token.Token, bool) {
	if p.current < len(p.tokens) {
		return p.tokens[p.current], true
	}
	return token.Token{}, false
}

func (p *Parser) nextToken() bool {
	p.current++
	if p.current >= len(p.tokens) {
		p.current--
		return false
	}
	return true
}

func (p *Parser) Parse() ([]objects.Object, error) {
	var (
		parsed []objects.Object
		obj    objects.Object
		err    error
	)

	if len(p.tokens) == 0 {
		return nil, nil
	}

	for {
		t, ok := p.token()

		fmt.Printf("## %v (%v)\n", t, ok)

		if !ok {
			return nil, errors.New("invalid index")
		}

		switch t.Type {
		case token.LPAREN:
			obj = objects.List{}
		case token.INT:
			i, err := strconv.Atoi(t.Literal)

			if err != nil {
				return parsed, err
			}

			obj = objects.Int{Val: i}
		case token.FLOAT:
			f, err := strconv.ParseFloat(t.Literal, 32)

			if err != nil {
				return parsed, err
			}

			obj = objects.Float{Val: f}
		case token.STRING:
			obj = objects.String{Val: t.Literal}
		case token.SYMBOL:
			obj = objects.Symbol{Name: t.Literal}
		}

		parsed = append(parsed, obj)

		if ok := p.nextToken(); !ok {
			break
		}
	}
	return parsed, err
}
