package parser

// import (
// 	"reflect"
// 	"testing"
// )

// func TestReadToken(t *testing.T) {
// 	var testCases = []struct {
// 		input    string
// 		expected Token
// 	}{
// 		{"foo", Token{"Symbol", "foo"}},
// 		{"\"Hello World!\"", Token{"String", "Hello World!"}},
// 		{"42", Token{"Int", "42"}},
// 	}

// 	for _, test := range testCases {
// 		lexer := NewLexer(test.input)
// 		token, _ := lexer.readToken()

// 		if token != test.expected {
// 			t.Errorf("Experced: '%s' , got: '%s'", test.expected, token)
// 		}
// 	}
// }

// func TestReadString(t *testing.T) {
// 	input := "\"Hello World!"
// 	lexer := NewLexer(input)
// 	result, err := readString(lexer)

// 	if err == nil {
// 		t.Errorf("Expected an error, got result: '%s'", result)
// 	}
// }

// func TestReadCode(t *testing.T) {
// 	input := "+ 2 97"
// 	expected := []Token{{"Symbol", "+"}, {"Int", "2"}, {"Int", "97"}}
// 	lexer := NewLexer(input)
// 	result, err := lexer.ReadCode()

// 	if err != nil {
// 		t.Errorf("Unexpected error: %s", err)
// 	}
// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("Experced: '%s' , got: '%s'", expected, result)
// 	}
// }
