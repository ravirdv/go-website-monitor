package main

type Job struct {
	Frequency   int64  `json:"frequency"`
	URL         string `json:"url"`
	CheckString string `json:"check_text"`
}

type Jobs []Job

type Config struct {
	JobsTest Jobs `json:"site_list"`
}
