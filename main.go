package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ScrapeConfigs []ScrapeConfig `yaml:"scrape_configs"`
}

type ScrapeConfig struct {
	JobName        string              `yaml:"job_name"`
	StaticConfigs  []StaticConfig      `yaml:"static_configs"`
	ScrapeInterval string              `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout  string              `yaml:"scrape_timeout,omitempty"`
	MetricsPath    string              `yaml:"metrics_path,omitempty"`
	Params         map[string][]string `yaml:"params,omitempty"`
	RelabelConfigs []RelabelConfig     `yaml:"relabel_configs,omitempty"`
}

type StaticConfig struct {
	Targets []string `yaml:"targets"`
}

type RelabelConfig struct {
	SourceLabels []string `yaml:"source_labels,omitempty"`
	TargetLabel  string   `yaml:"target_label,omitempty"`
	Replacement  string   `yaml:"replacement,omitempty"`
}

type Job struct {
	JobName       string         `json:"job_name"`
	StaticConfigs []StaticConfig `json:"static_configs"`
}

type AddJobRequest struct {
	JobName   string `json:"job_name"`
	IPAddress string `json:"ip_address"`
}

var config Config

func main() {
	router := gin.Default()

	// Define the endpoints
	router.GET("/jobs", listJobs)
	router.POST("/jobs", addJob)
	router.DELETE("/jobs/:job_name", removeJob)

	// Load the Prometheus configuration from the YAML file
	loadConfig()

	router.Run(":8080")
}

// read Prometheus configuration from prometheus.yml and unmarshals it into the config variable.
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

func addJob(c *gin.Context) {
	var newJobRequest AddJobRequest
	if err := c.BindJSON(&newJobRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the IP address already exists in any job
	for _, job := range config.ScrapeConfigs {
		for _, staticConfig := range job.StaticConfigs {
			for _, target := range staticConfig.Targets {
				if target == newJobRequest.IPAddress+":26" || target == newJobRequest.IPAddress+":27" {
					c.JSON(http.StatusConflict, gin.H{"error": "job existed"})
					return
				}
			}
		}
	}

	// Add new job
	newScrapeConfig := ScrapeConfig{
		JobName: newJobRequest.JobName,
		StaticConfigs: []StaticConfig{
			{
				Targets: []string{
					newJobRequest.IPAddress + ":26",
					newJobRequest.IPAddress + ":27",
				},
			},
		},
	}

	config.ScrapeConfigs = append(config.ScrapeConfigs, newScrapeConfig)
	saveConfig()
	c.JSON(http.StatusOK, gin.H{"status": "job added"})
}

func removeJob(c *gin.Context) {
	jobName := c.Param("job_name")
	for i, job := range config.ScrapeConfigs {
		if job.JobName == jobName {
			// Remove the job from the list
			config.ScrapeConfigs = append(config.ScrapeConfigs[:i], config.ScrapeConfigs[i+1:]...)
			saveConfig()
			c.JSON(http.StatusOK, gin.H{"status": "job removed"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
}
