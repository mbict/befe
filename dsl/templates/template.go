package templates

import (
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/expr"
	"html/template"
	"net/http"
)

type Option func(t *template.Template)

func FromFile(filenames ...string) Option {
	return func(t *template.Template) {
		if _, err := t.ParseFiles(filenames...); err != nil {
			panic(err)
		}
	}
}

func FromString(templateStr ...string) Option {
	return func(t *template.Template) {
		for _, s := range templateStr {
			if _, err := t.Parse(s); err != nil {
				panic(err)
			}
		}
	}
}

func FromPattern(pattern string) Option {
	return func(t *template.Template) {
		if _, err := t.ParseGlob(pattern); err != nil {
			panic(err)
		}
	}
}

func WithData(name string, value Valuer) Param {
	return dsl.WithParam(name, value)
}

type HTMLTemplate interface {
	RenderTemplate(name string, data ...Param) Action
	Render(data ...Param) Action
}

func New(options ...Option) HTMLTemplate {
	funcMap := sprig.FuncMap()

	funcMap["toHTML"] = func(s string) template.HTML {
		return template.HTML(s)
	}

	funcMap["toJS"] = func(s string) template.JS {
		return template.JS(s)
	}

	t := template.New("").Funcs(funcMap)
	for _, option := range options {
		option(t)
	}
	return &htmlTemplater{t: t}
}

type htmlTemplater struct {
	t *template.Template
}

func (h *htmlTemplater) Render(dataFn ...Param) Action {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		err := h.t.Execute(rw, generateTemplateData(dataFn, r))
		if err != nil {
			fmt.Println("template has error: ", err)
		}
		return err == nil, err

	})
}

func (h *htmlTemplater) RenderTemplate(name string, dataFn ...Param) Action {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		err := h.t.ExecuteTemplate(rw, name, generateTemplateData(dataFn, r))
		if err != nil {
			fmt.Println("template has error: ", err)
		}
		return err == nil, err
	})
}

func generateTemplateData(dataFn []Param, r *http.Request) map[string]interface{} {
	templateData := make(map[string]interface{})
	for _, df := range dataFn {
		k, v := df(r)
		templateData[k] = v
	}
	return templateData
}
