package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func jsonResponse(c *gin.Context, status string, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// SetUpRouter sets up the Gin router for testing.
func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
