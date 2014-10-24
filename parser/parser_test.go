package parser

import (
	p "github.com/Acconut/poul/program"
	"io"
	"reflect"
	"testing"
)

var programTest = `
# A comment
step1:
foo/bar -> dep/out {
command1
command2
}

 step-2 (pre1, pre2  /post1): 
  foo/*/$1/lol, here.file -> ../hi/ouz, three { 
	echo hello
}


`

func TestParser(t *testing.T) {
	program, err := Parse(string(programTest))
	if err != nil {
		t.Fatal(err)
	}

	expected := &p.Program{
		p.Step{
			Name:      "step1",
			Prehooks:  nil,
			Posthooks: nil,
			Sources: []string{
				"foo/bar",
			},
			Destinations: []string{
				"dep/out",
			},
			Code: `command1
command2
`,
		},
		p.Step{
			Name: "step-2",
			Prehooks: []string{
				"pre1",
				"pre2",
			},
			Posthooks: []string{
				"post1",
			},
			Sources: []string{
				"foo/*/$1/lol",
				"here.file",
			},
			Destinations: []string{
				"../hi/ouz",
				"three",
			},
			Code: `echo hello
`,
		},
	}

	//t.Error(program[0].Prehooks == expected[0].Prehooks)
	if !reflect.DeepEqual(program, expected) {
		t.Errorf("expectation failed: expected\n%s\ngot\n%s\n", expected, program)
	}
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
