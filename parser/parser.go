package parser

import (
	"errors"
	"io"
	"strconv"

	"github.com/twolodzko/goal/token"
	. "github.com/twolodzko/goal/types"
)

func Parse(r io.Reader) (List, error) {
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

func (p *Parser) Parse() (List, error) {
	var (
		parsed List
		obj    Any
		err    error
	)

	if len(p.tokens) == 0 {
		return nil, nil
	}

	for {
		t, ok := p.getToken()
		if !ok {
			return nil, errors.New("index out of bounds")
		}

		switch t.Type {
		case token.RPAREN:
			p.openBlocksCount--
			if p.openBlocksCount < 0 {
				return parsed, errors.New("missing opening brackets")
			}
			return parsed, nil
		case token.LPAREN:
			p.openBlocksCount++
			if ok := p.nextToken(); !ok {
				return parsed, errors.New("missing closing brackets")
			}
			obj, err = p.parseList()
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

		parsed = append(parsed, obj)

		if ok := p.nextToken(); !ok {
			if p.openBlocksCount > 0 {
				return parsed, errors.New("missing closing brakcet")
			}
			return parsed, nil
		}
	}
}

func ParseInt(s string) (Int, error) {
	i, err := strconv.Atoi(s)
	return Int(i), err
}

func ParseFloat(s string) (Float, error) {
	f, err := strconv.ParseFloat(s, 64)
	return Float(f), err
}

func (p *Parser) parseList() (List, error) {
	l, err := p.Parse()
	return List(l), err
}
