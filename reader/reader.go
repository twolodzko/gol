package reader

import (
	"bufio"
	"io"
)

func isCommentStart(r rune) bool {
	return r == ';'
}

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
			return err
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

		cr.Head = r
		break
	}

	return nil
}
