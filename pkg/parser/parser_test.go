package parser_test

import (
	"testing"

	"github.com/SafetyCulture/djinni-parser/pkg/parser"
)

func TestImports(t *testing.T) {
	src := `
		@import "relative/path/to/filename.djinni"
		@import "relative/path/to/filename2.djinni"
	`

	f, err := parser.ParseFile("", src)
	if err != nil {
		t.Fatal(err)
	}

	if len(f.Imports) != 2 {
		t.Fatalf("incorrect number of imports; expected 2, but got %d", len(f.Imports))
	}

	if f.Imports[0] != "relative/path/to/filename.djinni" {
		t.Errorf("incorrect import path: %q", f.Imports[0])
	}
	if f.Imports[1] != "relative/path/to/filename2.djinni" {
		t.Errorf("incorrect import path: %q", f.Imports[1])
	}
}
