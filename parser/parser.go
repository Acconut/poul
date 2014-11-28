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

var (
	ReTemplateStart = regexp.MustCompile(`^([A-Za-z0-9_\-]+)\s*(\([^\)]+\))?$`)
	ReStepName      = regexp.MustCompile(`^([^\(]+)(\(([^\)]+)\))?\s*->\s(.+)$`)
)

func Parse(code string) (*prog.Program, error) {

	// Split code into lines
	lines := strings.Split(code, Newline)

	program := prog.Program{
		Steps:     make([]prog.Step, 0),
		Templates: make(map[string]prog.Template),
	}

	inBlock := false
	name := ""
	body := ""
	blockStart := 0

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

		if !inBlock {
			// We currently aren't in a block and expect
			// a block beginning (line ending with opening bracket).
			if line[len(line)-1] != BracketOpen {
				return nil, ParseError{
					lineNumber + 1,
					"Expected block declaration",
				}
			}

			// Store trimed line without brackets as block name
			name = strings.TrimSpace(line[:len(line)-1])

			// Store line in which the block start
			blockStart = lineNumber
			inBlock = true
			continue
		} else {
			// When the line is a closing brackets
			// we have a block end
			if line == BracketClose {
				err := parseBlock(&program, name, body, blockStart)
				if err != nil {
					return nil, err
				}

				// Reset name, body and position
				inBlock = false
				name = ""
				body = ""
				continue
			}

			// The current line is part of the body
			body += line + Newline
		}

	}

	if inBlock {
		return nil, io.EOF
	}

	return &program, nil
}

func parseHooks(hooks string) ([]string, []string) {
	return split(hooks, Slash, Comma)
}

func parseSources(line string) (string, string, []string) {
	result := ReStepName.FindStringSubmatch(line)

	source := strings.TrimSpace(result[1])
	deps := splitSingle(result[3], ",")
	dest := strings.TrimSpace(result[4])

	// Don't return an array containing an empty string
	if len(deps) == 1 && deps[0] == "" {
		deps = nil
	}

	return source, dest, deps
}

func splitSingle(line, sep string) []string {
	parts := strings.Split(line, sep)

	for index, part := range parts {
		parts[index] = strings.TrimSpace(part)
	}

	return parts
}

func split(line, firstSep, secondSep string) ([]string, []string) {
	parts := strings.Split(line, firstSep)

	first := splitSingle(parts[0], secondSep)
	second := splitSingle(parts[1], secondSep)

	return first, second
}

func parseBlock(program *prog.Program, name, body string, lineNr int) error {
	if strings.Contains(name, Arrow) {
		// We found a step declaration (a line containing the arrow ->)
		source, dest, deps := parseSources(name)

		step := prog.Step{
			Source:       source,
			Destination:  dest,
			Code:         body,
			Dependencies: deps,
		}

		program.Steps = append(program.Steps, step)

		return nil
	}

	if ReTemplateStart.Match([]byte(name)) {
		// We found a template start
		result := ReTemplateStart.FindStringSubmatch(name)

		template := prog.Template{
			Name: result[1],
		}
		hooks := result[2]

		if hooks != "" {
			// Hooks include brackets so we remove
			// the first and last char
			preHooks, postHooks := parseHooks(hooks[1 : len(hooks)-1])

			template.Prehooks = preHooks
			template.Posthooks = postHooks
		}

		template.Destinations = strings.Split(strings.TrimSpace(body), Newline)

		program.Templates[template.Name] = template

		return nil
	}

	return ParseError{
		lineNr + 1,
		"Unknown block start",
	}
}
