package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// global list of monitoring jobs, should be maintained in K/V store like redis in case of distibuted setup.
var monitoringJobs = make(map[string]Job)
var schedulerTable = make(map[string]time.Time)

func Initialize(jobs Jobs) bool {
	for _, job := range jobs {
		// validate uri
		_, err := url.ParseRequestURI(job.URL)
		if err != nil {
			log.Print("WARN: Initialize: invalid uri : ", job.URL)
			return false
		}
		// check if host already exists
		if _, ok := monitoringJobs[job.URL]; ok {
			log.Print("WARN: Initialize: already exists with uri : ", job.URL)
			return true
		}
		monitoringJobs[job.URL] = job
	}
	return true
}

// Monitor website : monitors given website with specified frequency
func Monitor(job Job) {
	isShutDownRequested := false
	// let's run till shutdown is not requested.
	for !isShutDownRequested {
		if _, ok := monitoringJobs[job.URL]; !ok {
			isShutDownRequested = true
			log.Printf("Shutting down monitoring worker for url : %s", job.URL)
		}
		// refresh job config
		job = monitoringJobs[job.URL]
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
			processStringCheck(bodyString, &job)

			log.Printf("ResponseTime: %.2f ms, StatusCode: %d, ResponseLength: %d %s\n", secs, resp.StatusCode, len(body), job.URL)
			// set result values
			job.ResponseTime = secs
			job.LastChecked = time.Now().Unix()
			job.ActualStatusCode = resp.StatusCode
			processStatusCode(&job)
			// send it to all clients connected via websockets
			broadcast <- job
			//do something here
		} else {
			log.Println("WARN: Error proccessing job with url: ", job.URL)
			job.Result.Type = "Error"
			job.Result.Details = err.Error()
			log.Println(err.Error())
		}
		time.Sleep(time.Duration(time.Second * time.Duration(job.Frequency)))
	}
}

func processStatusCode(job *Job) {
	if job.ExpectedStatusCode == 0 {
		return
	}
	if job.ExpectedStatusCode == job.ActualStatusCode {
		job.Result.Type = "Success"
	} else {
		job.Result.Type = "Error"
		job.Result.Details = fmt.Sprintf("Status code doesn't match Expected:%d actual: %d", job.ExpectedStatusCode, job.ActualStatusCode)
	}
}

func processStringCheck(content string, job *Job) {
	if job.CheckString != "" {
		if strings.Contains(content, job.CheckString) {
			//log.Printf("response containss text %s\n", job.CheckString)
			job.CheckStringPresent = true
			job.Result.Type = fmt.Sprintf("Response contains : %s", job.CheckString)
		} else {
			job.Result.Type = "Error"
			job.Result.Details = job.Result.Details + "\n" + fmt.Sprintf("Response content doesn't contain string %s", job.CheckString)
		}
	}
}
