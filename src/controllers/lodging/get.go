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
	lodging.LoadFacilitiesDetail(c)

	c.JSON(http.StatusOK, response.SuccessData(lodging))
}

// GetMultiple lodging controller
func GetMultiple(c *gin.Context) {
	var (
		request struct {
			Search         *string `form:"search"`
			SortName       *string `form:"sort_title"`
			SortDate       *string `form:"sort_date"`
			Page           *int    `form:"page"`
			ContentPerPage *int    `form:"content_per_page"`
		}
		responseData struct {
			Data    interface{} `json:"data"`
			MaxPage uint        `json:"max_page"`
		}
		response models.Response
	)

	c.BindQuery(&request)

	lodgings := new(models.MultipleLodging)

	if request.Search != nil {
		lodgings.FilterBySearch(*request.Search)
	}
	if request.SortName != nil {
		lodgings.SortByName(*request.SortName)
	}
	if request.SortDate != nil {
		lodgings.SortByDate(*request.SortDate)
	}
	if request.Page != nil && request.ContentPerPage != nil {
		lodgings.FilterByPaginate(*request.Page, *request.ContentPerPage)
	}

	if err := lodgings.Get(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	responseData.Data = lodgings.Data()

	c.JSON(http.StatusOK, response.SuccessDataList(responseData))
}

// GetFacilities controller
func GetFacilities(c *gin.Context) {
	var (
		response models.Response
	)

	facilities := new(models.MultipleLodgingFacility)

	if err := facilities.Get(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessDataList(facilities.Data()))
}
