package travel

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// Create controller
func Create(c *gin.Context) {
	var (
		wg       sync.WaitGroup
		request  requestCreateUpdateTravel
		response models.Response
	)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	admin := middlewares.GetAdmin(c)

	travel := new(models.Travel)
	request.WriteToModel(travel)

	if err := travel.Create(c, admin); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
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

	c.JSON(http.StatusCreated, response.SuccessDataCreated(travel))
}
