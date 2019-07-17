// Package parser implements a parser for Djinni IDL source files.
//
package parser

import (
	"github.com/SafetyCulture/djinni-parser/pkg/ast"
	"github.com/SafetyCulture/djinni-parser/pkg/scanner"
	"github.com/SafetyCulture/djinni-parser/pkg/token"
)

type parser struct {
	scanner scanner.Scanner

	tok token.Token // last read token
	lit string      // token literal

	leadComment *ast.CommentGroup // last lead comment
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
	}

}

func (p *parser) parseImport() (i string) {
	p.next()
	// TODO: Deal with errors if we don't have a STRING
	if p.tok == token.STRING {
		// strip the quotes
		i = string(p.lit[1 : len(p.lit)-1])
	}
	p.next()
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
	// for p.tok != token.EOF {

	// }

	return &ast.IDLFile{
		Imports:   imports,
		TypeDecls: decls,
	}
}
