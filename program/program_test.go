package program

import (
	"fmt"
	"testing"
)

var prog = Program{
	Templates: map[string]Template{
		"echo": Template{
			Name: "echo",
			Destinations: []string{
				"foo/boo",
			},
		},
	},
	Steps: []Step{
		Step{
			Sources: []string{
				"src/$1",
				"src/package",
			},
			Destinations: []string{
				"foo/$1",
			},
			Code: `echo "Hello world!"
echo "Compiling ${POUL_SRC}."
echo "Building ${POUL_DEST}."
echo "Arg #1: ${POUL_ARG_1}."
printenv`,
		},
	},
}

func TestBuild(t *testing.T) {
	code, err := prog.Build("foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(code)
}

func TestCompile(t *testing.T) {
	code, err := prog.Compile("src/lol")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(code)
}

func TestCompileMulti(t *testing.T) {
	code, err := prog.CompileMulti([]string{
		"src/lol",
		"src/foo",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(code)
}

func TestRunTemplate(t *testing.T) {
	code, err := prog.RunTemplate("echo")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(code)
}
