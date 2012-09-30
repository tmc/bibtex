package bibtex

import (
	"errors"
	"fmt"
)

type parser struct {
	lexer     *lexer
	lexemes   chan lexeme
	lastToken lexeme
}

func ParseBibTeX(in string) (BibTeXEntry, error) {
	return parseBibTeX(in)
}

func parseBibTeX(in string) (BibTeXEntry, error) {

	lexer, lexemes := lexBibTeX(in)
	p := &parser{
		lexer:   lexer,
		lexemes: lexemes,
	}

	return p.parse()
}

func (p *parser) parse() (b BibTeXEntry, err error) {
	b = newBibTeXEntry()
	err = p.expect(tokenEntryStart)

	if err != nil {
		return b, err
	}
	b, err = p.entry_type(b)

	err = p.expect(tokenLeftBrace)
	if err != nil {
		return b, err
	}

	b, err = p.entry_body(b)
	if err != nil {
		return b, err
	}

	return b, err
}

func (p *parser) entry_type(b BibTeXEntry) (r BibTeXEntry, err error) {
	r = b
	err = p.expect(tokenIdentifier)
	r.Type = p.lastToken.val
	return r, nil
}

func (p *parser) entry_body(b BibTeXEntry) (r BibTeXEntry, err error) {
	r = b
	err = p.expect(tokenIdentifier)
	if err != nil {
		return r, err
	}

	r.Identifier = p.lastToken.val

	for p.lastToken.typ != tokenRightBrace && p.lastToken.typ != tokenEOF {
		p.accept(tokenComma)

		t := p.nextToken()
		if t.typ == tokenRightBrace || t.typ == tokenEOF {
			break
		}

		if t.typ != tokenIdentifier {
			err = errors.New(fmt.Sprintf("error: unexpected token: %s (expected %s)", p.lastToken, tokenIdentifier))
			return r, err
		}
		key := p.lastToken.val

		err = p.expect(tokenEquals)

		t = p.nextToken()
		if t.typ != tokenString && t.typ != tokenNumber {
			err = errors.New(fmt.Sprintf("error: unexpected token: %s (expected string or number)", p.lastToken))
			return r, err
		}

		if err != nil {
			return r, err
		}
		r.addAttribute(key, p.lastToken.Value())
	}

	return r, nil
}

func (p *parser) nextToken() (t lexeme) {
	select {
	case t = <-p.lexemes:
	}
	p.lastToken = t
	return
}

func (p *parser) accept(l tokenType) bool {
	if t := p.nextToken(); t.typ == l {
		return true
	}
	return false
}

func (p *parser) expect(l tokenType) error {
	if !p.accept(l) {
		return errors.New(fmt.Sprintf("error: unexpected token: %s (expected %s)", p.lastToken, l))
	}
	return nil
}
