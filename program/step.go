package program

import (
	"github.com/Acconut/poul/glob"
)

type Step struct {
	Sources      []string
	Destinations []string
	Code         string
}

func (step Step) Builds(dest string) (map[int]string, bool, error) {
	return globSlice(dest, step.Destinations)
}

func (step Step) Compiles(source string) (map[int]string, bool, error) {
	return globSlice(source, step.Sources)
}

func globSlice(element string, slice []string) (map[int]string, bool, error) {
	for _, item := range slice {
		pattern, err := glob.NewPattern(item)
		if err != nil {
			return nil, false, err
		}

		entry, matches := pattern.Match(element)
		if matches {
			return entry.Args, true, nil
		}
	}

	return nil, false, nil
}
