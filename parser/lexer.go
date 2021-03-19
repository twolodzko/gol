package parser

import (
	"fmt"
	"unicode"
)

func isWordBoundary(ch rune) bool {
	return ch == ' ' || ch == '\n' || ch == ')' || ch == 0
}

// ParsingError is thrown for invalid input when parsing the code
type ParsingError struct {
	msg string
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("%s", e.msg)
}

type Lexer struct {
	runes []rune
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{runes: []rune(input), pos: 0}
}

func readChar(l *Lexer) rune {
	pos := l.pos
	if pos >= len(l.runes) {
		return 0
	}
	// move the head one char forward
	l.pos++
	return l.runes[pos]
}

func readSymbol(l *Lexer) string {
	chars := []rune{}
	for {
		ch := readChar(l)
		if isWordBoundary(ch) {
			break
		}
		chars = append(chars, ch)
	}
	return string(chars)
}

func readString(l *Lexer) (string, error) {
	var ch rune
	chars := []rune{}

	if readChar(l) != '"' {
		return string(chars), &ParsingError{"Missing opening quotes"}
	}

	for {
		ch = readChar(l)
		if ch == '"' || ch == 0 {
			break
		}
		chars = append(chars, ch)
	}

	if ch != '"' {
		return string(chars), &ParsingError{"Missing closing quotes"}
	}

	return string(chars), nil
}

func readInt(l *Lexer) string {
	chars := []rune{}
	for {
		ch := readChar(l)
		if !unicode.IsDigit(ch) {
			break
		}
		chars = append(chars, ch)
	}
	return string(chars)
}

func (l *Lexer) readToken() (token Token, err error) {
	var (
		word string
		kind string
	)
	ch := l.runes[l.pos]

	if ch == '"' {
		word, err = readString(l)
		kind = "String"
	} else if unicode.IsDigit(ch) {
		word = readInt(l)
		kind = "Int"
	} else {
		word = readSymbol(l)
		kind = "Symbol"
	}

	return Token{kind, word}, err
}

func (l *Lexer) ReadCode() (tokens []Token, err error) {
	var token Token

	for {
		if l.pos >= len(l.runes) {
			break
		}
		token, err = l.readToken()
		tokens = append(tokens, token)

		if err != nil {
			break
		}
	}

	return tokens, err
}
