package handcraft

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetOne belanja controller
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

	handcraft := new(models.Handcraft)

	var found bool
	if request.ID != nil {
		found, _ = handcraft.GetByID(c, *request.ID)
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(handcraft))
}

// GetMultiple controller
func GetMultiple(c *gin.Context) {
	var response models.Response

	multiHandcraft := new(models.MultipleBelanja)

	if err := multiHandcraft.Get(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessDataList(multiHandcraft.Data()))
}
