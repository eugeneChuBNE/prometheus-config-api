package main

import (
	"net/http"
	"os"
	"strings"

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

	// Remove existing job with the same job name or IP address
	for i := len(config.ScrapeConfigs) - 1; i >= 0; i-- {
		job := config.ScrapeConfigs[i]
		if job.JobName == newJobRequest.JobName {
			config.ScrapeConfigs = append(config.ScrapeConfigs[:i], config.ScrapeConfigs[i+1:]...)
			continue
		}
		for _, staticConfig := range job.StaticConfigs {
			for _, target := range staticConfig.Targets {
				if target == newJobRequest.IPAddress+":26" || target == newJobRequest.IPAddress+":27" {
					config.ScrapeConfigs = append(config.ScrapeConfigs[:i], config.ScrapeConfigs[i+1:]...)
					break
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

	if os.Getenv("TEST_MODE") != "true" {
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

		if !isActive {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Docker container is not active", "data": nil})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Job added successfully", "data": newScrapeConfig})
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

			if os.Getenv("TEST_MODE") != "true" {
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

				if !isActive {
					c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Docker container is not active", "data": nil})
					return
				}
			}

			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Job removed successfully", "data": jobName})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Job not found", "data": nil})
}

func searchJobByIP(c *gin.Context) {
	ipAddress := c.Query("ip")
	if ipAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "IP address is required", "data": nil})
		return
	}

	if err := loadConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to load configuration", "data": nil})
		return
	}

	var foundJobs []Job
	for _, job := range config.ScrapeConfigs {
		for _, staticConfig := range job.StaticConfigs {
			for _, target := range staticConfig.Targets {
				if strings.Contains(target, ipAddress) {
					foundJobs = append(foundJobs, Job{
						JobName:       job.JobName,
						StaticConfigs: job.StaticConfigs,
					})
					break
				}
			}
		}
	}

	if len(foundJobs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "No job found with the specified IP address", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Job(s) retrieved successfully", "data": foundJobs})
}
