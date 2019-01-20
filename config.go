package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

//default consts
const (
	// AWS_PREFIX    = "https://s3.amazonaws.com/elevation-tiles-prod/terrarium/"
	SERVER_ADDR   = ":1323"
	DEF_LOG_LEVEL = "debug"
)

// Config is object that holds needed paths for req sending
// and API serving and logging level
type Config struct {
	// AWSPrefix  string
	ServerAddr string
	LogLevel   log.Level
}

// initConfig returns config object that holds server paths
func initConfig() Config {
	var cfg Config

	cfg.ServerAddr = os.Getenv("SERVER_ADDR")
	if len(cfg.ServerAddr) == 0 {
		log.Debug("SERVER_ADDR set to default")
		cfg.ServerAddr = SERVER_ADDR
	}

	lvlStr := os.Getenv("LOG_LEVEL")
	if len(lvlStr) == 0 {
		log.Debug("LOG_LEVEL set to default")
		lvlStr = DEF_LOG_LEVEL
	}

	lvl, err := log.ParseLevel(lvlStr)
	if err != nil {
		panic(err)
	}

	cfg.LogLevel = lvl

	return cfg
}
