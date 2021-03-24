package about

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

func AdminGet(c *gin.Context) {
	var response models.Response

	claims := middlewares.GetClaims(c)
	if claims.Role != 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorForbidden())
		return
	}

	about := new(models.About)
	found, _ := about.Get(c)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(about))
}

// ============================================================================

func AdminUpdate(c *gin.Context) {
	var (
		requestBody requestAdminUpdate
		response    models.Response
	)

	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	claims := middlewares.GetClaims(c)
	if claims.Role != 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorForbidden())
		return
	}

	about := new(models.About)
	found, _ := about.Get(c)

	requestBody.Write2Model(about)
	if !found {
		if err := about.Create(c); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
			return
		}
	} else {
		if err := about.Update(c); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
			return
		}
	}

	c.JSON(http.StatusOK, response.SuccessData(about))
}

// ============================================================================

func AdminDelete(c *gin.Context) {
	var response models.Response

	claims := middlewares.GetClaims(c)
	if claims.Role != 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorForbidden())
		return
	}

	about := new(models.About)
	found, _ := about.Get(c)
	if !found {
		return
	}
	if err := about.Delete(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(nil))
}
