package wisata

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// Create controller
func Create(c *gin.Context) {
	var request struct {
		Title       string `json:"title" binding:"required"`
		ImageCover  string `json:"image_cover" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	wisata := &models.Wisata{
		Title:       request.Title,
		ImageCover:  request.ImageCover,
		Description: request.Description,
	}

	if err := wisata.Create(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "wisata created",
		"data":    wisata,
	})
}
