package repl

import (
	"bufio"
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

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\n'
}

// ParsingError is thrown for invalid input when parsing the code
type ParsingError struct {
	msg string
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("%s", e.msg)
}

// Read input from REPL
func Read(reader *bufio.Reader) (string, error) {
	var (
		err        error = nil
		readerErr  error
		ch         rune
		openBlocks int = 0
	)
	input := []rune{}

	for {
		ch, _, readerErr = reader.ReadRune()

		if readerErr == io.EOF {
			break
		} else {
			if isBlockStart(ch) {
				openBlocks++
			} else if isBlockEnd(ch) {
				openBlocks--
			}

			// allow for breaking lines if
			// any block is still open
			if ch == '\n' && openBlocks == 0 {
				break
			}

			if unicode.IsPrint(ch) {
				input = append(input, ch)
			} else if isWhitespace(ch) {
				input = append(input, ' ')
			} else {
				err = &ParsingError{"Invalid character: " + fmt.Sprintf("%U", ch)}
				break
			}

			if openBlocks < 0 {
				err = &ParsingError{"Missing open bracket"}
				break
			}
		}
	}

	if openBlocks > 0 {
		err = &ParsingError{"Missing closing bracket"}
	}

	return string(input), err
}
