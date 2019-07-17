package ast

import (
	"strings"
)

// ----------------------------------------------------------------------------
// Comments

// Comment node represents a single #-style comment.
type Comment struct {
	Text string // comment text excluding '\n'
}

// CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.
type CommentGroup struct {
	List []*Comment
}

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

func stripTrailingWhitespace(s string) string {
	i := len(s)
	for i > 0 && isWhitespace(s[i-1]) {
		i--
	}
	return s[0:i]
}

// Text returns the text of the comment.
// Comment marker (#), the first space of a line comment, and
// leading and trailing empty lines are removed. Multiple empty lines are
// reduced to one, and trailing space on lines is trimmed. Unless the result
// is empty, it is newline-terminated.
func (g *CommentGroup) Text() string {
	if g == nil {
		return ""
	}
	comments := make([]string, len(g.List))
	for i, c := range g.List {
		comments[i] = c.Text
	}

	lines := make([]string, 0, 10) // most comments are less than 10 lines
	for _, c := range comments {
		// Remove comment marker #
		c = c[1:]
		// Strip first space
		if len(c) > 0 && c[0] == ' ' {
			c = c[1:]
		}

		// Split on newlines.
		cl := strings.Split(c, "\n")

		// Walk lines, stripping trailing white space and adding to list.
		for _, l := range cl {
			lines = append(lines, stripTrailingWhitespace(l))
		}
	}

	return strings.Join(lines, "\n")
}

// ----------------------------------------------------------------------------
// Interfaces

// All outer expression nodes implement the TypeDef interface.
type TypeDef interface {
	typeDefNode()
}

// ----------------------------------------------------------------------------
// Expressions

type (
	// Ident node represents an identifier.
	Ident struct {
		Name string
	}

	// Const node represents a constant.
	Const struct {
		Doc   *CommentGroup // associated documentation; or nil
		Ident Ident         // name of the constant
		Type  TypeExpr      // the type of the constant
		Value interface{}   // the value of the constant
	}

	// Ext represents the extension flags that are supported
	Ext struct {
		CPP  bool
		ObjC bool
		Java bool
	}

	// EnumOption represents a single option of an enumeration
	EnumOption struct {
		Doc   *CommentGroup // associated documentation; or nil
		Ident Ident         // name of the option
	}

	TypeExpr struct {
		Ident Ident      // expression type name, eg. i32, i64, string, map, set
		Args  []TypeExpr // arguments to any generic types like map, set and list; or nil
	}

	Field struct {
		Doc   *CommentGroup // associated documentation; or nil
		Ident Ident         // name of the field
		Type  TypeExpr      // the type of the field
	}

	Method struct {
		Doc    *CommentGroup // associated documentation; or nil
		Ident  Ident         // name of the method
		Params []Field
		Return TypeExpr
		Static bool // static method if true
		Const  bool // has been defined as a constant
	}
)

type (
	// Enum node represents an enumeration of options.
	Enum struct {
		Options []EnumOption // options for the enumernation; or nil
		Flags   bool         // true if the enum is defineds as flags
	}

	// Record reperesents a pure-data value object.
	Record struct {
		Ext    Ext // The extensions supported
		Fields []Field
		Consts []Const
	}

	// Interface defines an object with defined methods to call.
	Interface struct {
		Ext     Ext // The extensions supported
		Methods []Method
		Consts  []Const
	}
)

func (*Enum) typeDefNode()      {}
func (*Record) typeDefNode()    {}
func (*Interface) typeDefNode() {}

// ----------------------------------------------------------------------------
// Declarations

// A TypeDecl node represents an enum, flags, record or interface decleration
type TypeDecl struct {
	Doc   *CommentGroup // associated documentation; or nil
	Ident Ident         // name of the identifier
	Body  TypeDef       // decleration type
}

func (*TypeDecl) declNode() {}

// ----------------------------------------------------------------------------
// Files

type IDLFile struct {
	Imports   []string   // imports in this file
	TypeDecls []TypeDecl // top-level declarations; or nil
}
