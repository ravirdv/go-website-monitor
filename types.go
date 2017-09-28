package main

type Job struct {
	Frequency          int64   `json:"frequency"`
	URL                string  `json:"url"`
	CheckString        string  `json:"check_text"`
	LastChecked        int64   `json:"last_checked,omitempty"`
	ResponseTime       float64 `json:"response_time,omitempty"`
	CheckStringPresent bool    `json:"check_string_present,omitempty"`
	StatusCode         int     `json:"status_code,omitempty"`
}

type Jobs []Job

type Config struct {
	JobsTest Jobs `json:"site_list"`
}
