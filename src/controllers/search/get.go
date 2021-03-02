package search

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetSearch controller
func GetSearch(c *gin.Context) {
	type responseEachCategory struct {
		Data    interface{} `json:"data"`
		MaxPage uint        `json:"max_page"`
	}

	var (
		wg      sync.WaitGroup
		request struct {
			Search         string  `form:"search" binding:"required"`
			SortTitle      *string `form:"sort_title"`
			SortDate       *string `form:"sort_date"`
			Page           *int    `form:"page"`
			ContentPerPage *int    `form:"content_per_page"`
		}
		responseData struct {
			Article   responseEachCategory `json:"articles"`
			Culinary  responseEachCategory `json:"culinaries"`
			Handcraft responseEachCategory `json:"handcrafts"`
			Lodging   responseEachCategory `json:"lodgings"`
			Travel    responseEachCategory `json:"travels"`
		}
		response models.Response
	)

	if err := c.BindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	articles := new(models.Articles)
	culinaries := new(models.MultipleKuliner)
	handcrafts := new(models.MultipleBelanja)
	lodgings := new(models.MultipleLodging)
	travels := new(models.MultipleWisata)

	articles.FilterBySearch(request.Search)
	culinaries.FilterBySearch(request.Search)
	handcrafts.FilterBySearch(request.Search)
	lodgings.FilterBySearch(request.Search)
	travels.FilterBySearch(request.Search)

	if request.SortTitle != nil {
		articles.SortByTitle(*request.SortTitle)
		culinaries.SortByName(*request.SortTitle)
		handcrafts.SortByName(*request.SortTitle)
		lodgings.SortByName(*request.SortTitle)
		travels.SortByName(*request.SortTitle)
	}
	if request.SortDate != nil {
		articles.SortByDate(*request.SortDate)
		culinaries.SortByDate(*request.SortDate)
		handcrafts.SortByDate(*request.SortDate)
		lodgings.SortByDate(*request.SortDate)
		travels.SortByDate(*request.SortDate)
	}
	if request.Page != nil && request.ContentPerPage != nil {
		articles.FilterByPaginate(*request.Page, *request.ContentPerPage)
		culinaries.FilterByPaginate(*request.Page, *request.ContentPerPage)
		handcrafts.FilterByPaginate(*request.Page, *request.ContentPerPage)
		lodgings.FilterByPaginate(*request.Page, *request.ContentPerPage)
		travels.FilterByPaginate(*request.Page, *request.ContentPerPage)
	}

	wg.Add(5)
	go func() {
		defer wg.Done()
		articles.Get(c)
		responseData.Article.Data = articles.Data()
	}()

	go func() {
		defer wg.Done()
		culinaries.Get(c)
		responseData.Culinary.Data = culinaries.Data()
	}()

	go func() {
		defer wg.Done()
		handcrafts.Get(c)
		responseData.Handcraft.Data = handcrafts.Data()
	}()

	go func() {
		defer wg.Done()
		lodgings.Get(c)
		responseData.Lodging.Data = lodgings.Data()
	}()

	go func() {
		defer wg.Done()
		travels.Get(c)
		responseData.Travel.Data = travels.Data()
	}()

	wg.Wait()

	c.JSON(http.StatusOK, response.SuccessData(responseData))
}
