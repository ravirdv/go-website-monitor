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

// Monitor website : monitors given website with specified frequency
func Monitor(job Job, results chan<- string) {
	isShutDownRequested := false
	// let's run till shutdown is not requested.
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
		secs := time.Since(start).Seconds() * 1e3 // we want in milliseconds
		if err == nil {
			// read response body
			body, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(body)
			if job.CheckString != "" {
				if strings.Contains(bodyString, job.CheckString) {
					//log.Printf("response containss text %s\n", job.CheckString)
					job.CheckStringPresent = true
				} else {
					// sucess but response doesn't contain required string.
					//log.Printf("response doesn't contains text %s\n", job.CheckString)
				}
			}
			log.Printf("Time Elapsed: %.2f ms, Status Code: %d, Response Length: %d %s\n", secs, resp.StatusCode, len(body), job.URL)
			job.ResponseTime = secs
			job.LastChecked = time.Now().Unix()
			job.StatusCode = resp.StatusCode
			broadcast <- job
		} else {
			log.Println(err.Error())
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
	//time.Sleep( 2 * time.Second)
	os.Exit(0)
}
