package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

// Config represents the structure of the Prometheus configuration file.
type Config struct {
	ScrapeConfigs []ScrapeConfig `yaml:"scrape_configs"`
}

var config Config

// init loads the environment variables from the .env file.
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}
}

// loadConfig reads the Prometheus configuration from the file specified in the environment variable and unmarshals it into the config variable.
func loadConfig() error {
	configPath := os.Getenv("PROMETHEUS_CONFIG_PATH")
	if configPath == "" {
		log.Printf("PROMETHEUS_CONFIG_PATH not set in .env file")
		return errors.New("PROMETHEUS_CONFIG_PATH not set in .env file")
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	return nil
}

// saveConfig writes the current configuration to the file specified in the environment variable.
func saveConfig() error {
	configPath := os.Getenv("PROMETHEUS_CONFIG_PATH")
	if configPath == "" {
		log.Printf("PROMETHEUS_CONFIG_PATH not set in .env file")
		return errors.New("PROMETHEUS_CONFIG_PATH not set in .env file")
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	return nil
}
