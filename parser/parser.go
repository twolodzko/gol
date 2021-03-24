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

// Parser reads the code and parses it into the AST
type Parser struct {
	*CodeReader
}

// NewParser initializes the Parser
func NewParser(reader io.Reader) (*Parser, error) {
	cr, err := NewCodeReader(reader)
	return &Parser{cr}, err
}

// Read a quoted string until the closing quotation mark
func (p *Parser) readString() (String, error) {
	var err error
	str := []rune{}
	isEscaped := false

	if !isQuotationMark(p.Head) {
		return String{}, errors.New("missing opening quotation mark")
	}

	for {
		err = p.NextRune()
		r := p.Head

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
				err = p.NextRune()
				break
			}
		} else {
			// at next char after the escape, always cancel the escape
			isEscaped = false
		}

		str = append(str, r)
	}

	return String{string(str)}, err
}

func isWordBoundary(r rune) bool {
	return unicode.IsSpace(r) || isListEnd(r) || isListStart(r)
}

// Read characters until word boundary
func (p *Parser) readWord() (string, error) {
	var err error
	word := []rune{}

	for {
		r := p.Head

		if isWordBoundary(r) {
			break
		}

		word = append(word, r)

		err = p.NextRune()

		if err != nil {
			break
		}
	}

	return string(word), err
}

// readList reads the LISP-style list
func (p *Parser) readList() (list List, err error) {
	var node interface{}

	if !isListStart(p.Head) {
		return List{}, errors.New("missing list open bracket")
	}

	err = p.NextRune()

	if err != nil {
		return List{}, err
	}

	for {
		r := p.Head

		// FIXME
		// fmt.Println(string(r))

		if unicode.IsSpace(r) {
			err = p.NextRune()

			if err != nil {
				break
			}

			continue
		} else if isListEnd(r) {
			err = p.NextRune()

			if err != nil {
				break
			}

			break
		}

		node, err = p.ReadNext()

		list.Push(node)

		if err != nil {
			break
		}
	}

	return list, err
}

// ReadNext reads and parses the single element (atom, symbol, list)
func (p *Parser) ReadNext() (interface{}, error) {
	r := p.Head

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
			elem, numberParsingErr := stringToNumber(word)

			if numberParsingErr != nil {
				// if it starts with a digit, it needs to be a number
				if unicode.IsDigit(r) {
					return nil, numberParsingErr
				}

				// otherwise, treat it as a symbol
				return Symbol{word}, err
			}

			return elem, err
		}

		// symbol
		return Symbol{word}, err
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
