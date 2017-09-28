package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func ReadConfig() Jobs {
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	log.Printf("%s", raw)

	var c Config
	json.Unmarshal(raw, &c)
	return c.JobsTest
}

func WriteConfig(jobs Config) {
	data, _ := json.Marshal(jobs)

	err := ioutil.WriteFile("./config2.json", data, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
