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

	for {
		entry, err := p.parse()
		if err == nil {
			results = append(results, entry)
		} else {
			if p.lastToken.typ == tokenEOF || p.lastToken.typ == tokenError {
				break
			}
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
	if err != nil {
		return b, err
	}
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

	for p.lastToken.typ != tokenRightBrace && p.lastToken.typ != tokenEOF && p.lastToken.typ != tokenError {
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
		r.addAttribute(key, p.lastToken.Value())
	}

	return r, nil
}

func (p *parser) nextToken() (t lexeme) {
	select {
	case tok, ok := <-p.lexemes:
		if ok {
			t = tok
		} else {
			t = lexeme{tokenError, "Lexer done providing tokens"}
		}
		p.lastToken = t
	}
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
