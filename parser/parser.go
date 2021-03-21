package parser

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// Read element of a list as a string, without type conversion
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
				elem = ""
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

	return num, nil
}

// Read element of a list and convert type
func parseString(reader *strings.Reader) (String, error) {

	elem := ""
	firstIter := true

	for {
		ch, _, err := reader.ReadRune()

		if err == io.EOF {
			return String{}, errors.New("missing closing quotation sign")
		}
		if err != nil {
			return String{}, err
		}

		if firstIter {
			if ch != '"' {
				return String{}, errors.New("missing opening quotation sign")
			}
			firstIter = false
			continue
		} else if ch == '"' {
			break
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
			var raw string

			raw, err = readWord(reader)

			if raw == "+" || raw == "-" || raw == "." {
				elem = Symbol{raw}
			} else if ch == '.' || ch == '-' || ch == '+' || unicode.IsDigit(ch) {
				// maybe a number
				elem, err = stringToNumber(raw)

				if err != nil {
					// if it starts with a digit, it needs to be a number
					if unicode.IsDigit(ch) {
						return List{}, err
					}
					// symbol
					elem = Symbol{raw}
				}
			} else {
				// symbol
				elem = Symbol{raw}
			}
		}

		if err != nil {
			return List{}, err
		}

		list.Push(elem)
	}

	return list, nil
}
