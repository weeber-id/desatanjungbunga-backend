package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetOne article controller
func GetOne(c *gin.Context) {
	var request struct {
		ID   *string `form:"id"`
		Slug *string `form:"slug"`
	}

	c.BindQuery(&request)
	if request.ID == nil && request.Slug == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "id or slug must be filled",
		})
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
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    article,
	})
}

// GetMultiple controller
func GetMultiple(c *gin.Context) {
	var request struct {
		SortTitle      *string `form:"sort_title"`
		SortDate       *string `form:"sort_date"`
		Page           *int64  `form:"page"`
		ContentPerPage *int64  `form:"content_per_page"`
	}

	c.BindQuery(&request)

	articles := new(models.Articles)

	if request.SortTitle != nil {
		articles.SortByTitle(*request.SortTitle)
	}
	if request.SortDate != nil {
		articles.SortByDate(*request.SortDate)
	}
	if request.Page != nil && request.ContentPerPage != nil {
		articles.SetPagination(*request.Page, *request.ContentPerPage)
	}

	if err := articles.Get(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    articles.Data(),
	})
}
