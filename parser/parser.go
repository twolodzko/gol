package parser

import (
	"errors"
	"io"
	"strconv"

	"github.com/twolodzko/goal/objects"
	"github.com/twolodzko/goal/token"
)

func Parse(r io.Reader) ([]objects.Object, error) {
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
	if (p.current + 1) < len(p.tokens) {
		p.current++
		return true
	}
	return false
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
		if !ok {
			return nil, errors.New("index out of bonunds")
		}

		switch t.Type {
		case token.RPAREN:
			return parsed, nil
		case token.LPAREN:
			if ok := p.nextToken(); !ok {
				return parsed, errors.New("missing closing brakcet")
			}
			obj, err = p.parseList()
		case token.INT:
			obj, err = parseInt(t)
		case token.FLOAT:
			obj, err = parseFloat(t)
		case token.STRING:
			obj = objects.String{Val: t.Literal}
		case token.SYMBOL:
			obj = objects.Symbol{Name: t.Literal}
		}

		if err != nil {
			return parsed, err
		}

		parsed = append(parsed, obj)

		if ok := p.nextToken(); !ok {
			return parsed, nil
		}
	}
}

func parseInt(t token.Token) (objects.Int, error) {
	i, err := strconv.Atoi(t.Literal)
	return objects.Int{Val: i}, err
}

func parseFloat(t token.Token) (objects.Float, error) {
	f, err := strconv.ParseFloat(t.Literal, 64)
	return objects.Float{Val: f}, err
}

func (p *Parser) parseList() (objects.List, error) {
	l, err := p.Parse()
	return objects.List{Val: l}, err
}
