package testing

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func TestMain(main *testing.M) {
	log.Println("setup")

	main.Run()

	log.Println("teardown")
}

func TestSomeOutput(t *testing.T) {
	type Input struct {
		param string
		want  string
	}

	inputs := []Input{
		{"i", "I"},
	}

	for _, input := range inputs {
		if got := someOutput(input.param); got != input.want {
			t.Fail()
		}
	}
}
