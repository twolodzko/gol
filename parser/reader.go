package parser

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

// CodeReader reads runes but ignores code comments
type CodeReader struct {
	*bufio.Reader
	Head rune
}

// NewCodeReader initialize an instance of CodeReader
func NewCodeReader(r io.Reader) *CodeReader {
	return &CodeReader{bufio.NewReader(r), rune(0)}
}

// NextRune moves the head of the reader one rune forward and saves the state in CodeReader.Head
func (cr *CodeReader) NextRune() error {
	for {
		r, _, err := cr.ReadRune()

		if err != nil {
			cr.Head = r
			return err
		}
		if !isValidRune(r) {
			return fmt.Errorf("invalid character: %q", r)
		}
		if IsCommentStart(r) {
			err := cr.skipLine()

			if err != nil {
				return err
			}
			if unicode.IsPrint(cr.Head) && cr.Head != ' ' {
				cr.Head = ' '
				return nil
			}

			continue
		}

		cr.Head = r
		return nil
	}
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
