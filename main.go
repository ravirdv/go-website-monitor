package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"os"
	"os/signal"
	"syscall"
	"strings"
)

func main() {

	// channel for passing monitoring jobs to workers.
	jobChannel := make(chan string, 1000)
	// workers will postback results on this channel
	resultChannel:= make(chan string, 100)
    // let's start few workers, ideally this should be set via config.
	for workerId := 1; workerId <= 10; workerId++ {
		// we'll start a goroutine with id and channels
		go worker(workerId, jobChannel, resultChannel)
	}

	// move this to scheduler, which should be a goroutine.
	for i := 0; i < 2; i++ {
		url := ""
		jobChannel <- url
	}
	// let's wait for response
	listenForShutdown()
}

// Monitor website : monitors given website with specified frequency
func Monitor(url string, frequency int, containsText string) {
	// set timeout to 5 seconds, we don't want people to abuse this system.
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
    	Timeout: timeout,
	}
	// let's keep track of start time.
	start := time.Now()
	resp, err := client.Get(url)
	// let's get total time elasped for this operation.
	secs := time.Since(start).Seconds() *1e3 
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyString  := string(body)
		if containsText != "" {
			if strings.Contains(bodyString, containsText) {
				log.Println("response containss text %s", containsText)
			}
		}
		log.Printf("%.2f elapsed with response length: %d %s\n", secs, len(body), url)
	} else{
		log.Fatal(err)
	}
}

// workers will take a monitering job and result response.
func worker(id int, jobs <-chan string, results chan<- string) {
	for j := range jobs {
		log.Println("worker", id, "started  job", j)
		Monitor("http://google.com", 1, "")
		//results <- j
		log.Println("worker", id, "finished job", j)
	}
}

func listenForShutdown(){
	var gracefulStop = make(chan os.Signal)
    signal.Notify(gracefulStop, syscall.SIGTERM)
    signal.Notify(gracefulStop, syscall.SIGINT)
	log.Println("Application Started: Press Ctrl + C to exit.")
	sig := <-gracefulStop
	fmt.Println("caught sig: %+v", sig)
	fmt.Println("Waiting for workers to shutdown")
	time.Sleep(2*time.Second)
	os.Exit(0)
}
