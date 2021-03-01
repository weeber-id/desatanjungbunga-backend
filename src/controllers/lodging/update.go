package lodging

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
				Value string `json:"value" binding:"required"`
				Unit  string `json:"unit" binding:"required"`
			} `json:"price" binding:"required"`
			OperationTime string `json:"operation_time" binding:"required"`
			Links         []struct {
				Name string `json:"name" binding:"required"`
				Link string `json:"link" binding:"required"`
			} `json:"links" binding:"required"`
			FacilitiesID     []string `json:"facilities_id" binding:"required"`
			ShortDescription string   `json:"short_description" binding:"required"`
			Description      string   `json:"description" binding:"required"`
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

	lodging := new(models.Lodging)
	found, _ := lodging.GetByID(c, requestQuery.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	lodging.Name = requestBody.Name
	lodging.Image = requestBody.Image

	lodging.Price = struct {
		Value string "bson:\"value\" json:\"value\""
		Unit  string "bson:\"unit\" json:\"unit\""
	}(requestBody.Price)

	lodging.OperationTime = requestBody.OperationTime

	lodging.Links = []struct {
		Name string "bson:\"name\" json:\"name\""
		Link string "bson:\"link\" json:\"link\""
	}(requestBody.Links)

	lodging.FacilitiesID = requestBody.FacilitiesID
	lodging.ShortDescription = requestBody.ShortDescription
	lodging.Description = requestBody.Description

	if err := lodging.Update(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}
	lodging.LoadFacilitiesDetail(c)

	c.JSON(http.StatusOK, response.SuccessData(lodging))
}
