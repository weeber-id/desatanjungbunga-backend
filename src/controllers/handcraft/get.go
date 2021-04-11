package handcraft

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetOne belanja controller
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

	handcraft := new(models.Handcraft)

	var found bool
	if request.ID != nil {
		found, _ = handcraft.GetByID(c, *request.ID)
	}
	if request.Slug != nil {
		found, _ = handcraft.GetBySlug(c, *request.Slug)
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(handcraft))
}

// GetMultiple controller
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

	if err := c.BindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	multiHandcraft := new(models.MultipleBelanja)

	claims := middlewares.GetClaims(c)
	if claims.ID == "" {
		multiHandcraft.FilterOnlyActive()
	}

	if request.Search != nil {
		multiHandcraft.FilterBySearch(*request.Search)
	}
	if claims.Role != 0 {
		multiHandcraft.FilterByAuthorID(claims.ID)
	}
	multiHandcraft.SortByRecommendation()
	if request.SortName != nil {
		multiHandcraft.SortByName(*request.SortName)
	}
	if request.SortDate != nil {
		multiHandcraft.SortByDate(*request.SortDate)
	} else {
		multiHandcraft.SortByDate("desc")
	}
	if request.Page != nil && request.ContentPerPage != nil {
		multiHandcraft.FilterByPaginate(*request.Page, *request.ContentPerPage)
	}

	maxPage := make(chan uint)
	defer close(maxPage)
	go func() {
		maxPage <- multiHandcraft.CountMaxPage(c)
	}()

	if err := multiHandcraft.Get(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	responseData.Data = multiHandcraft.Data()
	responseData.MaxPage = <-maxPage

	c.JSON(http.StatusOK, response.SuccessDataList(responseData))
}
