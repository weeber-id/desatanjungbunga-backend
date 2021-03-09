package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// Create controller
func Create(c *gin.Context) {
	var (
		request struct {
			Title      string `json:"title" binding:"required"`
			ImageCover string `json:"image_cover" binding:"required"`
			Body       string `json:"body" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	admin := middlewares.GetAdmin(c)
	article := &models.Article{
		Title:      request.Title,
		ImageCover: request.ImageCover,
		Body:       request.Body,
	}

	if err := article.Create(c, admin); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessDataCreated(article))
}
