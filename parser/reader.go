package parser

import (
	"bufio"
	"io"
	"unicode"
)

// CodeReader reads runes but ignores non-printable characters, repeated white characters, and code comments
type CodeReader struct {
	reader   *bufio.Reader
	previous rune
}

// NewCodeReader initialize an instance of CodeReader
func NewCodeReader(r io.Reader) CodeReader {
	return CodeReader{bufio.NewReader(r), rune(0)}
}

func isCommentStart(ch rune) bool {
	return ch == ';'
}

// ReadRune single rune
func (cr *CodeReader) ReadRune() (r rune, size int, err error) {
	isCommented := false

	for {
		r, size, err = cr.reader.ReadRune()

		if err != nil {
			return r, size, err
		}

		// skip all the commented code
		if isCommentStart(r) {
			isCommented = true
			continue
		}
		if isCommented {
			if r == '\n' {
				isCommented = false
			}
			continue
		}

		if unicode.IsSpace(r) {
			// unify and strip white characters
			if unicode.IsSpace(cr.previous) {
				continue
			}
			r = ' '
		} else if !unicode.IsPrint(r) {
			// skip other non-printable chars
			continue
		}

		cr.previous = r
		break
	}

	return r, size, err
}

// UnreadRune the last read rune
func (cr *CodeReader) UnreadRune() error {
	return cr.reader.UnreadRune()
}
