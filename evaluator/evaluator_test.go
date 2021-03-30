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
		{"()", objects.List{}},
		{"2", objects.Int{Val: 2}},
		{"3.14", objects.Float{Val: 3.14}},
		{`"Hello World!"`, objects.String{Val: "Hello World!"}},
		// functions
		{`(str 3.14)`, objects.String{Val: "3.14"}},
		{`(int "3.14")`, objects.Int{Val: 3}},
		{`(float "1")`, objects.Float{Val: 1}},
		{`(list "Hello World!" 42 3.14)`, objects.NewList(objects.String{Val: "Hello World!"}, objects.Int{Val: 42}, objects.Float{Val: 3.14})},
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

func TestBooleans(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.Object
	}{
		// booleans: everything is true, unless it's false
		{`(true? "1")`, True},
		{`(true? 0)`, True},
		{`(true? 3.1415)`, True},
		// FIXME
		// {`(true? foo)`, True},
		// {`(true? true)`, True},
		{`(true? ())`, True},
		// {`(true? false)`, False},
		// {`(not true)`, False},
		// {`(not false)`, True},
		{`(not ())`, False},
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
