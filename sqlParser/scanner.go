package main

import (
	"bufio"
	"bytes"
	"io"
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r: bufio.NewReader(r),
	}
}

func (s *Scanner) read() rune {
	r, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return r
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) Scan() (Token, string) {

	r := s.read()

	if isWhitespace(r) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(r) {
		s.unread()
		return s.scanIdent()
	}

	switch r {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(r)
	case ',':
		return COMMA, string(r)
	}
	return ILLEGAL, string(r)
}

func (s *Scanner) scanWhitespace() (Token, string) {
	var buf bytes.Buffer

	buf.WriteRune(s.read())

	for {
		r := s.read()
		if r == eof {
			break
		}
		if !isWhitespace(r) {
			s.unread()
			break
		}
		buf.WriteRune(r)
	}
	return WS, buf.String()
}

func (s *Scanner) scanIdent() (Token, string) {
	var buf bytes.Buffer

	buf.WriteRune(s.read())

	for {
		r := s.read()
		if r == eof {
			break
		}
		if !isLetter(r) && !isDigit(r) && r != '_' {
			s.unread()
			break
		}
		buf.WriteRune(r)
	}

	switch string.ToUpper(buf.String()) {
	case "SELECT":
		return SELECT, buf.String()
	case "FROM":
		return FROM, buf.String()
	default:
		return IDENT, buf.String()
	}

}
