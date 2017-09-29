package main

import "log"

var config Config

func main() {
	// setup logging
	defer LogSetupAndDestruct()()

	// let's read config file
	config = ReadConfig()

	if Initialize(config.SiteList) != true {
		log.Fatal("Failed to initailze app, invalid configuration.")
	}

	// let's start few workers, ideally this should be set via config.
	// TODO: frequency is managed by worker itself, instead it should be managed by schedular
	for _, job := range monitoringJobs {
		log.Println("Starting monitor for url : ", job.URL)
		// we'll start a goroutine with id and channels
		go Monitor(job)
	}

	// read result here.
	StartServer()
}
