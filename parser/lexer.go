package parser

import (
	"errors"
	"io"
	"regexp"
	"unicode"

	"github.com/twolodzko/goal/token"
)

var (
	intRegex   = regexp.MustCompile(`^[+-]?\d+$`)
	floatRegex = regexp.MustCompile(`^[+-]?\d*\.?\d+[eE]?[+-]?\d+$`)
)

func IsListStart(r rune) bool {
	return r == '('
}

func IsListEnd(r rune) bool {
	return r == ')'
}

func IsQuotationMark(r rune) bool {
	return r == '"'
}

func isWordBoundary(r rune) bool {
	return unicode.IsSpace(r) || IsListEnd(r) || IsListStart(r)
}

func IsCommentStart(r rune) bool {
	return r == ';'
}

func isNumberStart(r rune) bool {
	return unicode.IsDigit(r) || r == '-' || r == '+' || r == '.'
}

func isValidRune(r rune) bool {
	return unicode.IsPrint(r) || unicode.IsSpace(r)
}

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

		if IsReaderError(err) {
			return tokens, err
		}

		tokens = append(tokens, t)

		if err == io.EOF {
			return tokens, err
		}
	}
}

func (l *Lexer) nextToken() (token.Token, error) {
	var (
		str       string
		err       error
		tokenType string
	)

	if err = l.NextRune(); IsReaderError(err) {
		return token.Token{}, err
	}

	if err := l.skipWhitespace(); err != nil {
		return token.Token{}, err
	}

	r := l.Head

	switch {
	case r == '(':
		return token.New(string(r), token.LPAREN), err
	case r == ')':
		return token.New(string(r), token.RPAREN), err
	case IsQuotationMark(r):
		str, err = l.readString()
		if IsReaderError(err) {
			return token.Token{}, err
		}
		return token.New(str, token.STRING), err
	default:
		str, err = l.readWord()
		if IsReaderError(err) {
			return token.Token{}, err
		}
	}

	switch {
	case intRegex.MatchString(str):
		tokenType = token.INT
	case floatRegex.MatchString(str):
		tokenType = token.FLOAT
	default:
		tokenType = token.SYMBOL
	}

	return token.New(str, tokenType), err
}

func (l *Lexer) skipWhitespace() error {
	for {
		r := l.Head

		if !unicode.IsSpace(r) {
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
		str       []rune
		isEscaped bool = false
	)

	if !IsQuotationMark(l.Head) {
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
			if IsQuotationMark(l.Head) {
				break
			}
		}

		isEscaped = false

		str = append(str, l.Head)
	}

	return string(str), err
}

func (l *Lexer) readWord() (string, error) {
	var (
		err   error
		runes []rune
	)

	if isWordBoundary(l.Head) {
		return "", errors.New("parsing error")
	}

	for {
		runes = append(runes, l.Head)

		err = l.NextRune()

		if err != nil {
			break
		}
		if isWordBoundary(l.Head) {
			if err = l.UnreadRune(); err != nil {
				return "", err
			}
			break
		}
	}

	return string(runes), err
}
