package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func executeCommand(s ...string) ([]byte, error) {
	// Gen command and Remove duplicate whitespace from a command
	space := regexp.MustCompile(`\s+`)
	command := space.ReplaceAllString(fmt.Sprint(strings.Join(s[:], " ")), " ")
	if len(command) == 0 {
		return nil, nil
	}
	log.Printf("_run: %s", command)

	str, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		log.Print(err)
		outStr := strings.Trim(string(str), "\n")
		return nil, errors.New(outStr)
	}
	return str, err
}

func checkDockerStatus() (bool, error) {
	statusCommand := `/usr/bin/docker inspect -f '{{.State.Running}}' $(/usr/bin/docker ps --filter "label=app=prometheus" --format "{{.ID}}")`
	output, err := executeCommand(statusCommand)
	if err != nil {
		return false, err
	}

	status := strings.TrimSpace(string(output))
	return status == "true", nil
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
