package lodging

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// Update controller
func Update(c *gin.Context) {
	var (
		requestQuery struct {
			ID string `form:"id" binding:"required"`
		}
		requestBody requestCreateUpdateLodging
		response    models.Response
	)

	if err := c.BindQuery(&requestQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	lodging := new(models.Lodging)
	found, _ := lodging.GetByID(c, requestQuery.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	requestBody.WriteToModel(lodging)

	if err := lodging.Update(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}
	lodging.LoadFacilitiesDetail(c)

	c.JSON(http.StatusOK, response.SuccessData(lodging))
}

// ChangeActiveDeactive content controller
func ChangeActiveDeactive(c *gin.Context) {
	var (
		request struct {
			ID     string `json:"id" binding:"required"`
			Active *bool  `json:"active" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	claims := middlewares.GetClaims(c)
	if claims.Role != 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorForbidden())
		return
	}

	lodging := new(models.Lodging)
	found, _ := lodging.GetByID(c, request.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	lodging.Active = *request.Active

	if err := lodging.Update(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(lodging))
}

// ChangeRecommendation content controller
func ChangeRecommendation(c *gin.Context) {
	var (
		request struct {
			ID             string `json:"id" binding:"required"`
			Recommendation *bool  `json:"recommendation" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	lodging := new(models.Lodging)
	found, _ := lodging.GetByID(c, request.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	lodging.Recommendation = *request.Recommendation

	if err := lodging.Update(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(lodging))
}
