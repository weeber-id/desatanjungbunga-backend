package kuliner

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

	kuliner := new(models.Kuliner)
	found, _ := kuliner.GetByID(c, requestQuery.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "data not found"})
		return
	}

	kuliner.Title = requestBody.Title
	kuliner.ImageCover = requestBody.ImageCover
	kuliner.Description = requestBody.Description

	if err := kuliner.Update(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "kuliner updated",
		"data":    kuliner,
	})

}
