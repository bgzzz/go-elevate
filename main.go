// go-elevate microservice for getting ground height
// based on the requested coordinates
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	cfg := initConfig()

	initLogging(cfg.LogLevel)

	e := echo.New()
	e.Use(LogRequest)
	e.Use(middleware.Recover())

	e.GET("/height", GetOneHeightHandler)
	//decided to make it as POST because it is more human readable
	//while doing manual requests
	e.POST("/heights", GetMultipleHeightsHandler)

	e.Start(cfg.ServerAddr)
}
