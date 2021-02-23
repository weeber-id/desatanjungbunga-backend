package kuliner

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetOne kuliner controller
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

	kuliner := new(models.Culinary)

	var found bool
	if request.ID != nil {
		found, _ = kuliner.GetByID(c, *request.ID)
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(kuliner))
}

// GetMultiple controller
func GetMultiple(c *gin.Context) {
	var response models.Response

	multiKuliner := new(models.MultipleKuliner)

	if err := multiKuliner.Get(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessDataList(multiKuliner.Data()))
}
