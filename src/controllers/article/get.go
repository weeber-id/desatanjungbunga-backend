package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetOne article controller
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

	c.BindQuery(&request)

	articles := new(models.Articles)

	if request.Search != nil {
		articles.FilterBySearch(*request.Search)
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
