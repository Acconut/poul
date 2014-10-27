package program

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
			Code: `
echo "POUL_SRC: ${POUL_SRC}"
echo "POUL_DEST: ${POUL_DEST}"
echo "POUL_ARG_1: ${POUL_ARG_1}"`,
		},
	},
}

func ExampleProgram_Build() {
	code, err := prog.Build("foo/bar")
	if err != nil {
		panic(err)
	}
	if code != 0 {
		panic("not null")
	}
	// Output:
	// POUL_SRC: src/bar src/package
	// POUL_DEST: foo/bar
	// POUL_ARG_1: bar
}

func ExampleProgram_Compile() {
	code, err := prog.Compile("src/lol")
	if err != nil {
		panic(err)
	}
	if code != 0 {
		panic("not null")
	}
	// Output:
	// POUL_SRC: src/lol
	// POUL_DEST: foo/lol
	// POUL_ARG_1: lol
}

func ExampleProgram_CompileMulti() {
	code, err := prog.CompileMulti([]string{
		"src/lol",
		"src/foo",
	})
	if err != nil {
		panic(err)
	}
	if code != 0 {
		panic("not null")
	}
	// Output:
	// POUL_SRC: src/lol
	// POUL_DEST: foo/lol
	// POUL_ARG_1: lol
	// POUL_SRC: src/foo
	// POUL_DEST: foo/foo
	// POUL_ARG_1: foo
}

func ExampleProgram_RunTemplate() {
	code, err := prog.RunTemplate("echo")
	if err != nil {
		panic(err)
	}
	if code != 0 {
		panic("not null")
	}
	// Output:
	// POUL_SRC: src/boo src/package
	// POUL_DEST: foo/boo
	// POUL_ARG_1: boo
}
