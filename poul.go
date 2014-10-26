package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "poul"
	app.Usage = "A make(1) inspired watching build system"

	app.Run(os.Args)
}
