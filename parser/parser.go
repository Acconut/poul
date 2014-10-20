package parser

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	Newline       = "\n"
	Comment uint8 = '#'
)

const (
	PartNone = 1 << iota
	PartName
	PartSource
	PartCode
)

var (
	ReName = regexp.MustCompile(`^([A-Za-z0-9_\-]+)\s*(\([^\)]+\))?:$`)
)

type Program []Step

type Step struct {
	Name         string
	Prehooks     []string
	Posthooks    []string
	Sources      []string
	Destinations []string
	Code         string
}

func Parse(code string) (*Program, error) {

	// Split code into lines
	lines := strings.Split(code, Newline)

	program := make(Program, 0)

	part := PartNone
	currentStep := Step{}
	buffer := ""

	for lineNumber, line := range lines {
		// Trim line
		line = strings.TrimSpace(line)

		// Ignore empty lines
		if len(line) == 0 {
			continue
		}

		// Ignore comments, e.g.
		// # That's a comment
		if line[0] == Comment {
			continue
		}

		// We expect a template name declaration
		if part == PartNone {
			if !ReName.Match([]byte(line)) {
				return nil, fmt.Errorf("Expected name declaration at line %d", lineNumber+1)
			}

			result := ReName.FindStringSubmatch(line)

			templateName := result[1]
			hooks := result[2]

			if hooks != "" {
			// Hooks include brackets so we remove
			// the first and last char
				preHooks, postHooks := parseHooks(hooks[1:len(hooks) - 1])

				currentStep.Prehooks = preHooks
				currentStep.Posthooks = postHooks
			}

			currentStep.Name = templateName

			part = PartName
			continue
		}

		// We expect a source declaration
		if part == PartName {
			if !strings.Contains(line, "->") || line[len(line)-1] != uint8('{') {
				return nil, fmt.Errorf("Expected source declaration at line %d", lineNumber+1)
			}

			sources, dests := parseSources(line[0:len(line)-1])

			currentStep.Sources = sources
			currentStep.Destinations = dests

			// Clear buffer
			buffer = ""

			part = PartSource
			continue
		}

		// We expect a source end or code
		if part == PartSource {
			// Code ends
			if line == "}" {
				currentStep.Code = buffer
				part = PartNone

				program = append(program, currentStep)
				currentStep = Step{}
				continue
			}

			buffer += line + Newline
			continue
		}

	}

	return &program, nil
}

func parseHooks(hooks string) ([]string, []string) {
	return split(hooks, "/", ",")
}

func parseSources(line string) ([]string, []string) {
	return split(line, "->", ",")
}

func split(line, firstSep, secondSep string) ([]string, []string) {
	parts := strings.Split(line, firstSep)

	first := strings.Split(parts[0], secondSep)
	second := strings.Split(parts[1], secondSep)

	for index, part := range first {
		first[index] = strings.TrimSpace(part)
	}

	for index, part := range second {
		second[index] = strings.TrimSpace(part)
	}

	return first, second
}
