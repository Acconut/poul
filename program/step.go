package program

import (
	"github.com/Acconut/poul/glob"
)

type Step struct {
	Source       string
	Destination  string
	Code         string
	Dependencies []string
}

type StepMatch struct {
	Step        Step
	Source      string
	Destination string
	Args        map[int]string
}

func (step Step) Builds(dest string) (map[int]string, bool, error) {
	return glob.SimpleMatch(step.Destination, dest)
}

func (step Step) Compiles(source string) (map[int]string, bool, error) {
	return glob.SimpleMatch(step.Source, source)
}

func (step Step) DependsOn(dep string) (bool, error) {
	for _, item := range step.Dependencies {
		pattern, err := glob.NewPattern(item)
		if err != nil {
			return false, err
		}

		_, matches := pattern.Match(dep)
		if matches {
			return true, nil
		}
	}

	return false, nil
}

func (step Step) FindSources() ([]StepMatch, error) {
	matches := make([]StepMatch, 0)
	pattern, err := glob.NewPattern(step.Source)
	if err != nil {
		return matches, err
	}
	entries, err := pattern.Glob()
	if err != nil {
		return matches, err
	}
	for _, entry := range entries {
		matches = append(matches, StepMatch{
			Step:        step,
			Source:      entry.Name,
			Destination: glob.Replace(step.Destination, entry.Args),
			Args:        entry.Args,
		})
	}
	return matches, nil
}
