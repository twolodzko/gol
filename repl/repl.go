package repl

import (
	"bufio"
	"io"

	"github.com/twolodzko/gol/evaluator"
	"github.com/twolodzko/gol/types"
)

type Any = types.Any

type Repl struct {
	reader *bufio.Reader
	*evaluator.Evaluator
}

func NewRepl(in io.Reader) *Repl {
	eval := evaluator.NewEvaluator()
	return &Repl{bufio.NewReader(in), eval}
}

func (repl *Repl) Repl() ([]Any, error) {
	cmd, err := repl.read()
	if err != nil {
		return nil, err
	}
	return repl.EvalString(cmd)
}
