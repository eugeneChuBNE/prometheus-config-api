package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
