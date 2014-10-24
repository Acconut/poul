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

var rePattern = regexp.MustCompile(`\$(\d+)`)

func toGlob(pattern string) string {
	return rePattern.ReplaceAllStringFunc(pattern, toGlobReplacer)
}

func toGlobReplacer(a string) string {
	return "*"
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

func Match(pattern string) ([]Entry, error) {
	pattern = path.Clean(pattern)

	entries := make([]Entry, 0)

	re, sourcemap, err := toRegexp(pattern)
	if err != nil {
		return entries, err
	}

	files, err := filepath.Glob(toGlob(pattern))
	if err != nil {
		return entries, err
	}

	for _, file := range files {
		file = path.Clean(file)

		if !re.Match([]byte(file)) {
			continue
		}
		matches := re.FindAllStringSubmatch(file, -1)[0][1:]

		skip := false
		args := make(map[int]string)
		for a, b := range sourcemap {
			value, ok := args[b]

			if ok && value != matches[a] {
				skip = true
				break
			}
			args[b] = matches[a]
		}

		if skip {
			continue
		}

		entries = append(entries, Entry{
			Name: file,
			Args: args,
		})
	}

	return entries, nil
}
