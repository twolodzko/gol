package parser

import "unicode"

func IsListStart(r rune) bool {
	return r == '('
}

func IsListEnd(r rune) bool {
	return r == ')'
}

func IsQuotationMark(r rune) bool {
	return r == '"'
}

func isWordBoundary(r rune) bool {
	return unicode.IsSpace(r) || IsListEnd(r) || IsListStart(r)
}

func IsCommentStart(r rune) bool {
	return r == ';'
}

func isNumberStart(r rune) bool {
	return unicode.IsDigit(r) || r == '-' || r == '+' || r == '.'
}

func isValidRune(r rune) bool {
	return unicode.IsPrint(r) || unicode.IsSpace(r)
}
