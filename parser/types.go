package parser

import (
	"strconv"
	"unicode"
)

func isBracket(ch rune) bool {
	return ch == '(' || ch == ')'
}

type Symbol struct {
	name string
}

func ParseSymbol(str string, pos int) (Symbol, int, error) {
	var (
		i  int
		ch rune
	)

	// find the length of the sequence
	for i, ch = range str[pos:] {
		if ch == ' ' || isBracket(ch) {
			i--
			break
		}
	}

	return Symbol{str[pos : pos+i+1]}, pos + i, nil
}

func ParseInteger(str string, pos int) (int, int, error) {
	var (
		i  int
		ch rune
	)

	// find the length of the sequence
	for i, ch = range str[pos:] {
		if !(unicode.IsDigit(ch) || (i == 0 && ch == '-')) {
			i--
			break
		}
	}

	out, err := strconv.Atoi(str[pos : pos+i+1])

	if err != nil {
		return 0, pos, err
	}

	return out, pos + i, nil
}
