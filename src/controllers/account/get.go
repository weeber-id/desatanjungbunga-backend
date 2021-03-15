package account

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// AdminList contreoller
func AdminList(c *gin.Context) {
	var (
		wg      sync.WaitGroup
		request struct {
			Page           *int `form:"page"`
			ContentPerPage *int `form:"content_per_page"`
			Role           *int `form:"role"`
		}
		responseData responseAdminList
		response     models.Response
	)

	if err := c.BindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	admins := new(models.Admins)
	if request.Role != nil {
		admins.FilterByRole(*request.Role)
	}
	if request.Page != nil && request.ContentPerPage != nil {
		admins.FilterByPaginate(*request.Page, *request.ContentPerPage)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := admins.Get(c); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
			return
		}
		responseData.Data = admins.Data()
	}()

	go func() {
		defer wg.Done()
		responseData.MaxPage = admins.CountMaxPage(c)
	}()

	wg.Wait()

	c.JSON(http.StatusOK, response.SuccessDataList(responseData))
}

// AdminInformation controller
func AdminInformation(c *gin.Context) {
	var (
		request struct {
			ID *string `form:"id"`
		}
		response models.Response
	)

	if err := c.BindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	claims := middlewares.GetClaims(c)
	admin := new(models.Admin)

	searchID := ""
	if request.ID != nil {
		searchID = *request.ID
	} else {
		searchID = claims.ID
	}

	found, _ := admin.GetByID(c, searchID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(admin))
}
