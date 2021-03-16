package discussion

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
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	discussion := new(models.Discussion)
	found, _ := discussion.GetByID(c, requestQuery.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	if err := discussion.Delete(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(nil))
}
