package repl

import (
	"bufio"
	"io"

	"github.com/twolodzko/goal/evaluator"
	. "github.com/twolodzko/goal/types"
)

type Repl struct {
	reader *bufio.Reader
	*evaluator.Evaluator
}

func NewRepl(in io.Reader) *Repl {
	eval := evaluator.NewEvaluator()
	return &Repl{bufio.NewReader(in), eval}
}

func (repl *Repl) Repl() ([]Any, error) {
	cmd, err := repl.Read()
	if err != nil {
		return nil, err
	}
	return repl.Eval(cmd)
}
