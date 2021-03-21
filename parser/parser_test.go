package parser

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parseElem(t *testing.T) {
	var testCases = []struct {
		input    string
		expected interface{}
	}{
		{"foo", Symbol{"foo"}},
		{"foo bar", Symbol{"foo"}},
		{"baz)", Symbol{"baz"}},
		{"baz(", Symbol{"baz"}},
		{"\"bar\" baz", String{"bar"}},
		{"2 (baz)", 2},
	}

	for _, tt := range testCases {
		reader := strings.NewReader(tt.input)
		result, err := parseElem(reader)

		if result != tt.expected {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func Test_parseString(t *testing.T) {
	var testCases = []struct {
		input    string
		expected interface{}
	}{
		{"\"\" foo bar", String{}},
		{"\"Hello World!\" foo bar", String{"Hello World!"}},
		//{"\"Hello \\\"Johnny\\\"!\"", String{"Hello \"Johnny\"!"}},
	}

	for _, tt := range testCases {
		reader := strings.NewReader(tt.input)
		result, err := parseString(reader)

		if result != tt.expected {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func Test_parseList(t *testing.T) {
	var testCases = []struct {
		input    string
		expected List
	}{
		{"", List{}},
		{"()", List{}},
		{"(a)", newList("a")},
		{"(foo bar baz)", newList("foo", "bar", "baz")},
		{"(foo (bar baz))", newList("foo", newList("bar", "baz"))},
		{"(foo 42 \"Hello World!\")", newList(Symbol{"foo"}, 42, String{"Hello World!"})},
	}

	for _, tt := range testCases {
		reader := strings.NewReader(tt.input)
		result, err := parseList(reader)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}
