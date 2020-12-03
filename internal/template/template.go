package template

import "html/template"

func LoadTemplates(templateDir string) (*template.Template, error) {
	return template.ParseGlob(templateDir)
}
