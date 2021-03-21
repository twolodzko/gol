package parser

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// Read element of a list as a string, without type conversion
func parseRaw(reader *strings.Reader) (string, error) {
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

// Read element of a list and convert type
func parseElem(reader *strings.Reader) (interface{}, error) {
	var elem interface{}

	raw, err := parseRaw(reader)

	if err != nil {
		return nil, err
	}

	runes := []rune(raw)

	if unicode.IsDigit(runes[0]) {
		// number
		elem, err = strconv.Atoi(raw)

		if err != nil {
			return 0, err
		}
	} else if runes[0] == '"' {
		// string
		if runes[len(runes)-1] == '"' {
			elem = String{raw[1 : len(runes)-1]}
		} else {
			return String{}, errors.New("missing closing quotation sign")
		}
	} else {
		// symbol
		elem = Symbol{raw}
	}

	return elem, nil
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
		if ch == ' ' {
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
			// symbol ???
			elem, err = parseRaw(reader)
		}

		if err != nil {
			return List{}, err
		}

		list.Push(elem)
	}

	return list, nil
}
