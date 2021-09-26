package main

import (
	"covidApi/controller"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

//Coordinate is a [longitude, latitude]
type Statelist struct {
	Statelist []State `json:"statewise"`
}

type State struct {
	Active          string `json:"active"`
	Confirmed       string `json:"confirmed"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	State           string `json:"state"`
}

type StateMapBox struct {
	StateMapBox []StateCoor `json:"features"`
}

type StateCoor struct {
	StateName string `json:"text"`
}

func stateCovid(c echo.Context) error {
	// Get team and member from the query string
	var x string = c.QueryParam("longitude")
	var y string = c.QueryParam("latitude")

	//longitude, err := strconv.ParseFloat(x, 8)
	//latitude, err := strconv.ParseFloat(y, 8)

	//log.Printf(longitude, latitude, err)
	fmt.Print(x, y)

	// var state string = controller.CordToState(x, y)
	// fmt.Print(state)
	// cases := controller.Statecases(state)
	// fmt.Print(cases)
	// casesMarshal, err := json.Marshal(cases)
	// fmt.Print(err)

	var resp string = controller.Statecases(x, y)

	return c.String(http.StatusOK, resp)
}

func fetchtoDB(c echo.Context) error {
	controller.Update()
	fmt.Print("updated")
	return c.String(http.StatusOK, "Updated")
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/state", stateCovid)
	e.GET("/fetchtodb", fetchtoDB)
	e.Logger.Fatal(e.Start(":1323"))

}
