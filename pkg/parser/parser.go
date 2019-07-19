// Package parser implements a parser for Djinni IDL source files.
//
package parser

import (
	"fmt"

	"github.com/SafetyCulture/djinni-parser/pkg/ast"
	"github.com/SafetyCulture/djinni-parser/pkg/scanner"
	"github.com/SafetyCulture/djinni-parser/pkg/token"
)

type parser struct {
	scanner scanner.Scanner

	tok token.Token // last read token
	lit string      // token literal

	leadComment *ast.CommentGroup // last lead comment

	errors errorsList
}

func (p *parser) init(src []byte) {
	p.scanner.Init(src)
	p.next()
}

func (p *parser) next() {
	p.leadComment = nil

	p.tok, p.lit = p.scanner.Scan()

	if p.tok == token.COMMENT {
		// TODO: Consume the comments to parse as docs for the next token
		p.next()
	}
}

func (p *parser) errorf(msg string, args ...interface{}) {

	// Track all errors and continue parsing.
	p.errors.add(fmt.Sprintf(msg, args...))

	// bailout if too many errors
	if len(p.errors) > 10 {
		// TODO: bailout
	}
}

func (p *parser) expect(tok token.Token) {
	if p.tok != tok {
		p.errorf("expected %q, got %q", tok, p.tok)
	}
	p.next()
}

func (p *parser) parseImport() (i string) {
	p.next()
	if p.tok != token.STRING {
		p.expect(token.STRING)
		return
	}
	// strip the quotes
	i = string(p.lit[1 : len(p.lit)-1])
	p.next()
	return
}

func (p *parser) parseLangExt() ast.Ext {
	ext := ast.Ext{}
	if !p.tok.IsLangExt() {
		p.next()
		return ext
	}
	for p.tok.IsLangExt() {
		switch p.tok {
		case token.CPP:
			ext.CPP = true
		case token.OBJC:
			ext.ObjC = true
		case token.JAVA:
			ext.Java = true
		}
		p.next()
	}
	return ext
}

func (p *parser) parseRecord() *ast.Record {
	p.next()
	ext := p.parseLangExt()
	p.expect(token.LBRACE)

	// TODO: handle all the record fields
	for p.tok != token.RBRACE && p.tok != token.EOF {
		p.next()
	}

	p.expect(token.RBRACE)

	return &ast.Record{
		Ext: ext,
	}
}

func (p *parser) parseInterface() *ast.Interface {
	p.next()
	ext := p.parseLangExt()
	p.expect(token.LBRACE)

	// TODO: handle all the interface methods
	for p.tok != token.RBRACE && p.tok != token.EOF {
		p.next()
	}

	p.expect(token.RBRACE)

	return &ast.Interface{
		Ext: ext,
	}
}

func (p *parser) parseEnum(isFlags bool) *ast.Enum {
	p.expect(token.LBRACE)

	// TODO: handle all options
	for p.tok != token.RBRACE && p.tok != token.EOF {
		p.next()
	}

	p.expect(token.RBRACE)

	return &ast.Enum{
		Flags: isFlags,
	}
}

func (p *parser) parseIdent() ast.Ident {
	name := "_"
	if p.tok == token.IDENT {
		name = p.lit
		p.next()
	} else {
		p.expect(token.IDENT)
	}

	return ast.Ident{Name: name}
}

func (p *parser) parseTypeDef() ast.TypeDef {
	if !p.tok.IsTypeDef() {
		p.errorf("expected one of %v, got %q", token.TypeDefTokens(), p.tok)
		p.next()
	}

	switch p.tok {
	case token.RECORD:
		return p.parseRecord()
	case token.INTERFACE:
		return p.parseInterface()
	case token.ENUM:
		return p.parseEnum(false)
	case token.FLAGS:
		return p.parseEnum(true)
	default:
		return nil
	}
}

// All decls should be in the form IDENT = KEYWORD [EXT] { }
func (p *parser) parseDecl() (decl ast.TypeDecl) {
	decl.Ident = p.parseIdent()
	p.expect(token.ASSIGN)
	decl.Body = p.parseTypeDef()
	return
}

func (p *parser) parseFile() *ast.IDLFile {

	// import decls
	var imports []string
	for p.tok == token.IMPORT {
		imports = append(imports, p.parseImport())
	}

	// rest of body
	var decls []ast.TypeDecl
	for p.tok != token.EOF {
		decls = append(decls, p.parseDecl())
	}

	return &ast.IDLFile{
		Imports:   imports,
		TypeDecls: decls,
	}
}
