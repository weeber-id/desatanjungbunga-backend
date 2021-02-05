package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck controller
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "connect"})
}
