package repl

import (
	"bufio"
	"errors"
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
	var escaped bool

	for _, r := range line {

		switch {
		// string - wait till closing the quote
		case reader.isQuoted && r == '\\':
			escaped = !escaped
		case reader.isQuoted && escaped:
			escaped = false
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

	p, err := parser.NewParser(strings.NewReader(s))

	if parser.IsReaderError(err) {
		return "", err
	}

	expr, err := p.Parse()

	if parser.IsReaderError(err) {
		return "", err
	}

	out := ""
	for _, obj := range expr {
		out += obj.String()
		out += "\n"
	}

	return out, nil
}
