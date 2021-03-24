package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/goal/objects"
)

func Test_stringToNumber(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.Object
	}{
		{"3.1415", objects.Float{Val: 3.1415}},
		{"1e-5", objects.Float{Val: 1e-5}},
		{"1.3e+5", objects.Float{Val: 1.3e+5}},
		{"-.34e-5", objects.Float{Val: -0.34e-5}},
		{".223", objects.Float{Val: 0.223}},
		{"+.45", objects.Float{Val: 0.45}},
		{"+42", objects.Int{Val: 42}},
		{"0", objects.Int{Val: 0}},
	}

	for _, tt := range testCases {
		result, err := stringToNumber(tt.input)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != tt.expected {
			t.Errorf("expected %v (%T), got %v (%T)",
				tt.expected, tt.expected, result, result,
			)
		}
	}
}

func Test_readString(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.Object
	}{
		{`"" ignore me`, objects.String{}},
		{`"Hello World!" not this`, objects.String{Val: "Hello World!"}},
		{`"To escape a char use \\" "ignore me"`, objects.String{Val: "To escape a char use \\"}},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		result, err := parser.readString()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != tt.expected {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func Test_readString_WithEOF(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.Object
	}{
		{`"Hello \"John\"!"`, objects.String{Val: `Hello "John"!`}},
		{`"It\'s alive!"`, objects.String{Val: "It's alive!"}},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		result, err := parser.readString()

		if err != io.EOF {
			t.Errorf("expected EOF error, got: %v", err)
		}
		if result != tt.expected {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func Test_readWord(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{"foo ", "foo"},
		{"bar)", "bar"},
		{"42)", "42"},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := parser.readWord()

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func Test_readWord_WithEOF(t *testing.T) {
	parser, err := NewParser(strings.NewReader("abc"))

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	result, err := parser.readWord()

	if err != io.EOF {
		t.Errorf("expected EOF error, got: %v", err)
	}
	if result != "abc" {
		t.Errorf("unexpected result: %s", result)
	}
}

func Test_readList(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.List
	}{
		{"(1 2))", objects.NewList(objects.Int{Val: 1}, objects.Int{Val: 2})},
		{"(1 2) (3 4)", objects.NewList(objects.Int{Val: 1}, objects.Int{Val: 2})},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := parser.readList()

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func Test_readList_WithEOF(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.List
	}{
		{"()", objects.List{}},
		{"(a)", objects.NewList(objects.Symbol{Name: "a"})},
		{"(1 2 (3 4) 5)", objects.NewList(objects.Int{Val: 1}, objects.Int{Val: 2}, objects.NewList(objects.Int{Val: 3}, objects.Int{Val: 4}), objects.Int{Val: 5})},
		{`(foo 42 "Hello World!" ())`, objects.NewList(objects.Symbol{Name: "foo"}, objects.Int{Val: 42}, objects.String{Val: "Hello World!"}, objects.List{})},
		{"(+ -2 +2)", objects.NewList(objects.Symbol{Name: "+"}, objects.Int{Val: -2}, objects.Int{Val: 2})},
		{"(lambda (x y) (+ x y))", objects.NewList(objects.Symbol{Name: "lambda"}, objects.NewList(objects.Symbol{Name: "x"}, objects.Symbol{Name: "y"}),
			objects.NewList(objects.Symbol{Name: "+"}, objects.Symbol{Name: "x"}, objects.Symbol{Name: "y"}))},
		{`((lambda (x y)
			   (+ x y))
			12 6.17)`,
			objects.NewList(objects.NewList(objects.Symbol{Name: "lambda"}, objects.NewList(objects.Symbol{Name: "x"}, objects.Symbol{Name: "y"}),
				objects.NewList(objects.Symbol{Name: "+"}, objects.Symbol{Name: "x"}, objects.Symbol{Name: "y"})), objects.Int{Val: 12}, objects.Float{Val: 6.17})},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := parser.readList()

		if err != io.EOF {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func TestReadNext(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.Object
	}{
		{"bar ", objects.Symbol{Name: "bar"}},
		{"42)", objects.Int{Val: 42}},
		{`"Hello World!" `, objects.String{Val: "Hello World!"}},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := parser.ReadNext()

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func TestReadNext_WithEOF(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.Object
	}{
		{"42", objects.Int{Val: 42}},
		{`"Hello World!"`, objects.String{Val: "Hello World!"}},
		{"+", objects.Symbol{Name: "+"}},
		{"foo", objects.Symbol{Name: "foo"}},
		{`(foo 42 "bar")`, objects.NewList(objects.Symbol{Name: "foo"}, objects.Int{Val: 42}, objects.String{Val: "bar"})},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := parser.ReadNext()

		if err != io.EOF {
			t.Errorf("expected EOF error, got: %v", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}
