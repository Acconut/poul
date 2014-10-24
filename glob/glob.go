package glob

import (
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
		sourcemap[num] = count

		count++
		return `([A-Za-z0-9\-_\.]+)`
	})

	pattern = strings.Replace(pattern, "*", `[A-Za-z0-9\-_\.]+`, -1)

	pattern = "^" + pattern + "$"

	re, err := regexp.Compile(pattern)
	return re, sourcemap, err
}

func Match(pattern string) ([]Entry, error) {
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
		matches := re.FindAllStringSubmatch(file, -1)[0][1:]

		args := make(map[int]string)
		for a, b := range sourcemap {
			args[a] = matches[b]
		}

		entries = append(entries, Entry{
			Name: file,
			Args: args,
		})
	}

	return entries, nil
}
