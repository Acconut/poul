package parser

import (
	p "github.com/Acconut/poul/program"
	"io"
	"reflect"
	"testing"
)

var programTest = `
# A comment


 template-empty (pre1, pre2  /post1) { 
}


foo/bar -> dep/out {
command1
command2
}

template-1 {
dist/foo.html
bar/lol.hi
}

  foo/*/$1/lol (  here.file, lol/hoo) -> ../hi/ouz { 
	echo hello
}


`

func TestParser(t *testing.T) {
	program, err := Parse(string(programTest))
	if err != nil {
		t.Fatal(err)
	}

	expected := &p.Program{
		Templates: map[string]p.Template{
			"template-1": p.Template{
				Name:      "template-1",
				Prehooks:  nil,
				Posthooks: nil,
				Destinations: []string{
					"dist/foo.html",
					"bar/lol.hi",
				},
			},
			"template-empty": p.Template{
				Name: "template-empty",
				Prehooks: []string{
					"pre1",
					"pre2",
				},
				Posthooks: []string{
					"post1",
				},
				Destinations: []string{
					"",
				},
			},
		},
		Steps: []p.Step{
			p.Step{
				Source:      "foo/bar",
				Destination: "dep/out",
				Code: `command1
command2
`,
			},
			p.Step{
				Source:      "foo/*/$1/lol",
				Destination: "../hi/ouz",
				Dependencies: []string{
					"here.file",
					"lol/hoo",
				},
				Code: `echo hello
`,
			},
		},
	}

	//t.Error(program[0].Prehooks == expected[0].Prehooks)
	if !reflect.DeepEqual(program, expected) {
		t.Errorf("expectation failed: expected\n%s\ngot\n%s\n", expected, program)
	}
}

func TestParserEOF(t *testing.T) {
	program, err := Parse(`
foo -> bar {
`)
	if err != io.EOF {
		t.Error("expected io.EOF")
	}
	if program != nil {
		t.Error("expected nil as return value")
	}
}

func TestError(t *testing.T) {
	_, err := Parse(`
foo
`)
	perr, ok := err.(ParseError)
	if !ok {
		t.Error("expected err to be ParseError")
	}
	if perr.Line != 2 {
		t.Error("expected error at line 2 not at %d", perr.Line)
	}
}
