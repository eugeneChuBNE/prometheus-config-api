package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/jobs", listJobs)
	router.POST("/jobs", addJob)
	router.DELETE("/jobs/:job_name", removeJob)
	router.GET("/jobs/search", searchJobByIP)
	return router
}

func TestListJobs(t *testing.T) {
	os.Setenv("TEST_MODE", "true")

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "success", response.Status)
}

func TestAddJob(t *testing.T) {
	os.Setenv("TEST_MODE", "true")

	job := AddJobRequest{
		JobName:   "test_job",
		IPAddress: "192.168.1.1",
	}
	jsonValue, _ := json.Marshal(job)

	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "success", response.Status)
}

func TestRemoveJob(t *testing.T) {
	os.Setenv("TEST_MODE", "true")

	job := AddJobRequest{
		JobName:   "existing_job",
		IPAddress: "192.168.1.1",
	}
	jsonValue, _ := json.Marshal(job)

	router := setupRouter()

	// Add the job first
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Now remove the job
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/jobs/existing_job", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "success", response.Status)
}

func TestSearchJobByIP(t *testing.T) {
	os.Setenv("TEST_MODE", "true")

	job := AddJobRequest{
		JobName:   "test_search_job",
		IPAddress: "192.168.1.1",
	}
	jsonValue, _ := json.Marshal(job)

	router := setupRouter()

	// Add the job first
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Now search for the job by IP
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/jobs/search?ip=192.168.1.1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "success", response.Status)
}
