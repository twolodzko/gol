package parser

import (
	"bufio"
	"io"
)

func isCommentStart(r rune) bool {
	return r == ';'
}

// CodeReader reads runes but ignores non-printable characters, repeated white characters, and code comments
type CodeReader struct {
	reader   *bufio.Reader
	previous rune
}

// NewCodeReader initialize an instance of CodeReader
func NewCodeReader(r io.Reader) *CodeReader {
	return &CodeReader{bufio.NewReader(r), rune(0)}
}

// ReadRune single rune
func (cr *CodeReader) ReadRune() (r rune, err error) {
	isCommented := false

	for {
		r, _, err = cr.reader.ReadRune()

		if err != nil {
			return r, err
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

		cr.previous = r
		break
	}

	return r, err
}

// UnreadRune moves the head of the reader one rune back
func (cr *CodeReader) UnreadRune() error {
	err := cr.reader.UnreadRune()

	if err != nil {
		return err
	}

	return nil
}

// PeekRune reads next rune without moving the head of the reader further
func (cr *CodeReader) PeekRune() (r rune, err error) {
	r, err = cr.ReadRune()
	if err != nil {
		return rune(0), err
	}

	err = cr.reader.UnreadRune()
	if err != nil {
		return rune(0), err
	}

	return r, nil
}
