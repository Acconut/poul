package glob

import (
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Entry struct {
	Name string
	Args map[int]string
}

type Pattern struct {
	re *regexp.Regexp
	sourcemap map[int]int
	glob string
	Pattern string
}

var rePattern = regexp.MustCompile(`\$(\d+)`)

func NewPattern(pattern string) (*Pattern, error) {
	// Remove ./ from the beginning
	pattern = path.Clean(pattern)

	// Transform the pattern into
	// filepath.Glob's one
	glob := rePattern.ReplaceAllStringFunc(pattern, func(a string) string {
		return "*"
	})

	// Compile it into a regexp
	re, sourcemap, err := toRegexp(pattern)
	if err != nil {
		return nil, err
	}

	return &Pattern{
		re: re,
		sourcemap: sourcemap,
		glob: glob,
		Pattern: pattern,
	}, nil
}

func toRegexp(pattern string) (*regexp.Regexp, map[int]int, error) {
	count := 0
	sourcemap := make(map[int]int)

	pattern = rePattern.ReplaceAllStringFunc(pattern, func(match string) string {
		num, _ := strconv.Atoi(match[1:])
		sourcemap[count] = num

		count++
		return `([A-Za-z0-9\-_\.]+)`
	})

	pattern = strings.Replace(pattern, "*", `[A-Za-z0-9\-_\.]+`, -1)

	pattern = "^" + pattern + "$"

	re, err := regexp.Compile(pattern)
	return re, sourcemap, err
}

func (pattern Pattern) Glob() ([]Entry, error) {
	entries := make([]Entry, 0)

	files, err := filepath.Glob(pattern.glob)
	if err != nil {
		return entries, err
	}

	for _, file := range files {
		file = path.Clean(file)

		entry, matches := pattern.Match(file)
		if matches {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func (pattern Pattern) Match(file string) (Entry, bool) {
	entry := Entry{}
	file = path.Clean(file)

	if !pattern.re.Match([]byte(file)) {
		return entry, false
	}
	matches := pattern.re.FindAllStringSubmatch(file, -1)[0][1:]

	args := make(map[int]string)
	for a, b := range pattern.sourcemap {
		value, ok := args[b]

		if ok && value != matches[a] {
			return entry, false
		}
		args[b] = matches[a]
	}


	entry.Name = file
	entry.Args = args

	return entry, true
}

func (pattern Pattern) String() string {
	return pattern.Pattern
}
