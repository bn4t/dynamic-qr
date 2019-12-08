package web

import (
	"git.bn4t.me/bn4t/dynamic-qr/app/utils"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log"
)

var server *echo.Echo

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// start labstack echo
func Start() {
	server = echo.New()

	// get current execution directory
	execDir, err := utils.GetExecutionDir()
	if err != nil {
		log.Fatal(err)
	}

	// register web routes
	RegisterRoutes(execDir)

	// setup template renderer
	t := &Template{
		templates: template.Must(template.ParseGlob(execDir + "/static/templates/*.html")),
	}
	server.Renderer = t

	server.Logger.Fatal(server.Start(":" + utils.GetEnv("PORT", "3000")))
}

func RegisterRoutes(execDir string) {
	server.POST("/create-qr", handleCreateQr)
	server.GET("/manage/:password", handleManage)
	server.GET("/link/:id", handleLink)
	server.POST("/update-qr", handleUpdateQr)
	server.Static("/", execDir+"/static/public")
}
