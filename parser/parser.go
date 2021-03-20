package parser

import "github.com/twolodzko/goal/list"

func Parse(str string) (list.List, error) {
	var l list.List

	for _, ch := range str {
		l.Push(ch)
	}

	return l, nil
}
