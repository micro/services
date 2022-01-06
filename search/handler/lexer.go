package handler

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type itemType int

const (
	itemError itemType = iota
	itemNumber
	itemIdentifier
	itemBoolean
	itemBooleanOp
	itemString
	itemOperator
	itemLeftParen
	itemRightParen
)

const (
	eof = -1
)

type item struct {
	typ itemType
	val string
}

type stateFn func(*lexer) stateFn

func (l *lexer) run() {
	for state := lexStartStatement; state != nil; {
		state = state(l)
	}
	close(l.items)
}

type lexer struct {
	name  string    //used only for error reports
	input string    // the string being scanned
	start int       // start position of this item
	pos   int       // current position in the input
	width int       // width of last rune read
	items chan item // channel of scanned items

}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l, l.items
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func (l *lexer) consumeSpace() {
	l.acceptRun(" ")
	l.ignore()
}

// next returns the next rune in the input
func (l *lexer) next() (rune int32) {
	if l.pos == len(l.input) {
		l.width = 0
		return eof
	}
	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

// ignore skips over the pending input before this point
func (l *lexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune
// Can be called only once per call of next
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input
func (l *lexer) peek() int32 {
	rune := l.next()
	l.backup()
	return rune
}

// accept consumes the next rune
// if it's from the valid set
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRune consumes a run of runes from the valid set
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {

	}
	l.backup()
}

func lexIdent(l *lexer) stateFn {
	l.consumeSpace()
	for {
		r := l.next()
		if r == eof {
			return l.errorf("Unexpected end of input %q", l.input[l.start:])
		}

		if unicode.IsSpace(r) || strings.IndexRune("=><", r) >= 0 {
			l.backup()
			break
		}
	}
	if l.pos > l.start {
		l.emit(itemIdentifier)
	}
	return lexOperator(l)
}

func lexValue(l *lexer) stateFn {
	l.consumeSpace()
	switch {
	case l.accept(`"'`):
		l.backup()
		return lexString(l)
	case strings.HasPrefix(l.input[l.start:], "true"), strings.HasPrefix(l.input[l.start:], "false"):
		return lexBool(l)
	default:
		// try it as a number
		return lexNumber(l)
	}
}

func lexString(l *lexer) stateFn {
	// TODO support single and double quotes with escaping
	if !l.accept(`"'`) {
		return l.errorf("Unexpected value %v, expected a quote", l.peek())
	}
	// ignore the quote
	openQuote := l.input[l.start:l.pos]
	l.ignore()
	lastRead := ""
	for {
		r := l.next()
		if r == eof { // should only happen in error case
			return l.errorf("Unexpected value %v, incorrectly terminated value %s %v", l.input[l.start:], openQuote, lastRead)
		}
		if string(r) == openQuote && lastRead != `\` {
			l.backup()
			l.emit(itemString)
			l.next()
			l.ignore() // ignore the quote
			return lexEndStatement(l)
		}
		lastRead = string(r)

	}
}

const (
	operatorEquals  = `==`
	operatorGreater = `>=`
	operatorLess    = `<=`
	parenLeft       = `(`
	parenRight      = `)`
)

func lexOperator(l *lexer) stateFn {
	l.consumeSpace()
	switch l.input[l.pos : l.pos+2] {
	case operatorEquals, operatorGreater, operatorLess:
		l.pos += 2
		l.emit(itemOperator)
		return lexValue(l)
	}
	// look for identifier
	return l.errorf("Unexpected operator %q", l.input[l.pos:l.pos+2])
}

func lexNumber(l *lexer) stateFn {
	l.consumeSpace()
	// optional leading sign
	l.accept("+-")
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	if l.start == l.pos {
		return l.errorf("Unexpected value %s", l.input[l.start:])
	}
	l.emit(itemNumber)
	return lexEndStatement(l)
}

func lexStartStatement(l *lexer) stateFn {
	l.consumeSpace()
	if string(l.peek()) == parenLeft {
		l.next()
		l.emit(itemLeftParen)
		return lexStartStatement(l)
	}
	return lexIdent(l)
}

func lexEndStatement(l *lexer) stateFn {
	// is this the end of the statement or can we find a boolean op
	l.consumeSpace()
	if string(l.peek()) == parenRight {
		l.next()
		l.emit(itemRightParen)
		return lexEndStatement(l)
	}

	for {
		r := l.next()
		if r == eof {
			break
		}
		if unicode.IsSpace(r) {
			l.backup()
			break
		}
	}
	if l.start == l.pos {
		return nil
	}

	if l.input[l.start:l.pos] == "and" || l.input[l.start:l.pos] == "AND" || l.input[l.start:l.pos] == "or" || l.input[l.start:l.pos] == "OR" {
		l.emit(itemBooleanOp)
		return lexStartStatement(l)
	}
	return l.errorf("Unexpected input %v", l.input[l.start:l.pos])
}

func lexBool(l *lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			break
		}
		if unicode.IsSpace(r) {
			l.backup()
			break
		}
	}
	if l.input[l.start:l.pos] == "true" || l.input[l.start:l.pos] == "false" {
		l.emit(itemBoolean)
		return lexEndStatement(l)
	}

	return l.errorf("Unexpected value %q, expecting a boolean", l.input[l.start:l.pos])
}
