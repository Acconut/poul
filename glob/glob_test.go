package glob

import (
	"reflect"
	"testing"
)

type globTest struct {
	pattern string
	entries []Entry
}

type matchTest struct {
	pattern string
	file    string
	matches bool
	entry   Entry
}

var globTests = []globTest{
	// No parameter
	{
		"./test/bar/bar",
		[]Entry{
			Entry{
				Name: "test/bar/bar",
				Args: make(map[int]string),
			},
		},
	},
	// One parameter and glob parameter
	{
		"./test/$1/*_test",
		[]Entry{
			Entry{
				Name: "test/bar/bar_test",
				Args: map[int]string{
					1: "bar",
				},
			},
			Entry{
				Name: "test/baz/baz_test",
				Args: map[int]string{
					1: "baz",
				},
			},
		},
	},
	// Multiple different parameters
	{
		"./test/$2/$1",
		[]Entry{
			Entry{
				Name: "test/bar/bar",
				Args: map[int]string{
					1: "bar",
					2: "bar",
				},
			},
			Entry{
				Name: "test/bar/bar_test",
				Args: map[int]string{
					1: "bar_test",
					2: "bar",
				},
			},
			Entry{
				Name: "test/baz/baz",
				Args: map[int]string{
					1: "baz",
					2: "baz",
				},
			},
			Entry{
				Name: "test/baz/baz_test",
				Args: map[int]string{
					1: "baz_test",
					2: "baz",
				},
			},
		},
	},
	// The same parameter
	{
		"./test/$1/$1_test",
		[]Entry{
			Entry{
				Name: "test/bar/bar_test",
				Args: map[int]string{
					1: "bar",
				},
			},
			Entry{
				Name: "test/baz/baz_test",
				Args: map[int]string{
					1: "baz",
				},
			},
		},
	},
}

var matchTests = []matchTest{
	// No parameter
	{
		"./test/bar/bar",
		"./test/bar/bar",
		true,
		Entry{
			Name: "test/bar/bar",
			Args: make(map[int]string),
		},
	},
	{
		"./test/bar/bar",
		"./test/bazfoo/bar",
		false,
		Entry{},
	},
	// One parameter and glob parameter
	{
		"./test/$1/*_test",
		"./test/baz/foo_test",
		true,
		Entry{
			Name: "test/baz/foo_test",
			Args: map[int]string{
				1: "baz",
			},
		},
	},
	// Multiple different parameters
	{
		"./test/$2/$1",
		"./test/lol/hihi",
		true,
		Entry{
			Name: "test/lol/hihi",
			Args: map[int]string{
				1: "hihi",
				2: "lol",
			},
		},
	},
	{
		"./test/$2/$1",
		"./tes/lol/hihi",
		false,
		Entry{},
	},
	// The same parameter
	{
		"./test/$1/$1_test",
		"./test/foo/foo_test",
		true,
		Entry{
			Name: "test/foo/foo_test",
			Args: map[int]string{
				1: "foo",
			},
		},
	},
	{
		"./test/$1/$1_test",
		"./test/foo/fba_test",
		false,
		Entry{},
	},
}

func TestGlob(t *testing.T) {
	for _, test := range globTests {
		pattern, err := NewPattern(test.pattern)
		if err != nil {
			t.Errorf("pattern %s failed: %s", test.pattern, err)
		}

		result, err := pattern.Glob()
		if err != nil {
			t.Errorf("pattern %s failed: %s", test.pattern, err)
		}
		if !reflect.DeepEqual(result, test.entries) {
			t.Errorf("pattern %s failed: expected %s\ngot %s", test.pattern, test.entries, result)
		}
	}
}

func TestMatch(t *testing.T) {
	for _, test := range matchTests {
		pattern, err := NewPattern(test.pattern)
		if err != nil {
			t.Errorf("pattern %s failed: %s", test.pattern, err)
		}

		entry, matches := pattern.Match(test.file)
		if matches != test.matches {
			t.Errorf("pattern %s failed: %s != %s", test.pattern, matches, test.matches)
		}
		if !reflect.DeepEqual(entry, test.entry) {
			t.Errorf("pattern %s failed: expected %s\ngot %s", test.pattern, test.entry, entry)
		}
	}
}
