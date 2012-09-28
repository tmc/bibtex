package bibtex

import (
	"fmt"
)

type BibTeXEntry struct {
	f string
}

const (
	tokenEntryStart tokenType = maxBuiltinToken + iota
	tokenIdentifier
	tokenLeftBrace
	tokenRightBrace
	tokenComma
	tokenEquals
	tokenString
	tokenNumber
)

var tokenTypeLabels = map[tokenType]string{
	tokenEntryStart: "entryStart",
	tokenIdentifier: "identifier",
	tokenLeftBrace:  "",
	tokenRightBrace: "",
	tokenComma:      "",
	tokenEquals:     "",
	tokenNumber:     "number",
	tokenString:     "string",
}

const (
	litEntryStart rune = '@'
	litLeftBrace       = '{'
	litRightBrace      = '}'
	litComma           = ','
	litEquals          = '='
	litQuote           = '"'
)

func (l lexeme) String() string {
	switch l.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return fmt.Sprintf("(error: %s)", l.val)
	}
	typeLabel, ok := tokenTypeLabels[l.typ]
	if !ok {
		typeLabel = "unknown_token"
	}

	// if no label print direct value
	if typeLabel == "" {
		return l.val
	}

	if len(l.val) > 30 {
		return fmt.Sprintf("%s %.30q...", typeLabel, l.val)
	}
	return fmt.Sprintf("%s %q", typeLabel, l.val)
}

func lexBibTeX(input string) (*lexer, chan lexeme) {
	return lex(input, lexTopLevel)
}

// lexing state functions follow

func lexTopLevel(l *lexer) stateFn {
	fmt.Println(" -> TopLevel")
	for {
		r := l.next()
		if r == litEntryStart {
			l.emit(tokenEntryStart)
			return lexEntryType
		}
		if r == eof {
			break
		}
		fmt.Println("tl nomatch", string(r))
	}
	l.emit(tokenEOF)
	return nil
}

func lexEntry(l *lexer) stateFn {
	fmt.Println(" -> Entry")
	for {
		switch r := l.next(); {
		case r == eof:
			fmt.Println("eof", r)
			return l.errorf("unclosed entry")
		case isWhitespace(r):
			l.ignore()
		case isAlphaNumeric(r):
			l.backup()
			return lexEntryType
		case r == litLeftBrace:
			l.emit(tokenLeftBrace)
			return lexEntryBody
		default:
			fmt.Println("No match:", string(r), isAlphaNumeric(r))
			return l.errorf("Unexpected input: %s\n", r)
		}
	}
	return nil
}

func lexEntryType(l *lexer) stateFn {
	fmt.Println(" -> Entry Type")
	for {
		switch r := l.next(); {
		case isWhitespace(r):
			l.ignore()
		case isAlphaNumeric(r):
			// consume
		default:
			l.backup()
			l.emit(tokenIdentifier)
			return lexEntry
		}
	}
	return nil
}

func lexIdentifier(l *lexer) stateFn {
	fmt.Println(" -> Identifier")
	for {
		switch r := l.next(); {
		case isWhitespace(r):
			l.ignore()
		case isAlphaNumeric(r):
			// consume
		default:
			l.backup()
			l.emit(tokenIdentifier)
			return lexEntryBody
		}
	}
	return nil
}

func lexString(l *lexer) stateFn {
	fmt.Println(" -> String")
	for {
		switch r := l.next(); {
		case r == litQuote:
			l.backup()
			l.emit(tokenString)
			l.accept(" ")
			l.ignore()
			return lexEntryBody
		case r == eof:
			return l.errorf("Unexpected EOF")
		}
	}
	return nil
}

func lexEntryBody(l *lexer) stateFn {
	fmt.Println(" -> Entry Body")
	for {
		switch r := l.next(); {
		case isWhitespace(r):
			l.ignore()
		case isAlphaNumeric(r):
			l.backup()
			return lexIdentifier
		case r == litRightBrace:
			l.emit(tokenRightBrace)
			return lexTopLevel
		case r == litComma:
			l.emit(tokenComma)
		case r == litEquals:
			l.emit(tokenEquals)
		case r == litQuote:
			l.ignore()
			return lexString
		default:
			fmt.Println("E unmatch:", string(r))
		}

	}
	return nil
}

/*
func FindBibTeXEntries(s string) []BibTeXEntry {
	results := make([]BibTeXEntry, 0)

	return results
}
*/
