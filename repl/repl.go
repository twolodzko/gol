package repl

import (
	"errors"
	"io"
	"unicode"

	"github.com/twolodzko/goal/parser"
)

func isBlockStart(ch rune) bool {
	return ch == '('
}

func isBlockEnd(ch rune) bool {
	return ch == ')'
}

// Read input from REPL
func Read(in io.Reader) (string, error) {
	var (
		err             error
		r               rune
		openBlocksCount int = 0
	)
	reader := parser.NewCodeReader(in)
	input := []rune{}

	for {
		r, _, err = reader.ReadRune()

		if err == io.EOF {
			err = nil
			break
		} else {
			if isBlockStart(r) {
				openBlocksCount++
			} else if isBlockEnd(r) {
				openBlocksCount--
			}

			if openBlocksCount < 0 {
				err = errors.New("missing open bracket")
				break
			}

			// break after the block is closed
			if unicode.IsSpace(r) && openBlocksCount == 0 {
				break
			}

			input = append(input, r)
		}
	}

	if openBlocksCount > 0 {
		err = errors.New("missing closing bracket")
	}

	return string(input), err
}
