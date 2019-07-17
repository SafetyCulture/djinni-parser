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

	literal_beg
	// Identifiers and basic type literals
	IDENT  // my_record
	INT    // 12345
	FLOAT  // 123.45
	STRING // "abc"
	literal_end

	operator_beg
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
	operator_end

	keyword_beg
	// Keywords
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

	langflag_beg
	// Language extension flags
	CPP
	OBJC
	JAVA
	langflag_end
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

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

// IsLangFlag returns true for tokens corresponding to language flag;
// it returns false otherwise.
func (tok Token) IsLangFlag() bool { return langflag_beg < tok && tok < langflag_end }
