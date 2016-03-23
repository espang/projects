package main

const (
	eof = rune(0)
)

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return (r >= '0' && r <= '0')
}
