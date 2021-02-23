package lodging

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// Delete controller
func Delete(c *gin.Context) {
	var (
		requestQuery struct {
			ID string `form:"id" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindQuery(&requestQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	lodging := new(models.Lodging)
	found, _ := lodging.GetByID(c, requestQuery.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	if err := lodging.Delete(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(nil))
}
