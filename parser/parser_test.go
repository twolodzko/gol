package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parseString(t *testing.T) {
	var testCases = []struct {
		input    string
		expected interface{}
	}{
		{`"" ignore me`, String{}},
		{`"Hello World!" not this`, String{"Hello World!"}},
		{`"Hello \"John\"!"`, String{"Hello \"John\"!"}},
		{`"It\'s alive!"`, String{"It's alive!"}},
		{`"To escape a char use \\" "ignore me"`, String{"To escape a char use \\"}},
	}

	for _, tt := range testCases {
		reader := NewCodeReader(strings.NewReader(tt.input))
		result, err := parseString(reader)

		if result != tt.expected {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func Test_stringToNumber(t *testing.T) {
	var testCases = []struct {
		input        string
		expected     interface{}
		expectedType string
	}{
		{"3.1415", 3.1415, "float"},
		{"1e-5", 1e-5, "float"},
		{"1.3e+5", 1.3e+5, "float"},
		{"-.34e-5", -0.34e-5, "float"},
		{".223", 0.223, "float"},
		{"+.45", 0.45, "float"},
		{"+42", 42, "int"},
		{"0", 0, "int"},
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
		switch result.(type) {
		case int:
			if tt.expectedType != "int" {
				t.Errorf("unexpected type %T for %s", result, tt.input)
			}
		case float64:
			if tt.expectedType != "float" {
				t.Errorf("unexpected type %T for %s", result, tt.input)
			}
		default:
			t.Errorf("unexpected type %T for %s", result, tt.input)
		}
	}
}

func Test_stringToNumber_InvalidInputs(t *testing.T) {
	var testCases = []struct {
		input string
	}{
		{"0x"},
		{"+a"},
		{"+"},
		{"."},
		{"1234x"},
		{"1e"},
		{"e+5"},
	}

	for _, tt := range testCases {
		result, err := stringToNumber(tt.input)

		if err == nil {
			t.Errorf("for %s expected an error, got %v", tt.input, result)
		}
	}
}

func Test_parseList(t *testing.T) {
	var testCases = []struct {
		input    string
		expected List
	}{
		{"()", List{}},
		{"(1 2) (3 4)", newList(1, 2)},
		{"(1 2 (3 4) 5)", newList(1, 2, newList(3, 4), 5)},
		{"(a)", newList(Symbol{"a"})},
		{"(foo 42 \"Hello World!\" ())", newList(Symbol{"foo"}, 42, String{"Hello World!"}, List{})},
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
		reader := NewCodeReader(strings.NewReader(tt.input))
		result, err := parseList(reader)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}

	reader := NewCodeReader(strings.NewReader(""))
	result, err := parseList(reader)

	if err != io.EOF {
		t.Errorf("unexpected error: %s", err)
	}
	if !cmp.Equal(result, List{}) {
		t.Errorf("unexpected result: %s (%T)", result, result)
	}
}

func Test_parseNode(t *testing.T) {
	var testCases = []struct {
		input    string
		expected interface{}
	}{
		{"42", 42},
		{`"Hello World!"`, String{"Hello World!"}},
		{"+", Symbol{"+"}},
		{"foo", Symbol{"foo"}},
		{"(foo 42 \"bar\")", newList(Symbol{"foo"}, 42, String{"bar"})},
	}

	for _, tt := range testCases {
		reader := NewCodeReader(strings.NewReader(tt.input))
		result, err := parseNode(reader)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}
