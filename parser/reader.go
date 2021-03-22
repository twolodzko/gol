package parser

import (
	"bufio"
	"unicode"
)

type CodeReader struct {
	reader   *bufio.Reader
	previous rune
}

func newCodeReader(r *bufio.Reader) CodeReader {
	return CodeReader{r, rune(0)}
}

func (r CodeReader) Read() (ch rune, err error) {
	isCommented := false

	for {
		ch, _, err = r.reader.ReadRune()

		if err != nil {
			return ch, err
		}

		// skip all the commented code
		if ch == ';' {
			isCommented = true
			continue
		}
		if isCommented {
			if ch == '\n' {
				isCommented = false
			}
			continue
		}

		if unicode.IsSpace(ch) {
			// unify and strip white characters
			if unicode.IsSpace(r.previous) {
				continue
			}
			ch = ' '
		} else if !unicode.IsPrint(ch) {
			// skip other non-printable chars
			continue
		}

		r.previous = ch
		break
	}

	return ch, err
}

func (r CodeReader) Unread() error {
	return r.reader.UnreadRune()
}
