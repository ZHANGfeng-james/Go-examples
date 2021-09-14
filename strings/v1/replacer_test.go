package gostrings

import (
	"os"
	"strings"
	"testing"
)

func TestReplacer(t *testing.T) {
	replacer := strings.NewReplacer("<", "&lt;", ">", "&gt")
	// fmt.Println(replacer.Replace("This is <b>HTML</b>!"))

	replacer.WriteString(os.Stdout, "This is <b>HTML</b>!\n")
}
