package program

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/Acconut/poul/glob"
)

var ErrStepNotFound = errors.New("program: step not found")
var ErrTemplateNotFound = errors.New("program: template not found")
var ErrNoMatch = errors.New("program: no matching step found")

type Program struct {
	Steps     []Step
	Templates map[string]Template
}

type Template struct {
	Name         string
	Prehooks     []string
	Posthooks    []string
	Destinations []string
}

func (prog Program) RunTemplate(name string) (int, error) {
	tpl, ok := prog.Templates[name]
	if !ok {
		return -1, ErrTemplateNotFound
	}

	// Run prehooks
	for _, hook := range tpl.Prehooks {
		code, err := prog.RunTemplate(hook)
		if err != nil || code != 0 {
			return code, err
		}
	}

	// Run steps for destinations
	for _, dest := range tpl.Destinations {
		code, err := prog.Build(dest)
		if err != nil || code != 0 {
			return code, err
		}
	}

	// Run posthooks
	for _, hook := range tpl.Posthooks {
		code, err := prog.RunTemplate(hook)
		if err != nil || code != 0 {
			return code, err
		}
	}

	return 0, nil
}

func (prog Program) Build(dest string) (int, error) {
	for _, step := range prog.Steps {
		args, matches, err := step.Builds(dest)
		if err != nil {
			return -1, err
		}

		if matches {
			sources := glob.ReplaceSlice(step.Sources, args)
			return prog.Run(step, sources, []string{
				dest,
			}, args)
		}
	}

	return -1, ErrNoMatch
}

func (prog Program) BuildMulti(dests []string) (int, error) {
	for _, dest := range dests {
		code, err := prog.Build(dest)
		if code != 0 {
			return code, err
		}
	}
	return 0, nil
}

func (prog Program) Compile(source string) (int, error) {
	hadMatch := false
	for _, step := range prog.Steps {
		args, matches, err := step.Compiles(source)
		if err != nil {
			return -1, err
		}

		if matches {
			hadMatch = true
			dests := glob.ReplaceSlice(step.Destinations, args)
			code, err := prog.Run(step, []string{
				source,
			}, dests, args)
			if code != 0 || err != nil {
				return code, err
			}
		}
	}

	if hadMatch {
		return 0, nil
	}

	return -1, ErrNoMatch
}

func (prog Program) CompileMulti(sources []string) (int, error) {
	for _, source := range sources {
		code, err := prog.Compile(source)
		if code != 0 {
			return code, err
		}
	}
	return 0, nil
}

func (prog Program) Run(step Step, sources, dests []string, args map[int]string) (int, error) {
	cmd := exec.Command("/bin/sh", "-e", "-c", step.Code)

	// Setup environment variables
	env := os.Environ()
	env = append(env, "POUL_SRC="+strings.Join(sources, " "))
	env = append(env, "POUL_DEST="+strings.Join(dests, " "))

	// Setup arguments
	for index, value := range args {
		env = append(env, "POUL_ARG_"+strconv.Itoa(index)+"="+value)
	}

	cmd.Env = env

	// Pipe output to stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	code, ok := getExitCode(err)
	if ok {
		return code, nil
	}

	return -1, err
}

// If possible get the exit code from an error
func getExitCode(err error) (int, bool) {
	if exiterr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0

		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), true
		}
	}

	return 0, true
}
