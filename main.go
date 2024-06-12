package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// Config represents the structure of the Prometheus configuration file.
type Config struct {
	ScrapeConfigs []ScrapeConfig `yaml:"scrape_configs"`
}

// ScrapeConfig represents a scrape configuration in the Prometheus configuration.
type ScrapeConfig struct {
	JobName        string              `yaml:"job_name"`
	StaticConfigs  []StaticConfig      `yaml:"static_configs"`
	ScrapeInterval string              `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout  string              `yaml:"scrape_timeout,omitempty"`
	MetricsPath    string              `yaml:"metrics_path,omitempty"`
	Params         map[string][]string `yaml:"params,omitempty"`
	RelabelConfigs []RelabelConfig     `yaml:"relabel_configs,omitempty"`
}

// StaticConfig represents the static configuration for targets in a scrape configuration.
type StaticConfig struct {
	Targets []string `yaml:"targets"`
}

// RelabelConfig represents the relabel configuration in a scrape configuration.
type RelabelConfig struct {
	SourceLabels []string `yaml:"source_labels,omitempty"`
	TargetLabel  string   `yaml:"target_label,omitempty"`
	Replacement  string   `yaml:"replacement,omitempty"`
}

// Job represents the simplified structure to be returned by the API, containing only the job name and static configurations.
type Job struct {
	JobName       string         `json:"job_name"`
	StaticConfigs []StaticConfig `json:"static_configs"`
}

var config Config

func main() {
	router := gin.Default()

	// Define the endpoint to list jobs
	router.GET("/jobs", listJobs)

	// Load the Prometheus configuration from the YAML file
	loadConfig()

	// Start the HTTP server on port 8080
	router.Run(":8080")
}

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

// listJobs handles the GET /jobs request and returns the filtered list of jobs in JSON format.
func listJobs(c *gin.Context) {
	var filteredJobs []Job
	for i, job := range config.ScrapeConfigs {
		// Skip the first job and any jobs with additional configurations
		if i == 0 || job.ScrapeInterval != "" || job.ScrapeTimeout != "" || job.MetricsPath != "" || len(job.Params) > 0 || len(job.RelabelConfigs) > 0 {
			continue
		}
		// Append the filtered job to the list, including only job_name and static_configs
		filteredJobs = append(filteredJobs, Job{
			JobName:       job.JobName,
			StaticConfigs: job.StaticConfigs,
		})
	}

	// Return the filtered jobs as a JSON response
	c.JSON(http.StatusOK, filteredJobs)
}
