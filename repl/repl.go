package repl

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/twolodzko/goal/parser"
)

type Reader struct {
	*bufio.Reader
	openBlocksCount int
	isQuoted        bool
}

func NewReader(in io.Reader) *Reader {
	return &Reader{bufio.NewReader(in), 0, false}
}

func (read *Reader) processLine(line string) (string, error) {
	var (
		runes []rune
		err   error
		cr    *parser.CodeReader
	)

	cr, err = parser.NewCodeReader(strings.NewReader(line))

	if err != nil {
		return "", err
	}

	for {
		head := cr.Head

		switch {
		// handling string
		case !read.isQuoted && head == '"':
			read.isQuoted = true
		case read.isQuoted:
			if head == '"' {
				read.isQuoted = false
			}
		// handling list
		case parser.IsListStart(head):
			read.openBlocksCount++
		case parser.IsListEnd(head):
			read.openBlocksCount--

			if read.openBlocksCount < 0 {
				return "", errors.New("missing opening bracket")
			}
		}

		runes = append(runes, head)

		err = cr.NextRune()

		if err == io.EOF {
			return string(runes), nil
		} else if err != nil {
			return string(runes), err
		}
	}
}

func (read *Reader) nextLine() (string, error) {
	line, err := read.ReadString('\n')

	if parser.IsReaderError(err) {
		return "", err
	}

	return read.processLine(line)
}

// Read the input allowing for opened lists to continiue on new lines
func (read *Reader) Read() (string, error) {
	var (
		err error
		out string
	)

	for {
		line, err := read.nextLine()

		if parser.IsReaderError(err) {
			return line, err
		}

		out += line

		if err == io.EOF || read.openBlocksCount == 0 {
			break
		}
	}

	if read.openBlocksCount > 0 {
		return out, errors.New("missing closing bracket")
	}

	return out, err
}

func Repl(in io.Reader) (string, error) {
	reader := NewReader(in)
	s, err := reader.Read()

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
