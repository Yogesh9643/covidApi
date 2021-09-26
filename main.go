package main

import (
	"covidApi/controller"
	"fmt"
	"net/http"

	_ "covidApi/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func stateCovid(c echo.Context) error {
	var x string = c.QueryParam("longitude")
	var y string = c.QueryParam("latitude")
	fmt.Print(x, y)
	var resp string = controller.Statecases(x, y)
	return c.String(http.StatusOK, resp)
}

func fetchtoDB(c echo.Context) error {
	controller.Update()
	fmt.Print("updated")
	return c.String(http.StatusOK, "Updated")
}

// @title Covid Api
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1323
// @BasePath /
// @schemes http
func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/state", stateCovid)
	e.GET("/fetchtodb", fetchtoDB)
	e.Logger.Fatal(e.Start(":1323"))
}
