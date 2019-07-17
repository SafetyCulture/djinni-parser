package parser

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"github.com/SafetyCulture/djinni-parser/pkg/ast"
)

func readSource(filename string, src interface{}) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			return ioutil.ReadAll(s)
		}
		return nil, errors.New("invalid source")
	}
	return ioutil.ReadFile(filename)
}

func ParseFile(filename string, src interface{}) (*ast.IDLFile, error) {
	source, err := readSource(filename, src)
	if err != nil {
		return nil, err
	}

	var p parser
	p.init(source)

	return p.parseFile(), nil
}
