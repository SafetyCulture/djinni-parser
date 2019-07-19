package parser

import "fmt"

type parseError struct {
	msg string
}

func (e parseError) Error() string {
	return e.msg
}

type errorsList []parseError

func (e errorsList) add(msg string) {
	e = append(e, parseError{msg})
}

func (e errorsList) Error() string {
	switch len(e) {
	case 0:
		return "no errors"
	case 1:
		return e[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", e[0], len(e)-1)
}
