package repl

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/twolodzko/goal/parser"
)

func Read(in io.Reader) (string, error) {
	var (
		err             error
		out             string
		openBlocksCount int
		isQuoted        bool = false
	)

	reader := bufio.NewReader(in)

	for {
	next:
		line, err := reader.ReadString('\n')

		if parser.IsReaderError(err) {
			return line, err
		}

		out += line

		var prev rune

		for _, r := range line {
			if r == ';' {
				prev = '\x00'
				goto next
			}

			switch {
			// handling string
			case !isQuoted && r == '"' && prev != '\\':
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

			prev = r
		}

		if err == io.EOF || openBlocksCount == 0 {
			break
		}
	}

	if openBlocksCount > 0 {
		return out, errors.New("missing closing bracket")
	}

	return out, err
}

func Repl(in io.Reader) (string, error) {
	s, err := Read(in)

	if parser.IsReaderError(err) {
		return "", err
	}

	p, err := parser.NewParser(strings.NewReader(s))

	if err != nil {
		return "", err
	}

	expr, err := p.Parse()

	if parser.IsReaderError(err) {
		return "", err
	}

	out := ""
	for _, obj := range expr {
		out += obj.String()
		out += "\n"
	}

	return out, nil
}
