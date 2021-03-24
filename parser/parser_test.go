package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_stringToNumber(t *testing.T) {
	var testCases = []struct {
		input    string
		expected interface{}
	}{
		{"3.1415", 3.1415},
		{"1e-5", 1e-5},
		{"1.3e+5", 1.3e+5},
		{"-.34e-5", -0.34e-5},
		{".223", 0.223},
		{"+.45", 0.45},
		{"+42", 42},
		{"0", 0},
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
		expected interface{}
	}{
		{`"" ignore me`, String{}},
		{`"Hello World!" not this`, String{"Hello World!"}},
		{`"To escape a char use \\" "ignore me"`, String{"To escape a char use \\"}},
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
		expected interface{}
	}{
		{`"Hello \"John\"!"`, String{`Hello "John"!`}},
		{`"It\'s alive!"`, String{"It's alive!"}},
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
		expected List
	}{
		{"(1 2))", newList(1, 2)},
		{"(1 2) (3 4)", newList(1, 2)},
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
		expected List
	}{
		{"()", List{}},
		{"(a)", newList(Symbol{"a"})},
		{"(1 2 (3 4) 5)", newList(1, 2, newList(3, 4), 5)},
		{`(foo 42 "Hello World!" ())`, newList(Symbol{"foo"}, 42, String{"Hello World!"}, List{})},
		{"(+ -2 +2)", newList(Symbol{"+"}, -2, 2)},
		{"(lambda (x y) (+ x y))", newList(Symbol{"lambda"}, newList(Symbol{"x"}, Symbol{"y"}),
			newList(Symbol{"+"}, Symbol{"x"}, Symbol{"y"}))},
		{`((lambda (x y)
			   (+ x y))
			12 6.17)`,
			newList(newList(Symbol{"lambda"}, newList(Symbol{"x"}, Symbol{"y"}),
				newList(Symbol{"+"}, Symbol{"x"}, Symbol{"y"})), 12, 6.17)},
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
		expected interface{}
	}{
		{"bar ", Symbol{"bar"}},
		{"42)", 42},
		{`"Hello World!" `, String{"Hello World!"}},
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
		expected interface{}
	}{
		{"42", 42},
		{`"Hello World!"`, String{"Hello World!"}},
		{"+", Symbol{"+"}},
		{"foo", Symbol{"foo"}},
		{`(foo 42 "bar")`, newList(Symbol{"foo"}, 42, String{"bar"})},
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
