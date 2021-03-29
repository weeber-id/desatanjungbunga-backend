package article

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetOne article controller
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

	article := new(models.Article)

	var found bool
	if request.ID != nil {
		found, _ = article.GetByID(c, *request.ID)
	}
	if request.Slug != nil {
		found, _ = article.GetBySlug(c, *request.Slug)
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		article.WithAuthor(c)
	}()
	go func() {
		defer wg.Done()
		article.WithRelated(c)
	}()
	wg.Wait()

	c.JSON(http.StatusOK, response.SuccessData(article))
}

// GetMultiple controller
func GetMultiple(c *gin.Context) {
	var (
		request struct {
			Search         *string `form:"search"`
			SortTitle      *string `form:"sort_title"`
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

	claims := middlewares.GetClaims(c)
	articles := new(models.Articles)

	if request.Search != nil {
		articles.FilterBySearch(*request.Search)
	}
	if claims.Role != 0 {
		articles.FilterByAuthorID(claims.ID)
	}
	if request.SortTitle != nil {
		articles.SortByTitle(*request.SortTitle)
	}
	if request.SortDate != nil {
		articles.SortByDate(*request.SortDate)
	}

	if request.Page != nil && request.ContentPerPage != nil {
		articles.FilterByPaginate(*request.Page, *request.ContentPerPage)
	}

	maxPage := make(chan uint)
	defer close(maxPage)
	go func() {
		maxPage <- articles.CountMaxPage(c)
	}()

	if err := articles.Get(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	responseData.Data = articles.Data()
	responseData.MaxPage = <-maxPage

	c.JSON(http.StatusOK, response.SuccessDataList(responseData))
}
