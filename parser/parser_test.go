package parser

import (
	"fmt"
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
