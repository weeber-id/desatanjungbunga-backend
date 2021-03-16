package discussion

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetMultiple controller
func GetMultiple(c *gin.Context) {
	var (
		request struct {
			ContentName    string  `form:"content_name" binding:"required"`
			ContentID      string  `form:"content_id" binding:"required"`
			ParentID       *string `form:"question_id"`
			ShowAnswer     bool    `form:"show_answer"`
			SortDate       *string `form:"sort_date"`
			Page           *int    `form:"page"`
			ContentPerPage *int    `form:"content_per_page"`
		}
		response models.Response
	)

	if err := c.BindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	discussion := new(models.Discussions)

	if err := discussion.FilterByContentNameID(request.ContentName, request.ContentID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest("invalid content name and id"))
		return
	}

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

	if err := discussion.Get(c, request.ShowAnswer); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessDataList(discussion.Data()))
}
