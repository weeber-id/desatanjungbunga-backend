package discussion

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// Create controller
func Create(c *gin.Context) {
	var (
		request struct {
			ParentID *string `json:"question_id"`

			ContentName string `json:"content_name" binding:"required"`
			ContentID   string `json:"content_id" binding:"required"`

			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required"`
			Body  string `json:"body" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	discussion := &models.Discussion{
		Name:  request.Name,
		Email: request.Email,
		Body:  request.Body,
	}

	if err := discussion.SetContentNameAndID(request.ContentName, request.ContentID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest("invalid content name and content id"))
		return
	}

	if request.ParentID != nil {
		found, err := discussion.SetParentID(*request.ParentID)
		if !found {
			c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
			return
		}
	}

	if err := discussion.Create(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessDataCreated(discussion))
}
