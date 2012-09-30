package bibtex

import (
	"fmt"
	"strings"
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
	tokenLeftBrace:  "{",
	tokenRightBrace: "}",
	tokenComma:      ",",
	tokenEquals:     "=",
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
		return fmt.Sprintf("error: %s", l.val)
	case tokenInvalid:
		return "INVALID"
	}
	typeLabel := l.typ.String()

	// if no label print direct value
	if typeLabel == l.val {
		return typeLabel
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
		case r == eof:
			return l.errorf("unclosed entry")
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
	allNumeric := true
	for {
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unclosed entry")
		case isIdentifierChar(r):
			// consume
			if !isNumeric(r) {
				allNumeric = false
			}
		default:
			l.backup()
			if allNumeric {
				l.emit(tokenNumber)
			} else {
				l.emit(tokenIdentifier)
			}
			return lexEntryBody
		}
	}
	return nil
}

func lexString(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unclosed string")
		case r == litQuote:
			l.backup()
			l.emit(tokenString)
			l.next()
			l.ignore()
			return lexEntryBody
		case r == eof:
			l.backup()
			return lexTopLevel
		}
	}
	return nil
}

func lexEntryBody(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unclosed entry")
		case isIdentifierChar(r):
			l.backup()
			return lexIdentifier
		case isWhitespace(r):
			l.ignore()
		case r == litLeftBrace:
			l.ignore()
			return lexBracedValue
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
			l.emit(tokenError)
			return lexTopLevel
		}

	}
	return nil
}

func lexBracedValue(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unclosed value")
		case r == litRightBrace:
			l.backup()
			l.emit(tokenString)
			l.next()
			l.ignore()
			return lexEntryBody
		}
	}
	return nil
}

func isIdentifierChar(r rune) bool {
	badChars := " \\\t{}\"@,=%#\n\r~"
	return strings.IndexRune(badChars, r) == -1
}
