package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func ReadConfig() Jobs {
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var c Config
	json.Unmarshal(raw, &c)
	log.Printf("Read configuration, found %d jobs\n", len(c.SiteList))

	return c.SiteList
}

func WriteConfig(jobs Config) {
	data, _ := json.Marshal(jobs)

	err := ioutil.WriteFile("./config2.json", data, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func LogSetupAndDestruct() func() {
	logFile, err := os.OpenFile("monitor.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Panicln(err)
	}

	log.SetOutput(io.MultiWriter(os.Stderr, logFile))

	return func() {
		e := logFile.Close()
		if e != nil {
			fmt.Fprintf(os.Stderr, "Problem closing the log file: %s\n", e)
		}
	}
}
