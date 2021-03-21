package parser

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// Read characters until word boundary
func readWord(reader *strings.Reader) (string, error) {
	elem := ""

	for {
		ch, _, err := reader.ReadRune()

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if ch == ' ' || ch == ')' || ch == '(' {
			err := reader.UnreadRune()

			if err != nil {
				return "", err
			}
			break
		}

		elem += string(ch)
	}

	return elem, nil
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

	elem := ""
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

		elem += string(ch)
	}

	return String{elem}, nil
}

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
			if ch == '(' {
				isFirstChar = false
				continue
			} else {
				return List{}, errors.New("missing opening bracket")
			}
		}
		if unicode.IsSpace(ch) {
			continue
		}
		if ch == ')' {
			break
		}

		// ther's no PeekRune, so we used ReadRune instead
		// we go one step back, so we start parsing from the current character
		err := reader.UnreadRune()

		if err != nil {
			return List{}, err
		}

		if ch == '(' {
			// list
			elem, err = parseList(reader)
		} else if ch == '"' {
			// string
			elem, err = parseString(reader)
		} else {
			var word string

			word, err = readWord(reader)

			if err != nil {
				return List{}, err
			}

			switch {
			case unicode.IsDigit(ch) || ch == '-' || ch == '+' || ch == '.':
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
			default:
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
