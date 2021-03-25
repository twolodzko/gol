package parser

import (
	"errors"
	"io"
	"strconv"
	"unicode"

	"github.com/twolodzko/goal/objects"
	"github.com/twolodzko/goal/reader"
)

func IsListStart(r rune) bool {
	return r == '('
}

func IsListEnd(r rune) bool {
	return r == ')'
}

func isQuotationMark(r rune) bool {
	return r == '"'
}

func isWordBoundary(r rune) bool {
	return unicode.IsSpace(r) || IsListEnd(r) || IsListStart(r)
}

// Parser reads the code and parses it into the AST
type Parser struct {
	*reader.CodeReader
}

// NewParser initializes the Parser
func NewParser(r io.Reader) (*Parser, error) {
	cr, err := reader.NewCodeReader(r)
	return &Parser{cr}, err
}

// Read a quoted string until the closing quotation mark
func (p *Parser) readString() (objects.String, error) {
	var (
		err       error
		str       []rune
		isEscaped bool = false
	)

	if !isQuotationMark(p.Head) {
		return objects.String{}, errors.New("missing opening quotation mark")
	}

	for {
		err = p.NextRune()
		r := p.Head

		if err == io.EOF {
			return objects.String{}, errors.New("missing closing quotation mark")
		}
		if err != nil {
			return objects.String{}, err
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

	return objects.String{Val: string(str)}, err
}

// Read characters until word boundary
func (p *Parser) readWord() (string, error) {
	var (
		err  error
		word []rune
	)

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
func (p *Parser) readList() (objects.List, error) {
	var (
		list []objects.Object
		err  error
	)

	if !IsListStart(p.Head) {
		return objects.List{}, errors.New("missing opening bracket")
	}

	err = p.NextRune()

	if err != nil && err != io.EOF {
		return objects.List{}, err
	}

	list, err = p.Parse()

	if !IsListEnd(p.Head) {
		return objects.List{}, errors.New("missing closing bracket")
	}
	if err != nil {
		return objects.List{}, err
	}

	err = p.NextRune()

	return objects.List{Val: list}, err
}

// readObject reads and parses the single element (atom, symbol, list)
func (p *Parser) readObject() (objects.Object, error) {
	for {
		r := p.Head

		if unicode.IsSpace(r) {
			err := p.NextRune()

			if err != nil {
				return nil, err
			}
			continue
		}

		switch {
		case IsListStart(r):
			// list
			return p.readList()
		case isQuotationMark(r):
			// string
			return p.readString()
		default:
			// number or symbol
			word, err := p.readWord()

			if err != nil && err != io.EOF {
				return objects.List{}, err
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
					return objects.Symbol{Name: word}, err
				}

				return elem, err
			}

			// symbol
			return objects.Symbol{Name: word}, err
		}
	}
}

// Parse the script, return list of expressions to evaluate
func (p *Parser) Parse() ([]objects.Object, error) {
	var (
		expr []objects.Object
		err  error
		obj  objects.Object
	)

	for {
		obj, err = p.readObject()

		if err != nil && err != io.EOF {
			break
		}

		if obj != nil {
			expr = append(expr, obj)
		}

		if err != nil || IsListEnd(p.Head) {
			break
		}
	}

	return expr, err
}

// Try parsing sting to an integer or a float
func stringToNumber(str string) (objects.Object, error) {
	var (
		err error
		f   float64
		i   int
	)

	// try parsing as an int
	i, err = strconv.Atoi(str)

	if err == nil {
		return objects.Int{Val: i}, err
	}

	// try parsing as a float
	f, err = strconv.ParseFloat(str, 64)

	if err == nil {
		return objects.Float{Val: f}, err
	}

	return nil, err
}
