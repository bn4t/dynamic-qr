package web

import "github.com/labstack/echo/v4"

var server *echo.Echo

func RegisterRoutes() {

}

func Start() {
	server = echo.New()

	RegisterRoutes()

	server.Logger.Fatal(server.Start(":3000"))
}
