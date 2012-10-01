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

// @todo move to Reader

func ParseBibTeXEntry(in string) (BibTeXEntry, error) {
	return parseBibTeXSingle(in)
}

func ParseBibTeXEntries(in string) []BibTeXEntry {
	return parseBibTeXMultiple(in)
}

func parseBibTeXSingle(in string) (BibTeXEntry, error) {
	lexer, lexemes := lexBibTeX(in)
	p := &parser{
		lexer:   lexer,
		lexemes: lexemes,
	}
	return p.parse()
}

func parseBibTeXMultiple(in string) (results []BibTeXEntry) {
	results = make([]BibTeXEntry, 0)

	lexer, lexemes := lexBibTeX(in)
	p := &parser{
		lexer:   lexer,
		lexemes: lexemes,
	}

	for p.lastToken.typ != tokenEOF {
		entry, err := p.parse()
		if err == nil {
			results = append(results, entry)
		}

	}
	return results
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

	return b, err
}

func (p *parser) entry_type(b BibTeXEntry) (r BibTeXEntry, err error) {
	r = b
	err = p.expect(tokenIdentifier)
	r.Type = p.lastToken.val
	return r, err
}

func (p *parser) entry_body(b BibTeXEntry) (r BibTeXEntry, err error) {
	r = b
	err = p.expect(tokenIdentifier)
	if err != nil {
		return r, err
	}
	r.Identifier = p.lastToken.val
	p.accept(tokenComma)

	return p.value_list(r)
}

func (p *parser) value_list(b BibTeXEntry) (r BibTeXEntry, err error) {
	r = b

	t := p.lastToken
	if t.typ == tokenRightBrace {
		return r, nil
	}
	if t.typ == tokenError || t.typ == tokenEOF {
		return r, errors.New(fmt.Sprint(t))
	}

	ok := p.accept(tokenIdentifier)
	if ok {
		t := p.lastToken
		if t.typ != tokenIdentifier {
			err = errors.New(fmt.Sprintf("error: unexpected token: %s (expected identifier)", p.lastToken))
			return r, err
		}
		key := t.val

		err = p.expect(tokenEquals)
		if err != nil {
			return r, err
		}

		t = p.nextToken()
		if t.typ != tokenString && t.typ != tokenNumber && t.typ != tokenIdentifier {
			err = errors.New(fmt.Sprintf("error: unexpected token: %s (expected string or number)", p.lastToken))
			return r, err
		}

		if err != nil {
			return r, err
		}
		if t.typ == tokenNumber {
			r.addAttribute(key, p.lastToken.val, Numeric)
		} else {
			r.addAttribute(key, p.lastToken.val, String)
		}

		p.accept(tokenComma)
	}

	return p.value_list(r)
}

func (p *parser) nextToken() lexeme {
	tok, ok := <-p.lexemes
	if !ok {
		tok = lexeme{tokenEOF, "Lexer done providing tokens"}
	}
	p.lastToken = tok
	return tok
}

func (p *parser) accept(l tokenType) bool {
	if p.nextToken().typ == l {
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
