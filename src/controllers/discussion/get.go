package discussion

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetMultiple controller
func GetMultiple(c *gin.Context) {
	var request struct {
		ParentID       *string `form:"question_id"`
		ShowQuestion   bool    `form:"show_question"`
		SortDate       *string `form:"sort_date"`
		Page           *int    `form:"page"`
		ContentPerPage *int    `form:"content_per_page"`
	}

	c.BindQuery(&request)

	discussion := new(models.Discussions)

	if request.SortDate != nil {
		discussion.SortByDate(*request.SortDate)
	}

	if request.ParentID != nil {
		discussion.FilterOnlyAnswer(*request.ParentID)
	} else {
		discussion.FilterOnlyQuestion()
	}

	if request.Page != nil && request.ContentPerPage != nil {
		discussion.FilterByPaginate(*request.Page, *request.ContentPerPage)
	}

	if err := discussion.Get(c, request.ShowQuestion); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    discussion.Data(),
	})
}
