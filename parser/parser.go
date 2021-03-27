package parser

import (
	"errors"
	"io"
	"strconv"
	"unicode"

	"github.com/twolodzko/goal/objects"
)

func IsReaderError(err error) bool {
	return err != nil && err != io.EOF
}

// Parser reads the code and parses it into the AST
type Parser struct {
	*CodeReader
}

// NewParser initializes the Parser
func NewParser(r io.Reader) (*Parser, error) {
	cr, err := NewCodeReader(r)
	return &Parser{cr}, err
}

// Parse the script, return list of expressions to evaluate
func (p *Parser) Parse() ([]objects.Object, error) {
	var (
		expr []objects.Object
		err  error
		obj  objects.Object
	)

	for {
		obj, err = p.nextObject()

		if IsReaderError(err) {
			break
		}
		if obj != nil {
			expr = append(expr, obj)
		}
		if err != nil || IsListEnd(p.Head) {
			break
		}
	}

	if len(expr) == 0 {
		return nil, err
	}

	return expr, err
}

// nextObject reads and parses the single element (atom, symbol, list)
func (p *Parser) nextObject() (objects.Object, error) {

	if err := p.skipWhitespace(); err != nil {
		return nil, err
	}

	r := p.Head

	switch {
	case IsListStart(r):
		return p.readList()
	case IsListEnd(r):
		return nil, nil
	case IsQuotationMark(r):
		return p.readString()
	case isNumberStart(r):
		return p.parseNumberOrSymbol()
	default:
		return p.readSymbol()
	}
}

// Read a quoted string until the closing quotation mark
func (p *Parser) readString() (objects.String, error) {
	var (
		err       error
		str       []rune
		isEscaped bool = false
	)

	if !IsQuotationMark(p.Head) {
		return objects.String{}, errors.New("missing opening quotation mark")
	}

	for {
		err = p.NextRune()

		if err != nil {
			if err == io.EOF {
				err = errors.New("missing closing quotation mark")
			}
			break
		}

		if !isEscaped {
			// skip the escape sign, unless it was escaped \\
			if p.Head == '\\' {
				isEscaped = true
				continue
			}
			// end of string, unless it was escaped \"
			if IsQuotationMark(p.Head) {
				err = p.NextRune()
				break
			}
		}

		isEscaped = false

		str = append(str, p.Head)
	}

	return objects.String{Val: string(str)}, err
}

// Read characters until word boundary
func (p *Parser) readWord() (string, error) {
	var (
		err   error
		runes []rune
	)

	for {
		if isWordBoundary(p.Head) {
			break
		}

		runes = append(runes, p.Head)

		err = p.NextRune()

		if err != nil {
			break
		}
	}

	return string(runes), err
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

	if IsReaderError(err) {
		return objects.List{}, err
	}

	list, err = p.Parse()

	switch {
	case IsReaderError(err):
		return objects.List{}, err
	case !IsListEnd(p.Head):
		return objects.List{}, errors.New("missing closing bracket")
	default:
		// if err == io.EOF reading next rune will throw io.EOF again
		err = p.NextRune()
		return objects.List{Val: list}, err
	}
}

func (p *Parser) parseNumberOrSymbol() (objects.Object, error) {

	head := p.Head
	word, err := p.readWord()

	if IsReaderError(err) {
		return nil, err
	}

	elem, numberParsingErr := stringToNumber(word)

	switch {
	case numberParsingErr == nil:
		return elem, err
	// symbols cannot start with a digit
	case unicode.IsDigit(head):
		return nil, errors.New("not a number")
	default:
		return objects.Symbol{Name: word}, err
	}
}

func (p *Parser) readSymbol() (objects.Object, error) {
	word, err := p.readWord()

	if IsReaderError(err) {
		return nil, err
	}

	return objects.Symbol{Name: word}, err
}

func (p *Parser) skipWhitespace() error {
	for {
		r := p.Head

		if !unicode.IsSpace(r) {
			return nil
		}

		if err := p.NextRune(); err != nil {
			return err
		}
	}
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
