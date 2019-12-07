package web

import (
	"git.bn4t.me/bn4t/dynamic-qr/app/utils"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

var server *echo.Echo

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Start() {
	server = echo.New()

	// register web routes
	RegisterRoutes()

	t := &Template{
		templates: template.Must(template.ParseGlob("static/templates/*.html")),
	}
	server.Renderer = t

	server.Logger.Fatal(server.Start(":" + utils.Getenv("PORT", "3000")))
}

func RegisterRoutes() {
	server.POST("/create-qr", handleCreateQr)
	server.GET("/manage/:password", handleManage)
	server.GET("/link/:id", handleLink)
	server.POST("/update-qr", handleUpdateQr)
	server.Static("/", "static/public")
}
