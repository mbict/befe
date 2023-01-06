package main

import (
	. "github.com/mbict/befe/dsl"
	"github.com/mbict/befe/dsl/templates"
)

func Program() Expr {
	template := templates.New(templates.FromString(`{{ .greeting }} {{ .name }}`))

	return template.Render(
		templates.WithData("greeting", String("hello")),
		DefaultParam(ParamFromQuery("name", "name"), String("world")),
	)
}
