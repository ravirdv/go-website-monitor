package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	// let's read config file
	jobs := ReadConfig()
	log.Printf("Read configuration, found %d jobs\n", len(jobs))

	// workers will postback results on this channel
	resultChannel := make(chan string, 100)

	// let's start few workers, ideally this should be set via config.
	for _, job := range jobs {
		log.Println("Starting monitor for url : ", job.URL)
		// we'll start a goroutine with id and channels
		go Monitor(job, resultChannel)
	}

	// read result here.

	// let's wait for response
	listenForShutdown()
}

// Monitor website : monitors given website with specified frequency
func Monitor(job Job, results chan<- string) {
	isShutDownRequested := false
	for !isShutDownRequested {
		// set timeout to 5 seconds, we don't want people to abuse this system.
		timeout := time.Duration(5 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		// let's keep track of start time.
		start := time.Now()
		resp, err := client.Get(job.URL)
		// let's get total time elasped for this operation.
		secs := time.Since(start).Seconds() * 1e3
		if err == nil {
			body, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(body)
			if job.CheckString != "" {
				if strings.Contains(bodyString, job.CheckString) {
					log.Printf("response containss text %s\n", job.CheckString)
				} else {
					// sucess but response doesn't contain required string.
					log.Printf("response doesn't contains text %s\n", job.CheckString)
				}
			}
			log.Printf("%.2f ms elapsed, statusCode:%d, Response length: %d %s\n", secs, resp.StatusCode, len(body), job.URL)
		} else {
			log.Fatal(err)
		}
		time.Sleep(time.Duration(time.Second * time.Duration(job.Frequency)))
	}
}

func listenForShutdown() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	log.Println("Application Started: Press Ctrl + C to exit.")
	sig := <-gracefulStop
	fmt.Printf("caught sig: %+v", sig)
	fmt.Println("Waiting for workers to shutdown")
	time.Sleep(2 * time.Second)
	os.Exit(0)
}
