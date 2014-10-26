package main

import (
	"encoding/json"
	"fmt"
	"github.com/Acconut/poul/parser"
	"github.com/Acconut/poul/program"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "poul"
	app.Usage = "A make(1) inspired watching build system"

	app.Commands = []cli.Command{
		{
			Name:   "dump",
			Usage:  "dump the content of Poulfile to stdout",
			Action: dump,
		},
		{
			Name:   "compile",
			Usage:  "compile a source file",
			Action: compile,
		},
		{
			Name:   "build",
			Usage:  "build a destination",
			Action: build,
		},
		{
			Name:   "run",
			Usage:  "run a templte",
			Action: run,
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "file",
			Value:  "./Poulfile",
			Usage:  "change the Poulfile to read from",
			EnvVar: "POUL_FILE",
		},
	}

	app.Run(os.Args)
}

func dump(c *cli.Context) {
	prog := readPoulfile(c)
	b, err := json.MarshalIndent(prog, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func readPoulfile(c *cli.Context) *program.Program {
	name := c.GlobalString("file")
	b, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	prog, err := parser.Parse(string(b))
	if err != nil {
		panic(err)
	}
	return prog
}

func compile(c *cli.Context) {
	if len(c.Args()) < 0 {
		fmt.Println("no source file(s) supplied.")
		os.Exit(1)
	}
	prog := readPoulfile(c)
	code, err := prog.CompileMulti(c.Args()[0:])
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func build(c *cli.Context) {
	if len(c.Args()) < 0 {
		fmt.Println("no destination supplied.")
		os.Exit(1)
	}
	prog := readPoulfile(c)
	code, err := prog.Build(c.Args()[0])
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func run(c *cli.Context) {
	if len(c.Args()) < 0 {
		fmt.Println("no template supplied.")
		os.Exit(1)
	}
	prog := readPoulfile(c)
	code, err := prog.RunTemplate(c.Args()[0])
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}
