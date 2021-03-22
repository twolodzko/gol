package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
)

func isBlockStart(ch rune) bool {
	return ch == '('
}

func isBlockEnd(ch rune) bool {
	return ch == ')'
}

func isCommentStart(ch rune) bool {
	return ch == ';'
}

// Read input from REPL
func Read(in io.Reader) (string, error) {
	var (
		err             error = nil
		readerErr       error
		ch              rune
		openBlocksCount int  = 0
		isCommented     bool = false
	)
	reader := bufio.NewReader(in)
	input := []rune{}

	for {
		ch, _, readerErr = reader.ReadRune()

		if readerErr == io.EOF {
			break
		} else {
			// after a comment ; ignore everything
			// until end of the line \n
			if isCommentStart(ch) {
				isCommented = true
				continue
			}
			if isCommented {
				if ch == '\n' {
					isCommented = false
				}
				continue
			}

			if isBlockStart(ch) {
				openBlocksCount++
			} else if isBlockEnd(ch) {
				openBlocksCount--
			}

			if openBlocksCount < 0 {
				err = errors.New("missing open bracket")
				break
			}

			// allow for breaking lines
			// if any block is still open
			if ch == '\n' && openBlocksCount == 0 {
				break
			}

			switch {
			case unicode.IsPrint(ch):
				input = append(input, ch)
			case unicode.IsSpace(ch):
				input = append(input, ' ')
			default:
				return "", fmt.Errorf("invalid character: %U", ch)
			}
		}
	}

	if openBlocksCount > 0 {
		err = errors.New("missing closing bracket")
	}

	return string(input), err
}
