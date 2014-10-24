package glob

import (
	"reflect"
	"testing"
)

type globTest struct {
	pattern string
	entries []Entry
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

func TestGlob(t *testing.T) {
	for _, test := range globTests {
		result, err := Match(test.pattern)
		if err != nil {
			t.Errorf("pattern %s failed: %s", test.pattern, err)
		}

		if !reflect.DeepEqual(result, test.entries) {
			t.Errorf("pattern %s failed: expected %s\ngot %s", test.pattern, test.entries, result)
		}
	}
}
