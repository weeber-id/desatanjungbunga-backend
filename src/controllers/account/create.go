package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// AdminCreate controller
func AdminCreate(c *gin.Context) {
	var request struct {
		Name     string `json:"name" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     int    `json:"role" binding:"required"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	newAdmin := &models.Admin{
		Name:     request.Name,
		Username: request.Username,
		Role:     request.Role,
	}
	newAdmin.SetPassword(request.Password)

	if err := newAdmin.Create(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "admin account created",
		"data":    newAdmin,
	})
}
