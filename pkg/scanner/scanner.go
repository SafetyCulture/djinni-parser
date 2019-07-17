// Package scanner implements a scanner for Djinni IDL source text.
//
package scanner

import (
	"github.com/SafetyCulture/djinni-parser/pkg/token"
)

// Scanner is a lexical scanner for the Djinni IDL.
type Scanner struct {
	src []byte

	// scanning state
	ch       rune // current character
	offset   int  // character offset
	rdOffset int  // reading offset (position after current character)
}

const bom = 0xFEFF // byte order mark, only permitted as very first character

// Init initiates a Scanner
func (s *Scanner) Init(src []byte) {
	s.src = src
	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0

	s.next()
	if s.ch == bom {
		s.next()
	}
}

// read the next Unicode from the source
// < 0 means end-of-file.
func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		s.ch = rune(s.src[s.rdOffset])
		s.rdOffset += 1
	} else {
		s.offset = len(s.src)
		s.ch = -1 // eof
	}
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

// Scan will scan the next rune and consume any literals
func (s *Scanner) Scan() (tok token.Token, lit string) {
	s.skipWhitespace()

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = token.Lookup(lit)
	case isDigit(ch):
		tok, lit = s.scanNumber()
	default:
		s.next() // always make progress
		switch ch {
		case '@':
			ident := s.scanIdentifier()
			if ident == "import" {
				tok = token.IMPORT
				lit = "@import"
			}
		case '"':
			tok = token.STRING
			lit = s.scanString()
		case '#':
			tok = token.COMMENT
			lit = s.scanComment()
		case '=':
			tok = token.ASSIGN
		case '(':
			tok = token.LPAREN
		case ')':
			tok = token.RPAREN
		case '{':
			tok = token.LBRACE
		case '}':
			tok = token.RBRACE
		case '<':
			tok = token.LANGLE
		case '>':
			tok = token.RANGLE
		case ',':
			tok = token.COMMA
		case ';':
			tok = token.SEMICOLON
		case ':':
			tok = token.COLON
		case '+':
			tok = s.scanLangFlag()
		case -1:
			tok = token.EOF
		default:
			tok = token.ILLEGAL
			lit = string(ch)
		}
	}

	return
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanNumber() (token.Token, string) {
	offs := s.offset
	tok := token.INT
	s.scanMantissa()
	if s.ch == '.' {
		tok = token.FLOAT
		s.next()
		s.scanMantissa()
	}
	return tok, string(s.src[offs:s.offset])
}

func (s *Scanner) scanMantissa() {
	for isDigit(s.ch) {
		s.next()
	}
}

func (s *Scanner) scanString() string {
	offs := s.offset - 1 // '"' opening already consumed
	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			break
		}
		s.next()
		if ch == '"' {
			break
		}
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanComment() string {
	offs := s.offset - 1 // '#' already consumed
	for s.ch != '\n' && s.ch >= 0 {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanLangFlag() (tok token.Token) {
	switch s.ch {
	case 'c':
		tok = token.CPP
		s.next()
	case 'o':
		tok = token.OBJC
		s.next()
	case 'j':
		tok = token.JAVA
		s.next()
	}
	return
}
