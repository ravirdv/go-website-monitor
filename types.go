package main

type Job struct {
	Frequency          int64   `json:"frequency"`
	URL                string  `json:"url"`
	CheckString        string  `json:"check_text,omitempty"`
	LastChecked        int64   `json:"last_checked,omitempty"`
	ResponseTime       float64 `json:"response_time,omitempty"`
	CheckStringPresent bool    `json:"check_string_present,omitempty"`
	ExpectedStatusCode int     `json:"expected_status_code,omitempty"`
	ActualStatusCode   int     `json:"actual_status_code,omitempty"`
	Result             Message `json:"result,omitempty"`
}

type Jobs []Job

type Config struct {
	SiteList       Jobs `json:"site_list"`
	RequestTimeOut int  `json:"request_timeout,omitempty"`
}

type Message struct {
	Type    string `json:"type,omitempty"`
	Details string `json:"details,omitempty"`
}
