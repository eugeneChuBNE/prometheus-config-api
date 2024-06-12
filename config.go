package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config represents the structure of the Prometheus configuration file.
type Config struct {
	ScrapeConfigs []ScrapeConfig `yaml:"scrape_configs"`
}

var config Config

// loadConfig reads the Prometheus configuration from prometheus.yml and unmarshals it into the config variable.
func loadConfig() {
	data, err := ioutil.ReadFile("prometheus.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

// saveConfig writes the current configuration to prometheus.yml.
func saveConfig() {
	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = ioutil.WriteFile("prometheus.yml", data, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
