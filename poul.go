package main

import (
	"encoding/json"
	"github.com/Acconut/poul/parser"
	"github.com/Acconut/poul/program"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	// Disable timestamp
	log.SetFlags(0)
}

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
	log.Println(string(b))
}

func readPoulfile(c *cli.Context) *program.Program {
	name := c.GlobalString("file")
	b, err := ioutil.ReadFile(name)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("unable to read poulfile: file '%s' does not exist", name)
		}
		panic(err)
	}
	prog, err := parser.Parse(string(b))
	if err != nil {
		if perr, ok := err.(parser.ParseError); ok {
			log.Fatalf("unable to parse poulfile: %s", perr)
		}
		panic(err)
	}
	return prog
}

func compile(c *cli.Context) {
	if len(c.Args()) < 0 {
		log.Fatal("no source file(s) supplied.")
	}
	prog := readPoulfile(c)
	code, err := prog.CompileMulti(c.Args()[0:])
	if err != nil {
		if err == program.ErrNoMatch {
			log.Fatal("no build step found.")
		}
		panic(err)
	}
	os.Exit(code)
}

func build(c *cli.Context) {
	if len(c.Args()) < 0 {
		log.Fatal("no destination(s) supplied.")
	}
	prog := readPoulfile(c)
	code, err := prog.BuildMulti(c.Args()[0:])
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func run(c *cli.Context) {
	if len(c.Args()) < 0 {
		log.Fatal("no template supplied.")
	}
	prog := readPoulfile(c)
	code, err := prog.RunTemplate(c.Args()[0])
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}
