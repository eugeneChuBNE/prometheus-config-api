package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func listJobs(c *gin.Context) {
	if err := loadConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to load configuration", "data": nil})
		return
	}
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
	jsonResponse(c, "success", "Jobs retrieved successfully", filteredJobs)
}

func addJob(c *gin.Context) {
	if err := loadConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to load configuration", "data": nil})
		return
	}
	var newJobRequest AddJobRequest
	if err := c.BindJSON(&newJobRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body", "data": nil})
		return
	}

	// Check if the job name already exists
	for _, job := range config.ScrapeConfigs {
		if job.JobName == newJobRequest.JobName {
			c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Job name already exists", "data": nil})
			return
		}
	}

	// Check if the IP address already exists in any job
	for _, job := range config.ScrapeConfigs {
		for _, staticConfig := range job.StaticConfigs {
			for _, target := range staticConfig.Targets {
				if target == newJobRequest.IPAddress+":26" || target == newJobRequest.IPAddress+":27" {
					c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Job with this IP address already exists", "data": nil})
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
	if err := saveConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save configuration", "data": nil})
		return
	}

	_, err := executeCommand("/usr/bin/docker container restart prometheus_prometheus_1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error(), "data": nil})
		return
	}

	isActive, err := checkDockerStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error(), "data": nil})
		return
	}

	if isActive {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Job added successfully", "data": newScrapeConfig})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "docker die roy", "data": nil})
	}
}

func removeJob(c *gin.Context) {
	if err := loadConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to load configuration", "data": nil})
		return
	}
	jobName := c.Param("job_name")
	for i, job := range config.ScrapeConfigs {
		if job.JobName == jobName {
			// Remove the job from the list
			config.ScrapeConfigs = append(config.ScrapeConfigs[:i], config.ScrapeConfigs[i+1:]...)
			if err := saveConfig(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save configuration", "data": nil})
				return
			}

			_, err := executeCommand("docker container restart prometheus_prometheus_1")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error(), "data": nil})
				return
			}

			isActive, err := checkDockerStatus()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error(), "data": nil})
				return
			}

			if isActive {
				c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Job removed successfully", "data": jobName})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "docker die roy", "data": nil})
			}
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Job not found", "data": nil})
}
