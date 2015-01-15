package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	extractAddress string
	entitiesPath   string
	logPath        string
}

func NewConfig() *Config {
	cfg := new(Config)

	cfg.extractAddress = getenvDefault("EXTRACTOR_EXTRACT_ADDR", ":3096")
	cfg.entitiesPath = getenvDefault("EXTRACTOR_ENTITIES_PATH", "/var/apps/entity-extractor/data/entities.jsonl")
	cfg.logPath = getenvDefault("EXTRACTOR_LOG_PATH", "STDERR")

	flag.Usage = usage
	flag.Parse()

	return cfg
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s\n", os.Args[0])
	helpstring := `
The following environment variables and defaults are available:

EXTRACTOR_EXTRACT_ADDR=:3096  Address on which to serve extraction requests
EXTRACTOR_ENTITIES_PATH=/var/apps/entity-extractor/data/entities.jsonl
                              Path of file holding entities in jsonlines format
EXTRACTOR_ERROR_LOG=STDERR    File to log errors to (in JSON format)
`
	fmt.Fprintf(os.Stderr, helpstring)
	os.Exit(2)
}

func getenvDefault(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultVal
	}
	return val
}
