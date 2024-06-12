package main

import (
	"github.com/gin-gonic/gin"
)

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
