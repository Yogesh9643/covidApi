package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func fetchtoDB(c echo.Context) error {

	return c.String(http.StatusOK, "Updated")
}
func stateCovid(c echo.Context) error {
	var x string = c.QueryParam("longitude")
	var y string = c.QueryParam("latitude")
	fmt.Print(x, y)
	var cases State
	state := GetStateData(x, y)
	fmt.Print(state)
	cases = stateCases(state)
	fmt.Print(cases)
	casesMarshal, err := json.Marshal(cases)
	fmt.Print(err)
	return c.String(http.StatusOK, string(casesMarshal))
}
