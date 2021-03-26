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
func NewCodeReader(r io.Reader) (*CodeReader, error) {
	cr := &CodeReader{bufio.NewReader(r), rune(0)}
	err := cr.NextRune()
	return cr, err
}

// NextRune moves the head of the reader one rune forward and saves the state in CodeReader.Head
func (cr *CodeReader) NextRune() error {
	isCommented := false

	for {
		r, _, err := cr.ReadRune()

		if err != nil {
			cr.Head = r
			return err
		}
		if !isValidRune(r) {
			cr.Head = r
			return fmt.Errorf("invalid character: %q", r)
		}

		// skip all the commented code
		if IsCommentStart(r) {
			isCommented = true
			continue
		}
		if isCommented {
			if r == '\n' {
				if unicode.IsPrint(cr.Head) && cr.Head != ' ' {
					cr.Head = ' '
					break
				}
				isCommented = false
			}
			continue
		}

		cr.Head = r
		break
	}

	return nil
}
