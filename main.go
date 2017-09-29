package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// let's read config file
	jobs := ReadConfig()
	log.Printf("Read configuration, found %d jobs\n", len(jobs))

	// workers will post back results on this channel
	resultChannel := make(chan string, 100)

	// let's start few workers, ideally this should be set via config.
	for _, job := range jobs {
		log.Println("Starting monitor for url : ", job.URL)
		// we'll start a goroutine with id and channels
		go Monitor(job, resultChannel)
	}

	// read result here.
	go StartServer()
	// let's wait for response
	listenForShutdown()
}

func listenForShutdown() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	log.Println("Application Started: Press Ctrl + C to exit.")
	sig := <-gracefulStop
	fmt.Printf("caught sig: %+v", sig)
	fmt.Println("Waiting for workers to shutdown")
	//time.Sleep( 2 * time.Second)
	os.Exit(0)
}
