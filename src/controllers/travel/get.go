package travel

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetOne wisata controller
func GetOne(c *gin.Context) {
	var (
		request struct {
			ID *string `form:"id" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	travel := new(models.Travel)

	var found bool
	if request.ID != nil {
		found, _ = travel.GetByID(c, *request.ID)
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(travel))
}

// GetMultiple controller
func GetMultiple(c *gin.Context) {
	var response models.Response

	multiTravel := new(models.MultipleWisata)

	if err := multiTravel.Get(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessDataList(multiTravel.Data()))
}
