package program

var prog = Program{
	Templates: map[string]Template{
		"echo": Template{
			Name: "echo",
			Destinations: []string{
				"test/out/boo",
			},
		},
	},
	Steps: []Step{
		Step{
			Source: "test/$1.txt",
			Dependencies: []string{
				"test/package",
			},
			Destination: "test/out/$1",
			Code: `
echo "POUL_SRC: ${POUL_SRC}"
echo "POUL_DEST: ${POUL_DEST}"
echo "POUL_ARG_1: ${POUL_ARG_1}"`,
		},
	},
}

func ExampleProgram_Build() {
	code, err := prog.Build("test/out/bar")
	if err != nil {
		panic(err)
	}
	if code != 0 {
		panic("not null")
	}
	// Output:
	// POUL_SRC: test/bar.txt
	// POUL_DEST: test/out/bar
	// POUL_ARG_1: bar
}

func ExampleProgram_Compile() {
	code, err := prog.Compile("test/foo.txt")
	if err != nil {
		panic(err)
	}
	if code != 0 {
		panic("not null")
	}
	// Output:
	// POUL_SRC: test/foo.txt
	// POUL_DEST: test/out/foo
	// POUL_ARG_1: foo
}

func ExampleProgram_CompileMulti() {
	code, err := prog.CompileMulti([]string{
		"test/foo.txt",
		"test/bar.txt",
	})
	if err != nil {
		panic(err)
	}
	if code != 0 {
		panic("not null")
	}
	// Output:
	// POUL_SRC: test/foo.txt
	// POUL_DEST: test/out/foo
	// POUL_ARG_1: foo
	// POUL_SRC: test/bar.txt
	// POUL_DEST: test/out/bar
	// POUL_ARG_1: bar
}

func ExampleProgram_RunTemplate() {
	code, err := prog.RunTemplate("echo")
	if err != nil {
		panic(err)
	}
	if code != 0 {
		panic("not null")
	}
	// Output
	// POUL_SRC: test/boo.txt
	// POUL_DEST: test/out/boo
	// POUL_ARG_1: boo
}

func ExampleProgram_CompileByDependency() {
	code, err := prog.CompileByDependency("test/package")
	if err != nil {
		panic(err)
	}
	if code != 0 {
		panic("not null")
	}
	// Output
	// POUL_SRC: test/bar.txt
	// POUL_DEST: test/out/bar
	// POUL_ARG_1: bar
	// POUL_SRC: test/foo.txt
	// POUL_DEST: test/out/foo
	// POUL_ARG_1: foo

}
