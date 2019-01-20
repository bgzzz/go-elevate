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
	log "github.com/sirupsen/logrus"
)

func main() {

	cfg := initConfig()

	initLogging(cfg.LogLevel)

	e := echo.New()
	e.Use(LogRequest)
	e.Use(middleware.Recover())

	//to use this API you have to encode your requests
	e.GET("/heights", GetHeightsHandler)
	//decided to make it as POST because it is more human readable
	//while doing manual requests
	e.POST("/heights", PostMultipleHeightsHandler)

	err := e.Start(cfg.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
}
