# Simple Web Monitor
 - This is a WIP Simple Web Monitor, it comes with a UI to keep track of monitoring activity.
 - REST APIs to add/delete/modify monitoring jobs.
 - Built using golang, gorilla websockets & vuejs for frontend.
 - Currently it supports Status Code and String check verification.


## Installation

    `go run *.go`

## Usage
 * Add Host
    - `curl -H "Content-Type: application/json" -X POST -d '{"url":"<url>", "frequency" : 15 }' http://localhost:8000/monitor/job/add`
 * Delete Host
    - `curl -H "Content-Type: application/json" -X POST -d '{"url":"<url>", "frequency" : 15 }' http://localhost:8000/monitor/job/remove`

## TODO:
- Distributed agents for monitoring via multiple locations
- Authentication & User Profiles
- More response verification conditions.
