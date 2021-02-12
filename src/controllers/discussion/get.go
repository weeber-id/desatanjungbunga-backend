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
		SortDate       *string `form:"sort_date"`
		Page           *int64  `form:"page"`
		ContentPerPage *int64  `form:"content_per_page"`
	}

	c.BindQuery(&request)

	discussion := new(models.Discussions)

	if request.SortDate != nil {
		discussion.SortByDate(*request.SortDate)
	}
	if request.Page != nil && request.ContentPerPage != nil {
		discussion.SetPagination(*request.Page, *request.ContentPerPage)
	}

	if request.ParentID != nil {
		discussion.FilterOnlyAnswer(*request.ParentID)
	} else {
		discussion.FilterOnlyQuestion()
	}

	if err := discussion.Get(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    discussion.Data(),
	})
}
