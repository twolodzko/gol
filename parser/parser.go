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

// Read characters until word boundary
func readWord(reader *CodeReader) (string, error) {
	word := []rune{}

	for {
		r, _, err := reader.ReadRune()

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if unicode.IsSpace(r) || isListEnd(r) || isListStart(r) {
			// we went outside the word boundary, exit
			err := reader.UnreadRune()

			if err != nil {
				return "", err
			}
			break
		}

		word = append(word, r)
	}

	return string(word), nil
}

// Try parsing sting to an integer or a float
func stringToNumber(str string) (num interface{}, err error) {
	num, err = strconv.Atoi(str)

	if err != nil {
		num, err = strconv.ParseFloat(str, 64)
	}

	return num, err
}

// Read a quoted string until the closing quotation mark
func parseString(reader *CodeReader) (String, error) {

	str := []rune{}
	isEscaped := false

	// check if it starts with "
	r, _, err := reader.ReadRune()

	if err != nil {
		return String{}, err
	}
	if !isQuotationMark(r) {
		return String{}, errors.New("missing opening quotation mark")
	}

	for {
		r, _, err := reader.ReadRune()

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

func parseNode(reader *CodeReader) (interface{}, error) {

	r, err := reader.PeekRune()

	if err != nil {
		return nil, err
	}

	switch {
	case isListStart(r):
		// list
		return parseList(reader)
	case isQuotationMark(r):
		// string
		return parseString(reader)
	default:
		// number or symbol
		word, err := readWord(reader)

		if err != nil {
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
		} else {
			// symbol
			return Symbol{word}, nil
		}
	}
}

// Parse a LISP list
func parseList(reader *CodeReader) (List, error) {
	var (
		r    rune
		err  error
		node interface{}
	)

	list := List{}

	r, _, err = reader.ReadRune()
	if err != nil {
		return List{}, err
	}
	if !isListStart(r) {
		return List{}, errors.New("missing opening bracket")
	}

	for {
		r, _, err = reader.ReadRune()

		if err == io.EOF {
			break
		}
		if err != nil {
			return List{}, err
		}

		if unicode.IsSpace(r) {
			continue
		} else if isListEnd(r) {
			break
		} else {
			err = reader.UnreadRune()
			if err != nil {
				return List{}, err
			}
		}

		node, err = parseNode(reader)

		list.Push(node)

		if err == io.EOF {
			break
		}
		if err != nil {
			return List{}, err
		}
	}

	return list, err
}
