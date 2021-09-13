package view

import (
	"os"
	"text/template"
)

type Person struct {
	UserName string
}

func templateView() {
	t := template.New("Template")
	t.Parse("hello {{.UserName}}!")
	p := Person{
		UserName: "Katyusha",
	}
	t.Execute(os.Stdout, p)
}
