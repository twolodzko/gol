package repl

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/twolodzko/goal/parser"
)

func isBlockStart(r rune) bool {
	return r == '('
}

func isBlockEnd(r rune) bool {
	return r == ')'
}

// Read input from REPL
func Read(in io.Reader) (s string, err error) {
	reader := bufio.NewReader(in)
	openBlocksCount := 0

	for {
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			return "", err
		}

		cr := parser.NewCodeReader(strings.NewReader(line))
		clean := []rune{}

		for {
			r, _, err := cr.ReadRune()

			if err == io.EOF {
				break
			} else if err != nil {
				return "", err
			}

			switch {
			case r == '(':
				openBlocksCount += 1
			case r == ')':
				openBlocksCount -= 1
			}

			clean = append(clean, r)
		}

		s += string(clean)

		if err == io.EOF || openBlocksCount == 0 {
			break
		}
	}

	switch {
	case openBlocksCount < 0:
		return "", errors.New("missing opening bracket")
	case openBlocksCount > 0:
		return "", errors.New("missing closing bracket")
	}

	return s, err
}
