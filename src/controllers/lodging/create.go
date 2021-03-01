package lodging

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// Create controller
func Create(c *gin.Context) {
	var (
		request struct {
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

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	lodging := &models.Lodging{
		Name:  request.Name,
		Image: request.Image,

		Price: struct {
			Value string "bson:\"value\" json:\"value\""
			Unit  string "bson:\"unit\" json:\"unit\""
		}(request.Price),

		OperationTime: request.OperationTime,

		Links: []struct {
			Name string "bson:\"name\" json:\"name\""
			Link string "bson:\"link\" json:\"link\""
		}(request.Links),

		FacilitiesID:     request.FacilitiesID,
		ShortDescription: request.ShortDescription,
		Description:      request.Description,
	}

	if err := lodging.Create(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessDataCreated(lodging))
}

// CreateFacilities controller
func CreateFacilities(c *gin.Context) {
	var (
		request struct {
			Name string `json:"name" binding:"required"`
			Icon string `json:"icon" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	facilities := &models.LodgingFacility{
		Name: request.Name,
		Icon: request.Icon,
	}

	if err := facilities.Create(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer((err)))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessDataCreated(facilities))
}
