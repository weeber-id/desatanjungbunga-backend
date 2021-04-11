package travel

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetOne wisata controller
func GetOne(c *gin.Context) {
	var (
		wg      sync.WaitGroup
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

	travel := new(models.Travel)

	var found bool
	if request.ID != nil {
		found, _ = travel.GetByID(c, *request.ID)
	}
	if request.Slug != nil {
		found, _ = travel.GetBySlug(c, *request.Slug)
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		travel.WithCulinaryDetails(c)
	}()
	go func() {
		defer wg.Done()
		travel.WithLodgingDetails(c)
	}()
	wg.Wait()

	c.JSON(http.StatusOK, response.SuccessData(travel))
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

	multiTravel := new(models.MultipleWisata)
	if request.Search != nil {
		multiTravel.FilterBySearch(*request.Search)
	}

	claims := middlewares.GetClaims(c)
	if claims.ID == "" {
		multiTravel.FilterOnlyActive()
	}

	if claims.Role != 0 {
		multiTravel.FilterByAuthorID(claims.ID)
	}
	multiTravel.SortByRecommendation()
	if request.SortName != nil {
		multiTravel.SortByName(*request.SortName)
	}
	if request.SortDate != nil {
		multiTravel.SortByDate(*request.SortDate)
	} else {
		multiTravel.SortByDate("desc")
	}
	if request.Page != nil && request.ContentPerPage != nil {
		multiTravel.FilterByPaginate(*request.Page, *request.ContentPerPage)
	}

	maxPage := make(chan uint)
	defer close(maxPage)
	go func() {
		maxPage <- multiTravel.CountMaxPage(c)
	}()

	if err := multiTravel.Get(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	responseData.Data = multiTravel.Data()
	responseData.MaxPage = <-maxPage

	c.JSON(http.StatusOK, response.SuccessDataList(responseData))
}
