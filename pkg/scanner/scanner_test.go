package scanner_test

import (
	"testing"

	"github.com/SafetyCulture/djinni-parser/pkg/scanner"
	"github.com/SafetyCulture/djinni-parser/pkg/token"
)

type el struct {
	tok token.Token
	lit string
}

var tokens = [...]el{
	{token.COMMENT, "# a comment \n"},

	{token.IDENT, "foobar"},
	{token.IDENT, "_foo"},
	{token.IDENT, "bar1234"},

	{token.INT, "123456"},
	{token.FLOAT, "1234.56"},
	{token.STRING, `"foobar"`},

	{token.ASSIGN, "="},

	{token.LPAREN, "("},
	{token.LBRACE, "{"},
	{token.RPAREN, ")"},
	{token.RBRACE, "}"},
	{token.LANGLE, "<"},
	{token.RANGLE, ">"},

	{token.COMMA, ","},
	{token.SEMICOLON, ";"},
	{token.COLON, ":"},

	{token.ENUM, "enum"},
	{token.FLAGS, "flags"},
	{token.RECORD, "record"},
	{token.INTERFACE, "interface"},

	{token.MAP, "map"},
	{token.SET, "set"},
	{token.LIST, "list"},

	{token.DERIVING, "deriving"},
	{token.EQUALITY, "eq"},
	{token.ORDERING, "ord"},
	{token.PARCELABLE, "parcelable"},

	{token.STATIC, "static"},
	{token.CONST, "const"},
	{token.IMPORT, "@import"},

	{token.CPP, "+c"},
	{token.OBJC, "+o"},
	{token.JAVA, "+j"},
}

const whitespace = "  \t  \n\n\n" // to separate tokens

var source = func() []byte {
	var src []byte
	for _, t := range tokens {
		src = append(src, t.lit...)
		src = append(src, whitespace...)
	}
	return src
}

func TestScan(t *testing.T) {

	var s scanner.Scanner
	s.Init(source())

	for _, e := range tokens {
		tok, lit := s.Scan()

		// check token
		if tok != e.tok {
			t.Errorf("bad token for %q: got %s, expected %s", lit, tok, e.tok)
		}
	}
}
