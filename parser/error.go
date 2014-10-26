package parser

import (
	"strconv"
)

type ParseError struct {
	Line int
	Desc string
}

func (err ParseError) Error() string {
	return err.Desc + " at line " + strconv.Itoa(err.Line)
}
