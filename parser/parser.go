package parser

import (
	"errors"
	"io"
	"strconv"

	"github.com/twolodzko/gol/token"
	. "github.com/twolodzko/gol/types"
)

func Parse(r io.Reader) ([]Any, error) {
	lexer := NewLexer(r)
	tokens, err := lexer.Tokenize()

	if err != nil && err != io.EOF {
		return nil, err
	}

	parser := NewParser(tokens)
	parsed, err := parser.Parse()

	if err != nil {
		return nil, err
	}

	return parsed, nil
}

type Parser struct {
	tokens          []token.Token
	current         int
	openBlocksCount int
}

func NewParser(t []token.Token) *Parser {
	return &Parser{t, 0, 0}
}

func (p *Parser) getToken() (token.Token, bool) {
	if p.current < len(p.tokens) {
		return p.tokens[p.current], true
	}
	return token.Token{}, false
}

func (p *Parser) nextToken() bool {
	if (p.current + 1) < len(p.tokens) {
		p.current++
		return true
	}
	return false
}

func (p *Parser) Parse() ([]Any, error) {
	var (
		parsed []Any
		obj    Any
		err    error
		quoted bool
	)

	if len(p.tokens) == 0 {
		return nil, nil
	}

	for {
		t, ok := p.getToken()
		if !ok {
			return nil, errors.New("index out of bounds")
		}

		if t.Type == token.QUOTE {
			quoted = true
			if ok := p.nextToken(); !ok {
				return parsed, errors.New("missing quoted object")
			}
			continue
		}

		switch t.Type {
		case token.RPAREN:
			p.openBlocksCount--
			if p.openBlocksCount < 0 {
				return parsed, errors.New("missing opening brackets")
			}
			if quoted {
				return parsed, errors.New("missing quoted object")
			}
			return parsed, nil
		case token.LPAREN:
			p.openBlocksCount++
			if ok := p.nextToken(); !ok {
				return parsed, errors.New("missing closing brackets")
			}
			obj, err = p.parseList()
		case token.NIL:
			obj, err = nil, nil
		case token.BOOL:
			obj, err = Bool(t.Literal == "true"), nil
		case token.INT:
			obj, err = ParseInt(t.Literal)
		case token.FLOAT:
			obj, err = ParseFloat(t.Literal)
		case token.STRING:
			obj = String(t.Literal)
		case token.SYMBOL:
			obj = Symbol(t.Literal)
		}

		if err != nil {
			return parsed, err
		}

		if quoted {
			quoted = false
			obj = quote(obj)
		}

		parsed = append(parsed, obj)

		if ok := p.nextToken(); !ok {
			if p.openBlocksCount > 0 {
				return parsed, errors.New("missing closing bracket")
			}
			return parsed, nil
		}
	}
}

func ParseInt(s string) (Int, error) {
	return strconv.Atoi(s)
}

func ParseFloat(s string) (Float, error) {
	return strconv.ParseFloat(s, 64)
}

func (p *Parser) parseList() (List, error) {
	l, err := p.Parse()
	return List(l), err
}

func quote(obj Any) List {
	return List{Symbol("quote"), obj}
}
