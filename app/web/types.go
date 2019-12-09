package web

import "html/template"

type managementPageData struct {
	Target   string
	Password string
	Id       int
	QrCode   string
	Link     string
	Csrf     string
}

type indexPageData struct {
	Csrf string
}

type Template struct {
	templates *template.Template
}
