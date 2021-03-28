package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/twolodzko/goal/parser"
)

type Reader struct {
	*bufio.Reader
	openBlocksCount int
	isQuoted        bool
}

func (reader *Reader) shouldStop(line string) bool {
	for _, r := range line {

		if r == '\\' {
			continue
		}

		switch {
		case parser.IsQuotationMark(r):
			reader.isQuoted = !reader.isQuoted
		// comment - ignore rest of the line
		case parser.IsCommentStart(r):
			return false
		// list - wait till closing brace
		case parser.IsListStart(r):
			reader.openBlocksCount++
		case parser.IsListEnd(r):
			reader.openBlocksCount--

			if reader.openBlocksCount <= 0 {
				return true
			}
		}
	}

	return reader.openBlocksCount == 0
}

func Read(in io.Reader) (string, error) {
	var (
		err       error
		out, line string
	)

	reader := Reader{bufio.NewReader(in), 0, false}

	for {
		line, err = reader.ReadString('\n')

		if parser.IsReaderError(err) {
			return out, err
		}

		out += line

		if reader.shouldStop(line) || err == io.EOF {
			break
		}
	}

	switch {
	case reader.openBlocksCount > 0:
		err = errors.New("missing closing bracket")
	case reader.openBlocksCount < 0:
		err = errors.New("missing opening bracket")
	}

	return out, err
}

func Repl(in io.Reader) (string, error) {
	s, err := Read(in)

	if parser.IsReaderError(err) {
		return "", err
	}

	l := parser.NewLexer(strings.NewReader(s))

	expr, err := l.Tokenize()

	if parser.IsReaderError(err) {
		return "", err
	}

	out := fmt.Sprintf("%v\n", expr)

	return out, nil
}
