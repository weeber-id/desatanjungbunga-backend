package discussion

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// Create controller
func Create(c *gin.Context) {
	var request struct {
		ParentID *string `json:"question_id"`
		Name     string  `json:"name" binding:"required"`
		Email    string  `json:"email" binding:"required"`
		Body     string  `json:"body" binding:"required"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	discussion := &models.Discussion{
		Name:  request.Name,
		Email: request.Email,
		Body:  request.Body,
	}

	if request.ParentID != nil {
		found, err := discussion.SetParentID(*request.ParentID)
		if !found {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "data not found"})
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "search question data"})
			return
		}
	}

	if err := discussion.Create(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Discussion created",
		"data":    discussion,
	})
}
