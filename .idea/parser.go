package old_parser

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	LPAREN = "("
	RPAREN = ")"
)

type ParsingError struct {
	msg string
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("%s", e.msg)
}

func Parse(input string) ([]string, error) {
	input = strings.TrimSpace(input)
	var parsed []string
	word := ""

	if rune(input[0]) != OPEN {
		return nil, &ParsingError{"Missing opening bracket " + string(CLOSE)}
	}

	for _, char := range input {
		if char == OPEN {
			continue
		} else if char == CLOSE {
			parsed = append(parsed, word)
			return parsed, nil
		}

		if unicode.IsSpace(char) == true {
			if word != "" {
				parsed = append(parsed, word)
				word = ""
			}
			continue
		}
		word += string(char)
	}
	parsed = append(parsed, word)

	return parsed, &ParsingError{"Missing closing bracket " + string(CLOSE)}
}
