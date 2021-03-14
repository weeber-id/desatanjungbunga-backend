package lodging

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// Create controller
func Create(c *gin.Context) {
	var (
		request  requestCreateUpdateLodging
		response models.Response
	)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	admin := middlewares.GetAdmin(c)

	lodging := new(models.Lodging)
	request.WriteToModel(lodging)

	if err := lodging.Create(c, admin); err != nil {
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
