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
		{`"Hello \"John\"!"`, objects.String{Val: `Hello "John"!`}},
		{`"It\'s alive!"`, objects.String{Val: "It's alive!"}},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		result, err := parser.readString()

		if IsReaderError(err) {
			t.Errorf("unexpected error: %v", err)
		}
		if result != tt.expected {
			t.Errorf("expected: %v (%T), got: %v (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func Test_readString_InvalidInput(t *testing.T) {
	var testCases = []string{
		`Hello World!"`,
		`"Hello World!`,
		"\"Hello World!\n",
		"\n\"Hello World!",
		"Hello World!\n\t\"",
	}

	for _, input := range testCases {
		parser, err := NewParser(strings.NewReader(input))

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		result, err := parser.readString()

		if err == nil || err == io.EOF {
			t.Errorf("expected an error, got: %v", result)
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
		{"baz", "baz"},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := parser.readWord()

		if IsReaderError(err) {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %v (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func Test_readList(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.List
	}{
		{"()", objects.List{}},
		{"(1 2))", objects.NewList(objects.Int{Val: 1}, objects.Int{Val: 2})},
		{"((1 2))", objects.NewList(objects.NewList(objects.Int{Val: 1}, objects.Int{Val: 2}))},
		{"(1 2) (3 4)", objects.NewList(objects.Int{Val: 1}, objects.Int{Val: 2})},
		{"(a)", objects.NewList(objects.Symbol{Name: "a"})},
		{`(")")`, objects.NewList(objects.String{Val: ")"})},
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

		if IsReaderError(err) {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %v (%T) %q %T", tt.expected, tt.expected, result, result, result.Val[0], result.Val[0])
		}
	}
}

func Test_readList_FailOnMissingBrackets(t *testing.T) {
	var testCases = []string{
		"(1 2",
		"(1 2 ",
		"(1 2\n",
		`(1 2 (3 4)`,
		`((1)`,
		"(1 2 ((3 4)) 5",
		`(1 2 ")"`,
		"1 2)",
		`"(" 1 2)`,
	}

	for _, input := range testCases {
		parser, err := NewParser(strings.NewReader(input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := parser.readList()

		if err == nil || err == io.EOF {
			t.Errorf("for '%s' expected an error, got result: %v (error=%v)", input, result, err)
		}
	}
}

func TestNewParser_EmptyString(t *testing.T) {
	_, err := NewParser(strings.NewReader(""))

	if err != io.EOF {
		t.Errorf("expected EOF error")
	}
}

func Test_readObject(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.Object
	}{
		{"bar ", objects.Symbol{Name: "bar"}},
		{"baz\n", objects.Symbol{Name: "baz"}},
		{"42)", objects.Int{Val: 42}},
		{`"Hello World!"`, objects.String{Val: "Hello World!"}},
		{"42", objects.Int{Val: 42}},
		{"+", objects.Symbol{Name: "+"}},
		{" \n\t bar", objects.Symbol{Name: "bar"}},
		{`(foo 42 "bar")`, objects.NewList(objects.Symbol{Name: "foo"}, objects.Int{Val: 42}, objects.String{Val: "bar"})},
		{"  \n\t(\nfoo \n\n42\t\"bar\")", objects.NewList(objects.Symbol{Name: "foo"}, objects.Int{Val: 42}, objects.String{Val: "bar"})},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := parser.readObject()

		if IsReaderError(err) {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %v (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func Test_readObject_WithEOF(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.Object
	}{
		{"42", objects.Int{Val: 42}},
		{`"Hello World!"`, objects.String{Val: "Hello World!"}},
		{"+", objects.Symbol{Name: "+"}},
		{"  \n\t(\nfoo \n\n42\t\"bar\")", objects.NewList(objects.Symbol{Name: "foo"}, objects.Int{Val: 42}, objects.String{Val: "bar"})},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := parser.readObject()

		if err != io.EOF {
			t.Errorf("expected EOF error, got: %v", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %v (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func Test_readObject_EmptyInput(t *testing.T) {
	var testCases = []struct {
		input    string
		expected objects.Object
	}{
		{" ", nil},
		{"\n\n\t \n", nil},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := parser.readObject()

		if err != io.EOF {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %v (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func TestParse(t *testing.T) {
	var testCases = []struct {
		input    string
		expected []objects.Object
	}{
		{" ", nil},
		{"\n", nil},
		{"bar ", []objects.Object{objects.Symbol{Name: "bar"}}},
		{"foo bar\n", []objects.Object{objects.Symbol{Name: "foo"}, objects.Symbol{Name: "bar"}}},
		{"42)", []objects.Object{objects.Int{Val: 42}}},
		{`"Hello World!" `, []objects.Object{objects.String{Val: "Hello World!"}}},
		{"42", []objects.Object{objects.Int{Val: 42}}},
		{" \n\t bar", []objects.Object{objects.Symbol{Name: "bar"}}},
		{"((1 2))", []objects.Object{objects.NewList(objects.NewList(objects.Int{Val: 1}, objects.Int{Val: 2}))}},
		{`(foo 42 "bar")`, []objects.Object{objects.NewList(objects.Symbol{Name: "foo"}, objects.Int{Val: 42}, objects.String{Val: "bar"})}},
		{"  \n\t(\nfoo \n\n42\t\"bar\")", []objects.Object{objects.NewList(objects.Symbol{Name: "foo"}, objects.Int{Val: 42}, objects.String{Val: "bar"})}},
	}

	for _, tt := range testCases {
		parser, err := NewParser(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err := parser.Parse()

		if IsReaderError(err) {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}
