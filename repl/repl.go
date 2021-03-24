package repl

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/twolodzko/goal/reader"
)

func isBlockStart(r rune) bool {
	return r == '('
}

func isBlockEnd(r rune) bool {
	return r == ')'
}

// Read input from REPL
func Read(in io.Reader) (string, error) {
	var err error
	lineReader := bufio.NewReader(in)
	openBlocksCount := 0
	s := ""

	for {
		line, err := lineReader.ReadString('\n')

		if err != nil && err != io.EOF {
			return "", err
		}

		cr, err := reader.NewCodeReader(strings.NewReader(line))

		if err != nil {
			return "", err
		}

		clean := []rune{}

		for {
			r := cr.Head

			switch {
			case r == '(':
				openBlocksCount++
			case r == ')':
				openBlocksCount--
			}

			if openBlocksCount < 0 {
				return "", errors.New("missing opening bracket")
			}

			clean = append(clean, r)

			err = cr.NextRune()

			if err == io.EOF {
				break
			} else if err != nil {
				return "", err
			}
		}

		s += string(clean)

		if err == io.EOF || openBlocksCount == 0 {
			break
		}
	}

	if openBlocksCount > 0 {
		return "", errors.New("missing closing bracket")
	}
	if err != nil && err != io.EOF {
		return "", err
	}

	return s, err
}
