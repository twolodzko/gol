package repl

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/twolodzko/goal/environment"
	"github.com/twolodzko/goal/evaluator"
	"github.com/twolodzko/goal/parser"
	. "github.com/twolodzko/goal/types"
)

type REPL struct {
	reader *bufio.Reader
	env    *environment.Env
}

func NewREPL(in io.Reader, env *environment.Env) *REPL {
	return &REPL{bufio.NewReader(in), env}
}

func (repl *REPL) Repl() ([]Any, error) {
	s, err := repl.Read()
	if err != nil {
		return nil, err
	}

	parsed, err := parser.Parse(strings.NewReader(s))
	if err != nil {
		return nil, err
	}

	evaluated, err := evaluator.EvalAll(parsed, repl.env)
	if err != nil {
		return nil, err
	}

	return evaluated, nil
}

type BlockReader struct {
	*bufio.Reader
	openBlocksCount int
	isQuoted        bool
}

func (repl *REPL) Read() (string, error) {
	var (
		err       error
		out, line string
	)

	reader := BlockReader{repl.reader, 0, false}

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

func (reader *BlockReader) shouldStop(line string) bool {
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
