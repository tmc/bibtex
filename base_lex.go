// Simple lexer essentials
// Based on Rob Pike's "Lexical Scanning in Go" talk.

package bibtex

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	maxBuiltinToken
)

// represents a token (type and string representing value)
type lexeme struct {
	typ tokenType
	val string
}

var eof rune

// holds lexer state
type lexer struct {
	input   string      // the input being scanned
	start   int         // start position of current item
	pos     int         // the current position in the input
	width   int         // width of the last rune read (kept for backup)
	lexemes chan lexeme // channel of scanned lexemes
}

// represents the state of the lexer and returns the next state when run
type stateFn func(*lexer) stateFn

// gives a lexer and a channel of produced lexemes generated from input
func lex(input string, startState stateFn) (*lexer, chan lexeme) {
	l := &lexer{
		input:   input,
		lexemes: make(chan lexeme),
	}
	go l.run(startState)
	return l, l.lexemes
}

// the the lexer until we meet a null state
func (l *lexer) run(startState stateFn) {
	for state := startState; state != nil; {
		state = state(l)
	}
	// signal that no more lexmes will be produced
	close(l.lexemes)
}

// produce a lexeme
func (l *lexer) emit(t tokenType) {
	// slice out the current match
	l.lexemes <- lexeme{t, l.input[l.start:l.pos]}
	// and shift start forward
	l.start = l.pos
}

// return the next rune in the input or eof
func (l *lexer) next() (result rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	result, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return
}

// accept consumes the next rune if it's in the valid string
func (l *lexer) accept(valid string) bool {
	// if rune is in valid list then return true, advancing
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	// otherwise backup and return false
	l.backup()
	return false
}

// acceptRun consumes a run of runes if they are in the valid string
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	// once we encounter the eventually failure, backup
	l.backup()
}

func (l *lexer) acceptRunFunc(validFn func(rune) bool) {
	for validFn(l.next()) {
	}
	// once we encounter the eventual failure, backup
	l.backup()
}

// generate an error token with fmt.Sprintf
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.lexemes <- lexeme{
		tokenError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

// skips over pending input at this point
func (l *lexer) ignore() {
	l.start = l.pos
}

// back up one rune
func (l *lexer) backup() {
	l.pos -= l.width
	l.width = 0
}

// returns the next rune and automatically backs up
func (l *lexer) peek() rune {
	defer l.backup()
	return l.next()
}

// utility functions

func isWhitespace(r rune) bool {
	whitespace := " \t\n\r"
	return strings.IndexRune(whitespace, r) > -1
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isNumeric(r rune) bool {
	return unicode.IsDigit(r)
}
