package travel

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// Create controller
func Create(c *gin.Context) {
	var (
		request struct {
			Name          string `json:"name" binding:"required"`
			Image         string `json:"image" binding:"required"`
			Price         string `json:"price" binding:"required"`
			OperationTime struct {
				From struct {
					Day  string `json:"day" binding:"required"`
					Time string `json:"time" binding:"required"`
				} `json:"from" binding:"required"`
				To struct {
					Day  string `json:"day" binding:"required"`
					Time string `json:"time" binding:"required"`
				} `json:"to" binding:"required"`
			} `json:"operation_time" binding:"required"`
			ShortDescription string `json:"short_description" binding:"required"`
			Description      string `json:"description" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	travel := &models.Travel{
		Name:  request.Name,
		Image: request.Image,
		Price: request.Price,
		OperationTime: struct {
			From struct {
				Day  string "bson:\"day\" json:\"day\""
				Time string "bson:\"time\" json:\"time\""
			} "bson:\"from\" json:\"from\""
			To struct {
				Day  string "bson:\"day\" json:\"day\""
				Time string "bson:\"time\" json:\"time\""
			} "bson:\"to\" json:\"to\""
		}(request.OperationTime),
		ShortDescription: request.ShortDescription,
		Description:      request.Description,
	}

	if err := travel.Create(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessDataCreated(travel))
}
