package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func executeCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT not set in .env file")
	}

	router := gin.Default()

	// Define the endpoints
	router.GET("/jobs", listJobs)
	router.POST("/jobs", addJob)
	router.DELETE("/jobs/:job_name", removeJob)

	// Load the Prometheus configuration from the YAML file
	loadConfig()

	router.Run(":" + port)
}
