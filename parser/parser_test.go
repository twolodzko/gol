package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/twolodzko/goal/types"
)

func TestParse(t *testing.T) {
	var testCases = []struct {
		input    string
		expected []Any
	}{
		{" ", nil},
		{"\n", nil},
		{"true", []Any{Bool(true)}},
		{"bar ", []Any{Symbol("bar")}},
		{"foo bar\n", []Any{Symbol("foo"), Symbol("bar")}},
		{"42", []Any{Int(42)}},
		{`"Hello World!" `, []Any{String("Hello World!")}},
		{"1e-7", []Any{Float(1e-7)}},
		{" \n\t bar", []Any{Symbol("bar")}},
		{"((1 2))", []Any{List{List{Int(1), Int(2)}}}},
		{`(foo 42 "bar")`, []Any{List{Symbol("foo"), Int(42), String("bar")}}},
		{"  \n\t(\nfoo \n\n42\t\"bar\")", []Any{List{Symbol("foo"), Int(42), String("bar")}}},
		{"(1 2) (3 4)", []Any{List{Int(1), Int(2)}, List{Int(3), Int(4)}}},
		{"(1 2)\n\n(3\n4)", []Any{List{Int(1), Int(2)}, List{Int(3), Int(4)}}},
	}

	for _, tt := range testCases {
		result, err := Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}

	result, err := Parse(strings.NewReader("1 2) (3 4)"))

	if err == nil || err == io.EOF {
		t.Errorf("expected and error, got result: %v (error=%v)", result, err)
	}
}

func TestParse_InvalidInput(t *testing.T) {
	var testCases = []string{
		"(",
		")",
		")(",
		"1 2) (3 4)",
		"(1 2) (3 4",
		"(1 2) 3 4)",
		"(1 2) (",
		"(1 2) )",
	}

	for _, input := range testCases {
		result, err := Parse(strings.NewReader(input))

		if err == nil || err == io.EOF {
			t.Errorf("expected and error, got result: %v (error=%v)", result, err)
		}
	}
}
