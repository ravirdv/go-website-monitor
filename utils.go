package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// reads config file and returns list of sites
func ReadConfig() Config {
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var config Config
	json.Unmarshal(raw, &config)
	log.Printf("Read configuration, found %d jobs\n", len(config.SiteList))

	return config
}

// setup our log file & stdout
func LogSetupAndDestruct() func() {
	// create file
	logFile, err := os.OpenFile("monitor.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Panicln(err)
	}
	// set logger output
	log.SetOutput(io.MultiWriter(os.Stderr, logFile))

	return func() {
		e := logFile.Close()
		if e != nil {
			fmt.Fprintf(os.Stderr, "Problem closing the log file: %s\n", e)
		}
	}
}
