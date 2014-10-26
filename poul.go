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