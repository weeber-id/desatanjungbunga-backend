package lodging

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetOne lodging controller
func GetOne(c *gin.Context) {
	var (
		request struct {
			ID   *string `form:"id"`
			Slug *string `form:"slug"`
		}
		response models.Response
	)

	c.BindQuery(&request)
	if request.ID == nil && request.Slug == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest("id atau slug harus diisi"))
		return
	}

	lodging := new(models.Lodging)

	var found bool
	if request.ID != nil {
		found, _ = lodging.GetByID(c, *request.ID)
	}
	if request.Slug != nil {
		found, _ = lodging.GetBySlug(c, *request.Slug)
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(lodging))
}

// GetMultiple lodging controller
func GetMultiple(c *gin.Context) {
	var response models.Response

	lodgings := new(models.MultipleLodging)

	if err := lodgings.Get(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessDataList(lodgings.Data()))
}
