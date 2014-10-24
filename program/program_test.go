package program

import (
	"testing"
	"fmt"
)

var prog = Program{
	Step{
		Name: "echo",
		Prehooks: []string{
			"pre",
		},
		Code: `echo "Hello world!"
echo "Running step ${POUL_STEP}."
exit 4`,
	},
	Step{
		Name: "pre",
		Code: `echo "I'm pre: ${POUL_STEP}"`,
	},
}

func TestFindName(t *testing.T) {
	step, ok := prog.FindName("echo")
	if ok != true {
		t.Error("ok should be true")
	}
	if step.Name != "echo" {
		t.Error("wrong name returned")
	}

	step, ok = prog.FindName("doesnotexist")
	if ok != false {
		t.Error("ok should be false")
	}
	if step.Name != "" {
		t.Error("empty step should be returned")
	}
}

func TestRun(t *testing.T) {
	code, err := prog.RunName("echo")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(code)
}
