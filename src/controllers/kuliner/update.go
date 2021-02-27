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
			Name  string `json:"name" binding:"required"`
			Image string `bson:"image" json:"image"`
			Price struct {
				Start string `json:"start" binding:"required"`
				End   string `json:"end" binding:"required"`
				Unit  string `json:"unit" binding:"required"`
			} `json:"price" binding:"required"`
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
			Links []struct {
				Name string `json:"name" binding:"required"`
				Link string `json:"link" binding:"required"`
			} `json:"links" binding:"required"`
			ShortDescription string `json:"short_description" binding:"required"`
			Description      string `json:"description" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindQuery(&requestQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	kuliner := new(models.Culinary)
	found, _ := kuliner.GetByID(c, requestQuery.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	kuliner.Name = requestBody.Name
	kuliner.Image = requestBody.Image

	kuliner.Price = struct {
		Start string "bson:\"start\" json:\"start\""
		End   string "bson:\"end\" json:\"end\""
		Unit  string "bson:\"unit\" json:\"unit\""
	}(requestBody.Price)

	kuliner.OperationTime = struct {
		From struct {
			Day  string "bson:\"day\" json:\"day\""
			Time string "bson:\"time\" json:\"time\""
		} "bson:\"from\" json:\"from\""
		To struct {
			Day  string "bson:\"day\" json:\"day\""
			Time string "bson:\"time\" json:\"time\""
		} "bson:\"to\" json:\"to\""
	}(requestBody.OperationTime)

	kuliner.Links = []struct {
		Name string "bson:\"name\" json:\"name\""
		Link string "bson:\"link\" json:\"link\""
	}(requestBody.Links)

	kuliner.ShortDescription = requestBody.ShortDescription
	kuliner.Description = requestBody.Description

	if err := kuliner.Update(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(kuliner))

}
