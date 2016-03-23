package main

import (
	"fmt"
	"io"
)

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token
		lit string
		n   int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) scan() (Token, string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	tok, lit = p.s.Scan()

	p.buf.tok, p.buf.lit = tok, lit

	return tok, lit
}

func (p *Parser) unscan() {
	p.buf.n = 1
}

func (p *Parser) scanIgnoreWhitespace() (Token, string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return tok, lit
}

func (p *Parser) Parse() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf(
			"found '%s', expected SELECT", lit,
		)
	}

	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errof(
				"found '%s', expected field", lit,
			)
		}
		stmt.Fields = append(stmt.Fields, lit)
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return nil, fmt.Errorf(
			"found '%s', expected FROM", lit,
		)
	}

	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf(
			"found '%s', expected table name", lit,
		)
	}
	stmt.TableName = lit
	return stmt, nil
}
