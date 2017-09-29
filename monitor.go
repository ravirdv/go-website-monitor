package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

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

func processResponse() {

}
