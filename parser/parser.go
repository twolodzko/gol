package parser

import (
	"errors"
	"io"
	"strconv"
	"unicode"
)

func isListStart(r rune) bool {
	return r == '('
}

func isListEnd(r rune) bool {
	return r == ')'
}

func isQuotationMark(r rune) bool {
	return r == '"'
}

type Parser struct {
	*CodeReader
}

func newParser(reader io.Reader) *Parser {
	return &Parser{NewCodeReader(reader)}
}

// Read a quoted string until the closing quotation mark
func (p *Parser) readString() (String, error) {

	str := []rune{}
	isEscaped := false

	r, err := p.ReadRune()
	if err != nil {
		return String{}, err
	}

	if !isQuotationMark(r) {
		return String{}, errors.New("missing opening quotation mark")
	}

	for {
		r, err = p.ReadRune()

		if err == io.EOF {
			return String{}, errors.New("missing closing quotation mark")
		}
		if err != nil {
			return String{}, err
		}

		if !isEscaped {
			// skip the escape sign, unless it was escaped \\
			if r == '\\' {
				isEscaped = true
				continue
			}
			// end of string, unless it was escaped \"
			if isQuotationMark(r) {
				break
			}
		} else {
			// at next char after the escape, always cancel the escape
			isEscaped = false
		}

		str = append(str, r)
	}

	return String{string(str)}, nil
}

func isWordBoundary(r rune) bool {
	return unicode.IsSpace(r) || isListEnd(r) || isListStart(r)
}

// Read characters until word boundary
func (p *Parser) readWord() (string, error) {
	var (
		r   rune
		err error
	)
	word := []rune{}

	for {
		r, err = p.ReadRune()

		if isWordBoundary(r) {
			err = p.UnreadRune()
			break
		}
		if err != nil {
			break
		}

		word = append(word, r)
	}

	return string(word), err
}

func (p *Parser) readList() (list List, err error) {
	var (
		node interface{}
		r    rune
	)

	r, err = p.ReadRune()

	if err != nil {
		return list, err
	}
	if !isListStart(r) {
		return List{}, errors.New("missing list open bracket")
	}

	for {
		r, err = p.ReadRune()

		if err != nil {
			break
		}

		if unicode.IsSpace(r) {
			continue
		} else if isListEnd(r) {
			break
		}

		err := p.UnreadRune()
		if err != nil {
			break
		}

		node, err = p.readNode()

		list.Push(node)

		if err != nil {
			break
		}
	}

	return list, err
}

func (p *Parser) readNode() (interface{}, error) {
	r, err := p.PeekRune()

	if err != nil {
		return nil, err
	}

	switch {
	case isListStart(r):
		// list
		return p.readList()
	case isQuotationMark(r):
		// string
		return p.readString()
	default:
		// number or symbol
		word, err := p.readWord()

		if err != nil && err != io.EOF {
			return List{}, err
		}

		if unicode.IsDigit(r) || r == '-' || r == '+' || r == '.' {
			// try to parse it as a number
			elem, err := stringToNumber(word)

			if err != nil {
				// if it starts with a digit, it needs to be a number
				if unicode.IsDigit(r) {
					return nil, err
				}

				// otherwise, treat it as a symbol
				return Symbol{word}, nil
			}

			return elem, nil
		}

		// symbol
		return Symbol{word}, nil
	}
}

// Try parsing sting to an integer or a float
func stringToNumber(str string) (num interface{}, err error) {
	num, err = strconv.Atoi(str)

	if err != nil {
		num, err = strconv.ParseFloat(str, 64)
	}

	return num, err
}
