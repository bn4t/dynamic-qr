package web

import (
	"git.bn4t.me/bn4t/dynamic-qr/app/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"log"
)

var server *echo.Echo

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

	// use the csrf middleware
	server.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{TokenLength: 32, TokenLookup: "form:csrf", CookieName: "_csrf", CookieMaxAge: 86400}))

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
	server.GET("/manage/:password", handleManagePage)
	server.GET("/link/:id", handleLink)
	server.POST("/update-qr", handleUpdateQr)
	server.GET("/", handleIndexPage)
	server.Static("/", execDir+"/static/public")
}
