package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Acconut/poul/parser"
	"github.com/Acconut/poul/program"
	"github.com/codegangsta/cli"
	"gopkg.in/fsnotify.v1"
)

var (
	stderr = log.New(os.Stderr, "--> ", 0)
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
		{
			Name:   "watch",
			Usage:  "watch a directory for changes on sources and recompile",
			Action: watch,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "exclude",
					Usage:  "exclude directories from being watched",
					EnvVar: "POUL_EXCLUDE",
				},
			},
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

func watch(c *cli.Context) {
	prog := readPoulfile(c)
	dir := "./"
	if len(c.Args()) > 0 {
		dir = c.Args()[0]
	}
	excludes := excludeMap(c.String("exclude"))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	recentChanges := make(map[string]time.Time)

	done := make(chan bool)
	go func() {
		for {
			select {
			case evt := <-watcher.Events:
				if !isChangeOp(evt.Op) {
					continue
				}
				recentChanges[evt.Name] = time.Now()
			case err := <-watcher.Errors:
				log.Fatal(err)
			case <-time.Tick(1 * time.Second):
				now := time.Now()
				for name, lastEvent := range recentChanges {
					diff := now.Sub(lastEvent)
					if diff > (time.Millisecond * 500) {
						delete(recentChanges, name)
						go processFile(prog, name)
					}
				}
			}
		}
	}()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		// Skip excluded dirs
		if _, ok := excludes[filepath.Clean(path)]; ok {
			return filepath.SkipDir
		}
		stderr.Printf("Watching directory '%s'.\n", path)

		return watcher.Add(path)
	})
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func isChangeOp(op fsnotify.Op) bool {
	return op&fsnotify.Create == fsnotify.Create ||
		op&fsnotify.Write == fsnotify.Write ||
		op&fsnotify.Rename == fsnotify.Rename
}

func excludeMap(str string) map[string]bool {
	arr := strings.Split(str, ",")
	Map := make(map[string]bool)
	for _, value := range arr {
		if value == "" {
			continue
		}
		Map[filepath.Clean(value)] = true
	}
	return Map
}

func processFile(prog *program.Program, fileName string) {
	stderr.Println("")
	stderr.Printf("Event(%s): recompiling...", fileName)
	code, err := prog.Compile(fileName)
	processCompilation(code, err, "No build step found.")

	stderr.Println("Recompiling sources by dependency...")
	code, err = prog.CompileByDependency(fileName)
	processCompilation(code, err, "Not as dependency used.")
}

func processCompilation(code int, err error, message string) {
	if err != nil {
		if err == program.ErrNoMatch {
			stderr.Println(message)
		} else {
			log.Fatal(err)
		}
	}
	if err != program.ErrNoMatch {
		stderr.Printf("(%d)\n", code)
	}
}
