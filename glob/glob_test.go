package glob

import (
	"fmt"
	"testing"
)

func TestFoo(t *testing.T) {
	fmt.Println(toGlob("/foo/$2/$1.html"))

	fmt.Println(Match("../*/$1_test.go"))

	fmt.Println(Match("../$2/$1.go"))
}
