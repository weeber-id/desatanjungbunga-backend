package article

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
		requestBody struct {
			Title      string `json:"title" binding:"required"`
			ImageCover string `json:"image_cover" binding:"required"`
			Body       string `json:"body" binding:"required"`
			Relateds   []struct {
				Source string `json:"source"`
				ID     string `json:"id"`
			} `json:"relateds"`
		}
		response models.Response
	)

	if err := c.BindQuery(&requestQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	article := new(models.Article)
	found, _ := article.GetByID(c, requestQuery.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	article.Title = requestBody.Title
	article.ImageCover = requestBody.ImageCover
	article.Body = requestBody.Body
	article.ResetRelateds()
	for _, related := range requestBody.Relateds {
		article.SetRelatedRow(related.Source, related.ID)
	}

	if err := article.Update(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	article.WithRelated(c)

	c.JSON(http.StatusOK, response.SuccessData(article))

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

	article := new(models.Article)
	found, _ := article.GetByID(c, request.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	article.Active = *request.Active

	if err := article.Update(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(article))
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

	claims := middlewares.GetClaims(c)
	if claims.Role != 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorForbidden())
		return
	}

	article := new(models.Article)
	found, _ := article.GetByID(c, request.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	article.Recommendation = *request.Recommendation

	if err := article.Update(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(article))
}
