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

// TODO: add scheduler and fixed workers
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
	// let's run till shutdown is not requested.
	for !job.ShutDownRequest {
		// refresh job config
		job = monitoringJobs[job.URL]

		// set timeout, we don't want people to abuse this system.
		timeout := time.Duration(time.Duration(config.RequestTimeOut) * time.Second) //TODO: read via config
		client := http.Client{
			Timeout: timeout,
		}
		// let's keep track of start time.
		start := time.Now()
		resp, err := client.Get(job.URL)
		// let's get total time elasped for this operation.
		secs := time.Since(start).Seconds() * 1e3 // we want in milliseconds
		if err == nil {
			// read response body, TODO: limit this to max configurable size
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
			broadcast <- job // TODO: instead of directly sending to client, this should be sent in batches.
		} else {
			// oops, failed to send get request
			log.Println("WARN: Error proccessing job with url: ", job.URL)
			job.Result.Type = "Error"
			job.Result.Details = err.Error()
			log.Println(err.Error())
		}
		frequency := job.Frequency
		if frequency < 5 {
			// too aggresive, resetting to 5
			log.Printf("WARN: frequency(%d) too low for job %s, resetting to 5 seconds", job.Frequency, job.URL)
			frequency = 5
		}

		time.Sleep(time.Duration(time.Second * time.Duration(frequency)))
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
			job.CheckStringPresent = true
			job.Result.Type = fmt.Sprintf("Response contains : %s", job.CheckString)
		} else {
			job.Result.Type = "Error"
			job.Result.Details = job.Result.Details + "\n" + fmt.Sprintf("Response content doesn't contain string %s", job.CheckString)
		}
	}
}
