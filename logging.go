package main

import (
	"os"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// initLogging sets up logger with
// appropriate logging level
func initLogging(lvl log.Level) {

	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	log.SetLevel(lvl)

	log.SetReportCaller(true)

}

// LogRequest is middleware function that logs handlers performance
// and outputs req/rsp meta data
func LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {

		//get the objects
		req := c.Request()
		res := c.Response()
		start := time.Now()

		// wait while handler will be processed
		if err = next(c); err != nil {
			log.Error(err.Error())
			c.Error(err)
		}

		// calculate the time and add it to the log
		stop := time.Now()
		log.WithFields(log.Fields{
			"method":    req.Method,
			"remote_ip": c.RealIP(),
			"uri":       req.RequestURI,
			"status":    res.Status,
			"latency":   stop.Sub(start).String(),
		}).Info("request processed")

		return nil
	}
}
