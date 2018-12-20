package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

//default consts
const (
	AWS_PREFIX    = "https://s3.amazonaws.com/elevation-tiles-prod/terrarium/"
	SERVER_ADDR   = ":1323"
	DEF_LOG_LEVEL = "debug"
)

//Config is object that holds needed paths for req sending
//and API serving and logging level
type Config struct {
	AWSPrefix  string
	ServerAddr string
	LogLevel   log.Level
}

//initConfig return config object that holds server paths
func initConfig() Config {
	var cfg Config

	cfg.AWSPrefix = os.Getenv("AWS_PREFIX")
	if len(cfg.AWSPrefix) == 0 {
		log.Warning("AWS_PREFIX set to default")
		cfg.AWSPrefix = AWS_PREFIX
	}

	cfg.ServerAddr = os.Getenv("SERVER_ADDR")
	if len(cfg.ServerAddr) == 0 {
		log.Warning("SERVER_ADDR set to default")
		cfg.ServerAddr = SERVER_ADDR
	}

	tmp := os.Getenv("LOG_LEVEL")
	if len(tmp) == 0 {
		log.Warning("LOG_LEVEL set to default")
		tmp = DEF_LOG_LEVEL
	}

	lvl, err := log.ParseLevel(tmp)
	if err != nil {
		panic(err)
	}

	cfg.LogLevel = lvl

	return cfg
}
