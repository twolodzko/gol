package repl

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/twolodzko/goal/parser"
)

// Read the input allowing for opened lists to continiue on new lines
func Read(in io.Reader) (string, error) {
	var (
		err             error
		s               string
		openBlocksCount int
		isQuoted        bool
	)
	lineReader := bufio.NewReader(in)

	for {
		line, err := lineReader.ReadString('\n')

		if err != nil && err != io.EOF {
			return "", err
		}

		cr, err := parser.NewCodeReader(strings.NewReader(line))

		if err != nil {
			return "", err
		}

		runes := []rune{}

		for {
			r := cr.Head

			switch {
			// handling string
			case !isQuoted && r == '"':
				isQuoted = true
			case isQuoted:
				if r == '"' {
					isQuoted = false
				}
			// handling list
			case parser.IsListStart(r):
				openBlocksCount++
			case parser.IsListEnd(r):
				openBlocksCount--

				if openBlocksCount < 0 {
					return "", errors.New("missing opening bracket")
				}
			}

			runes = append(runes, r)

			codeReaderError := cr.NextRune()

			if codeReaderError == io.EOF {
				break
			} else if codeReaderError != nil {
				return "", codeReaderError
			}
		}

		s += string(runes)

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

func Repl(in io.Reader) (string, error) {
	s, err := Read(in)

	if err != nil {
		return "", err
	}

	p, err := parser.NewParser(strings.NewReader(s))

	if err != nil {
		return "", err
	}

	expr, err := p.Parse()

	if err != io.EOF {
		return "", err
	}

	out := ""
	for _, obj := range expr {
		out += obj.String()
		out += "\n"
	}

	return out, nil
}
