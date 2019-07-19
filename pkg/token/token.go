// Package token defines constants representing the lexical tokens of the Djinni IDL
//
package token

import "strconv"

// Token is the set of lexical tokens of the Djinni IDL
type Token int

// The list of Tokens
const (

	// Special Tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	// Identifiers and basic type literals
	IDENT  // my_record
	INT    // 12345
	FLOAT  // 123.45
	STRING // "abc"

	ASSIGN // =

	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }
	LANGLE // <
	RANGLE // >

	COMMA     // ,
	SEMICOLON // ;
	COLON     // :

	keyword_beg
	// Type Keywords
	ENUM
	FLAGS
	RECORD
	INTERFACE

	MAP
	SET
	LIST

	DERIVING
	EQUALITY
	ORDERING
	PARCELABLE

	STATIC
	CONST
	IMPORT
	keyword_end

	// Language extension flags
	ext_beg
	CPP
	OBJC
	JAVA
	ext_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

	ASSIGN: "=",

	LPAREN: "(",
	RPAREN: ")",
	LBRACE: "{",
	RBRACE: "}",
	LANGLE: "<",
	RANGLE: ">",

	COMMA:     ",",
	SEMICOLON: ";",
	COLON:     ":",

	ENUM:      "enum",
	FLAGS:     "flags",
	RECORD:    "record",
	INTERFACE: "interface",

	MAP:  "map",
	SET:  "set",
	LIST: "list",

	DERIVING:   "deriving",
	EQUALITY:   "eq",
	ORDERING:   "ord",
	PARCELABLE: "parcelable",

	STATIC: "static",
	CONST:  "const",
	IMPORT: "@import",

	CPP:  "+c",
	OBJC: "+o",
	JAVA: "+j",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token CPP, the string is
// "+c"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	// helps with debugging
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

// IsTypeDef returns true for tokens corresponding to type defs;
// it returns false otherwise.
func (tok Token) IsTypeDef() bool { return ENUM <= tok && tok <= INTERFACE }

// TypeDefTokens gets the top level type defs
func TypeDefTokens() []Token {
	var defs []Token
	for i := ENUM; i <= INTERFACE; i++ {
		defs = append(defs, i)
	}
	return defs
}

// IsLangExt returns true for tokens that are language extesions;
// it returns false otherwise.
func (tok Token) IsLangExt() bool { return ext_beg < tok && tok < ext_end }
