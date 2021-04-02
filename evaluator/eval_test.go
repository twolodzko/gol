package evaluator

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/goal/parser"

	. "github.com/twolodzko/goal/types"
)

func TestEvalExpr(t *testing.T) {
	var testCases = []struct {
		input    string
		expected Any
	}{
		// objects
		{`nil`, nil},
		{`()`, List{}},
		{`2`, Int(2)},
		{`3.14`, Float(3.14)},
		{`"Hello World!"`, String("Hello World!")},
		{`true`, Bool(true)},
		// functions
		{`(if true 1 2)`, Int(1)},
		{`(if false 1 2)`, Int(2)},
		{`(if (true? false) (error "this should not fail!") "ok")`, String("ok")},
	}

	for _, tt := range testCases {
		expr, err := parser.Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := EvalExpr(expr[0], BaseEnv)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("for %v expected: %v (%T), got: %s (%T)", tt.input, tt.expected, tt.expected, result, result)
		}
	}
}
