package evaluator

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/goal/objects"
	"github.com/twolodzko/goal/parser"
)

func TestEval(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.Object
	}{
		// atoms
		{"()", objects.List{}},
		{"2", objects.Int{Val: 2}},
		{"3.14", objects.Float{Val: 3.14}},
		{`"Hello World!"`, objects.String{Val: "Hello World!"}},
		// functions
		{`(str 3.14)`, objects.String{Val: "3.14"}},
		{`(int "3.14")`, objects.Int{Val: 3}},
		{`(float "1")`, objects.Float{Val: 1}},
	}

	for _, tt := range testCases {
		expr, err := parser.Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := Eval(expr[0])

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}
