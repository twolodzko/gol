package parser

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"unicode"

	"github.com/twolodzko/gol/token"
)

var (
	intRegex   = regexp.MustCompile(`^[+-]?\d+$`)
	floatRegex = regexp.MustCompile(`^[+-]?\d*\.?(?:\d+[eE]?[+-]?)?\d+$`)
)

type Lexer struct {
	*CodeReader
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{NewCodeReader(r)}
}

func (l *Lexer) Tokenize() ([]token.Token, error) {
	var tokens []token.Token

	for {
		t, err := l.nextToken()

		if err != nil {
			return tokens, err
		}

		tokens = append(tokens, t)
	}
}

func (l *Lexer) nextToken() (token.Token, error) {
	var (
		str string
		err error
	)

	if err = l.NextRune(); err != nil {
		return token.Token{}, err
	}
	if err := l.skipWhitespace(); err != nil {
		return token.Token{}, err
	}

	r := l.Head

	switch r {
	case '(':
		return token.New(string(r), token.LPAREN), err
	case ')':
		return token.New(string(r), token.RPAREN), err
	case '\'':
		return token.New(string(r), token.QUOTE), err
	case '`':
		return token.New(string(r), token.TICK), err
	case ',':
		return token.New(string(r), token.COMMA), err
	case '"':
		str, err = l.readString()
		return token.New(str, token.STRING), err
	default:
		str, err = l.readWord()
		return token.New(str, guessType(str)), err
	}
}

func IsInt(str string) bool {
	return intRegex.MatchString(str)
}

func IsFloat(str string) bool {
	return floatRegex.MatchString(str)
}

func guessType(str string) string {
	switch {
	case str == "true", str == "false":
		return token.BOOL
	case str == "nil":
		return token.NIL
	case IsInt(str):
		return token.INT
	case IsFloat(str):
		return token.FLOAT
	default:
		return token.SYMBOL
	}
}

func (l *Lexer) skipWhitespace() error {
	for {
		if !unicode.IsSpace(l.Head) {
			return nil
		}
		if err := l.NextRune(); err != nil {
			return err
		}
	}
}

func (l *Lexer) readString() (string, error) {
	var (
		err       error
		runes     []rune
		isEscaped bool = false
	)

	if l.Head != '"' {
		return "", errors.New("missing opening quotation mark")
	}

	for {
		err = l.NextRune()

		if err != nil {
			if err == io.EOF {
				err = errors.New("missing closing quotation mark")
			}
			break
		}

		if !isEscaped {
			// skip the escape sign, unless it was escaped \\
			if l.Head == '\\' {
				isEscaped = true
				continue
			}
			// end of string, unless it was escaped \"
			if l.Head == '"' {
				break
			}
		}

		isEscaped = false

		runes = append(runes, l.Head)
	}

	return string(runes), err
}

func (l *Lexer) readWord() (string, error) {
	var (
		err   error
		runes []rune
	)

	if IsWordBoundary(l.Head) {
		return "", fmt.Errorf("unexpected character: %q", l.Head)
	}

	for {
		runes = append(runes, l.Head)

		err = l.NextRune()

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if IsWordBoundary(l.Head) {
			if err = l.UnreadRune(); err != nil {
				return "", err
			}
			break
		}
	}

	return string(runes), err
}

func IsWordBoundary(r rune) bool {
	return unicode.IsSpace(r) || r == '(' || r == ')'
}
