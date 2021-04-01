package parser

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type CodeReader struct {
	*bufio.Reader
	Head rune
}

func NewCodeReader(r io.Reader) *CodeReader {
	return &CodeReader{bufio.NewReader(r), rune(0)}
}

func (cr *CodeReader) NextRune() error {
	r, _, err := cr.ReadRune()

	switch {
	case err != nil:
		cr.Head = rune(0)
		return err
	case !isValidRune(r):
		cr.Head = rune(0)
		return fmt.Errorf("invalid character: %q", r)
	case IsCommentStart(r):
		err := cr.skipLine()
		cr.Head = ' '
		return err
	}

	cr.Head = r
	return nil
}

func (cr *CodeReader) skipLine() error {
	for {
		r, _, err := cr.ReadRune()

		if err != nil {
			return err
		}
		if r == '\n' {
			return nil
		}
	}
}

func isValidRune(r rune) bool {
	return unicode.IsPrint(r) || unicode.IsSpace(r)
}

func IsCommentStart(r rune) bool {
	return r == ';'
}
