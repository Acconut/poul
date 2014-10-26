package parser

import (
	prog "github.com/Acconut/poul/program"
	"io"
	"regexp"
	"strings"
)

const (
	Newline            = "\n"
	Comment      uint8 = '#'
	Arrow              = "->"
	Comma              = ","
	Slash              = "/"
	BracketOpen  uint8 = '{'
	BracketClose       = "}"
)

const (
	PartNone = 1 << iota
	PartTemplateStart
	PartStepStart
)

var (
	ReTemplateStart = regexp.MustCompile(`^([A-Za-z0-9_\-]+)\s*(\([^\)]+\))?\s*{$`)
)

func Parse(code string) (*prog.Program, error) {

	// Split code into lines
	lines := strings.Split(code, Newline)

	program := prog.Program{
		Steps:     make([]prog.Step, 0),
		Templates: make(map[string]prog.Template),
	}

	part := PartNone
	currentStep := prog.Step{}
	currentTemplate := prog.Template{}
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

		// We expect a template or step start declaration
		if part == PartNone {
			if ReTemplateStart.Match([]byte(line)) {
				// We found a template start
				result := ReTemplateStart.FindStringSubmatch(line)

				templateName := result[1]
				hooks := result[2]

				if hooks != "" {
					// Hooks include brackets so we remove
					// the first and last char
					preHooks, postHooks := parseHooks(hooks[1 : len(hooks)-1])

					currentTemplate.Prehooks = preHooks
					currentTemplate.Posthooks = postHooks
				}

				currentTemplate.Name = templateName

				part = PartTemplateStart
				continue
			} else if strings.Contains(line, Arrow) && line[len(line)-1] == BracketOpen {
				// We found a step declaration
				sources, dests := parseSources(line[0 : len(line)-1])

				currentStep.Sources = sources
				currentStep.Destinations = dests

				part = PartStepStart
				continue
			} else {
				return nil, ParseError{
					lineNumber + 1,
					"Expected template or step declaration",
				}
			}
		}

		// We expect a body end or content
		if part == PartTemplateStart || part == PartStepStart {
			// Code ends
			if line == BracketClose {
				if part == PartStepStart {
					currentStep.Code = buffer
					part = PartNone

					program.Steps = append(program.Steps, currentStep)
					currentStep = prog.Step{}
				} else if part == PartTemplateStart {
					currentTemplate.Destinations = strings.Split(strings.TrimSpace(buffer), Newline)
					part = PartNone

					program.Templates[currentTemplate.Name] = currentTemplate
					currentTemplate = prog.Template{}
				}

				// Clear buffer
				buffer = ""
				continue
			}

			buffer += line + Newline
			continue
		}

	}

	if part != PartNone {
		return nil, io.EOF
	}

	return &program, nil
}

func parseHooks(hooks string) ([]string, []string) {
	return split(hooks, Slash, Comma)
}

func parseSources(line string) ([]string, []string) {
	return split(line, Arrow, Comma)
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
