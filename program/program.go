package program

type Program []Step

type Step struct {
	Name         string
	Prehooks     []string
	Posthooks    []string
	Sources      []string
	Destinations []string
	Code         string
}
