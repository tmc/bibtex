package bibtex

import (
	"fmt"
)

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

func (tt tokenType) String() string {
	typeLabel, ok := tokenTypeLabels[tt]
	if !ok {
		typeLabel = "unknown_token"
	}
	return typeLabel
}

func (l lexeme) String() string {
	switch l.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return fmt.Sprintf("(error: %s)", l.val)
	}
	typeLabel := l.typ.String()

	// if no label print direct value
	if typeLabel == "" {
		return l.val
	}

	if len(l.val) > 30 {
		return fmt.Sprintf("%s %.30q...", typeLabel, l.val)
	}
	return fmt.Sprintf("%s %q", typeLabel, l.val)
}

func (l lexeme) Value() string {
	switch l.typ {
	case tokenString:
		return fmt.Sprintf("\"%s\"", l.val)
	default:
		return fmt.Sprint(l.val)
	}
	return "(unknown)"
}

func lexBibTeX(input string) (*lexer, chan lexeme) {
	return lex(input, lexTopLevel)
}

// lexing state functions follow

func lexTopLevel(l *lexer) stateFn {
	for {
		r := l.next()
		if r == litEntryStart {
			l.emit(tokenEntryStart)
			return lexEntryType
		}
		l.ignore() // ignore if we're outside of toplevel
		if r == eof {
			break
		}
	}
	l.emit(tokenEOF)
	return nil
}

func lexEntry(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
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
			l.emit(tokenError)
			return lexTopLevel
		}
	}
	return nil
}

func lexEntryType(l *lexer) stateFn {
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
	for {
		switch r := l.next(); {
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
	for {
		switch r := l.next(); {
		case r == litQuote:
			l.backup()
			l.emit(tokenString)
			l.next()
			l.ignore()
			return lexEntryBody
		case r == eof:
			return l.errorf("Unexpected EOF")
		}
	}
	return nil
}

func lexNumber(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isNumeric(r):
			// consume
		case r == litComma || r == litRightBrace:
			l.backup()
			l.emit(tokenNumber)
			return lexEntryBody
		default:
			l.emit(tokenError)
			return lexTopLevel
		}
	}
	return nil
}

func lexEntryBody(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isWhitespace(r):
			l.ignore()
		case isNumeric(r):
			l.backup()
			return lexNumber
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
			return l.errorf("Unexpected input in entry body: %s\n", string(r))
		}

	}
	return nil
}
