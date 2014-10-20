package parser

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestParserWeb(t *testing.T) {

	file, err := ioutil.ReadFile("../examples/web")
	if err != nil {
		t.Fatal(err)
	}

	program, err := Parse(string(file))
	fmt.Println(program)
	fmt.Println(err)
}

func TestParserEOF(t *testing.T) {
	program, err := Parse(`
template:
foo -> bar {
`)
	if err != io.EOF {
		t.Error("expected io.EOF")
	}
	if program != nil {
		t.Error("expected nil as return value")
	}
}
