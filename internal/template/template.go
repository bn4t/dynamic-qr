package template

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

var templateNames = map[string]string{
	"index":  "/templates/index.html",
	"manage": "/templates/manage.html",
}

func LoadTemplates(statikFs http.FileSystem) (*template.Template, error) {
	var tmpl *template.Template

	// load all templates
	for name, path := range templateNames {
		f, err := statikFs.Open(path)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}

		if tmpl == nil {
			tmpl, err = template.New(name).Parse(string(b))
			if err != nil {
				return nil, err
			}
		} else {
			tmpl, err = tmpl.New(name).Parse(string(b))
			if err != nil {
				return nil, err
			}
		}
	}

	return tmpl, nil
}
