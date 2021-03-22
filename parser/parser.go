package parser

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

func isListStart(ch rune) bool {
	return ch == '('
}

func isListEnd(ch rune) bool {
	return ch == ')'
}

// Read characters until word boundary
func readWord(reader *strings.Reader) (string, error) {
	word := []rune{}

	for {
		ch, _, err := reader.ReadRune()

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if unicode.IsSpace(ch) || isListEnd(ch) || isListStart(ch) {
			// we went outside the word boundary, exit
			err := reader.UnreadRune()

			if err != nil {
				return "", err
			}
			break
		}

		word = append(word, ch)
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

// Read a quoted string until the closing quotation sign
func parseString(reader *strings.Reader) (String, error) {

	str := []rune{}
	escaped := false

	// check if it starts with "
	ch, _, err := reader.ReadRune()

	if err != nil {
		return String{}, err
	}
	if ch != '"' {
		return String{}, errors.New("missing opening quotation sign")
	}

	for {
		ch, _, err := reader.ReadRune()

		if err == io.EOF {
			return String{}, errors.New("missing closing quotation sign")
		}
		if err != nil {
			return String{}, err
		}

		if !escaped {
			// skip the escape sign, unless it was escaped \\
			if ch == '\\' {
				escaped = true
				continue
			}
			// end of string, unless it was escaped \"
			if ch == '"' {
				break
			}
		} else {
			// at next char after the escape, always cancel the escape
			escaped = false
		}

		str = append(str, ch)
	}

	return String{string(str)}, nil
}

// Parse a LISP list
func parseList(reader *strings.Reader) (List, error) {
	var (
		elem interface{}
		err  error
		ch   rune
	)
	list := List{}
	isFirstChar := true

	for {
		ch, _, err = reader.ReadRune()

		if err == io.EOF {
			break
		}
		if err != nil {
			return List{}, err
		}

		if isFirstChar {
			if isListStart(ch) {
				isFirstChar = false
				continue
			} else {
				return List{}, errors.New("missing opening bracket")
			}
		}
		if unicode.IsSpace(ch) {
			continue
		}
		if isListEnd(ch) {
			break
		}

		err = reader.UnreadRune()
		if err != nil {
			return List{}, err
		}

		switch {
		case isListStart(ch):
			// list
			elem, err = parseList(reader)
		case ch == '"':
			// string
			elem, err = parseString(reader)
		default:
			// number or symbol
			var word string

			word, err = readWord(reader)

			if err != nil {
				return List{}, err
			}

			if unicode.IsDigit(ch) || ch == '-' || ch == '+' || ch == '.' {
				// try to parse it as a number
				elem, err = stringToNumber(word)

				if err != nil {
					// if it starts with a digit, it needs to be a number
					if unicode.IsDigit(ch) {
						return List{}, err
					}

					// otherwise, treat it as a symbol
					elem = Symbol{word}
					// it was not an error
					err = nil
				}
			} else {
				// symbol
				elem = Symbol{word}
			}
		}

		if err != nil {
			return List{}, err
		}

		list.Push(elem)
	}

	return list, nil
}
