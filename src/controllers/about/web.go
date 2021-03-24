package about

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

func Get(c *gin.Context) {
	var response models.Response

	about := new(models.About)
	about.Get(c)

	c.JSON(http.StatusOK, response.SuccessData(about))
}
