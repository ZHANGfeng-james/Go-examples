package main

import (
	"html/template"
	"log"
	"os"
)

const templateText = `Output 0: {{title .Name1}}
Output 1: {{title .Name2}}
Output 2: {{.Name3 | title}}`

const originText = `Output 0: {{.Name1}}
Output 1: {{.Name2}}
Output 2: {{.Name3}}`

func main() {
	// funcMap := template.FuncMap{"title": strings.Title}
	tpl := template.New("go-programming-tour")
	tpl, err := tpl. /* .Funcs(funcMap) */ Parse(originText)
	if err != nil {
		log.Fatal(err.Error())
	}

	data := map[string]string{
		"Name1": "go",
		"Name2": "programming",
		"Name3": "tour",
	}

	_ = tpl.Execute(os.Stdout, data)
}
