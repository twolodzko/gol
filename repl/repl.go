package repl

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/twolodzko/goal/evaluator"
	"github.com/twolodzko/goal/parser"
	. "github.com/twolodzko/goal/types"
)

type Reader struct {
	*bufio.Reader
	openBlocksCount int
	isQuoted        bool
}

func Read(in io.Reader) (string, error) {
	var (
		err       error
		out, line string
	)

	reader := Reader{bufio.NewReader(in), 0, false}

	for {
		line, err = reader.ReadString('\n')

		if err != nil && err != io.EOF {
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

	return reader.openBlocksCount <= 0 && !reader.isQuoted
}

func Repl(in io.Reader) ([]Any, error) {
	s, err := Read(in)

	if err != nil {
		return nil, err
	}

	parsed, err := parser.Parse(strings.NewReader(s))

	if err != nil {
		return nil, err
	}

	// FIXME enviroment needs to be fixed per session!
	env := evaluator.InitBuildin()
	evaluated, err := evaluator.EvalAll(parsed, env)

	if err != nil {
		return nil, err
	}

	return evaluated, nil
}
