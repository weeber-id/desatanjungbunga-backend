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

	kuliner := new(models.Culinary)

	var found bool
	if request.ID != nil {
		found, _ = kuliner.GetByID(c, *request.ID)
	}
	if request.Slug != nil {
		found, _ = kuliner.GetBySlug(c, *request.Slug)
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(kuliner))
}

// GetMultiple controller
func GetMultiple(c *gin.Context) {
	var (
		request struct {
			SortName       *string `form:"sort_title"`
			SortDate       *string `form:"sort_date"`
			Page           *int    `form:"page"`
			ContentPerPage *int    `form:"content_per_page"`
		}
		response models.Response
	)

	c.BindQuery(&request)

	multiKuliner := new(models.MultipleKuliner)

	if request.SortName != nil {
		multiKuliner.SortByName(*request.SortName)
	}
	if request.SortDate != nil {
		multiKuliner.SortByDate(*request.SortDate)
	}
	if request.Page != nil && request.ContentPerPage != nil {
		multiKuliner.FilterByPaginate(*request.Page, *request.ContentPerPage)
	}

	if err := multiKuliner.Get(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessDataList(multiKuliner.Data()))
}
