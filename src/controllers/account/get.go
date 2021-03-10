package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// AdminList contreoller
func AdminList(c *gin.Context) {
	var response models.Response

	admins := new(models.Admins)
	if err := admins.Get(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessDataList(admins.Data()))
}

// AdminInformation controller
func AdminInformation(c *gin.Context) {
	var response models.Response

	claims := middlewares.GetClaims(c)

	admin := new(models.Admin)
	found, _ := admin.GetByID(c, claims.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(admin))
}
