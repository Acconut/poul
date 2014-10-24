package program

import (
	"os"
	"os/exec"
	"syscall"
	"errors"
	"strings"
)

var ErrStepNotFound = errors.New("program: step not found")

type Program []Step

type Step struct {
	Name         string
	Prehooks     []string
	Posthooks    []string
	Sources      []string
	Destinations []string
	Code         string
}

func (prog Program) FindName(name string) (Step, bool) {
	for _, step := range prog {
		if step.Name == name {
			return step, true
		}
	}
	return Step{}, false
}

func (prog Program) RunName(name string) (int, error) {
	return prog.Run(name, nil, nil, nil)
}

func (prog Program) Run(name string, sources, dests, args []string) (int, error) {
	step, ok := prog.FindName(name)
	if !ok {
		return -1, ErrStepNotFound
	}

	// Run prehooks
	for _, hook := range step.Prehooks {
		code, err := prog.RunName(hook)
		if err != nil || code != 0 {
			return code, err
		}
	}

	cmd := exec.Command("/bin/sh", "-c", step.Code)

	// Setup environment variables
	env := make([]string, len(args) + 3)
	env[0] = "POUL_STEP=" + step.Name
	env[1] = "POUL_SOURCES=" + strings.Join(sources, " ")
	env[2] = "POUL_DESTINATIONS=" + strings.Join(dests, " ")

	// Setup arguments
	i := 3
	for index, value := range args {
		env[i] = "POUL_ARG_" + string(index) + "=" + value
		i++
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

	// Run posthooks
	for _, hook := range step.Posthooks {
		code, err := prog.RunName(hook)
		if err != nil || code != 0 {
			return code, err
		}
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