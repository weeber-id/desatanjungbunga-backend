package wisata

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// Update controller
func Update(c *gin.Context) {
	var (
		requestQuery struct {
			ID string `form:"id" binding:"required"`
		}
		requestBody struct {
			Title       string `json:"title" binding:"required"`
			ImageCover  string `json:"image_cover" binding:"required"`
			Description string `json:"description" binding:"required"`
		}
	)

	if err := c.BindQuery(&requestQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	wisata := new(models.Wisata)
	found, _ := wisata.GetByID(c, requestQuery.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "data not found"})
		return
	}

	wisata.Title = requestBody.Title
	wisata.ImageCover = requestBody.ImageCover
	wisata.Description = requestBody.Description

	if err := wisata.Update(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "wisata updated",
		"data":    wisata,
	})

}
