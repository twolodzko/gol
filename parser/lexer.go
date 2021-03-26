package parser

import "unicode"

func IsListStart(r rune) bool {
	return r == '('
}

func IsListEnd(r rune) bool {
	return r == ')'
}

func isQuotationMark(r rune) bool {
	return r == '"'
}

func isWordBoundary(r rune) bool {
	return unicode.IsSpace(r) || IsListEnd(r) || IsListStart(r)
}

func isCommentStart(r rune) bool {
	return r == ';'
}
