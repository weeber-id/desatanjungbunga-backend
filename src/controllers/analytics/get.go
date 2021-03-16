package analytics

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// ContentCount controller for analytics
func ContentCount(c *gin.Context) {
	var (
		wg           sync.WaitGroup
		responseData responseContentCount
		response     models.Response
	)

	claims := middlewares.GetClaims(c)

	articles := new(models.Articles)
	culinary := new(models.MultipleKuliner)
	handcraft := new(models.MultipleBelanja)
	lodging := new(models.MultipleLodging)
	travel := new(models.MultipleWisata)

	if claims.Role != 0 {
		articles.FilterByAuthorID(claims.ID)
		culinary.FilterByAuthorID(claims.ID)
		handcraft.FilterByAuthorID(claims.ID)
		lodging.FilterByAuthorID(claims.ID)
		travel.FilterByAuthorID(claims.ID)
	}

	wg.Add(5)
	go func() {
		defer wg.Done()
		responseData.Article.Count = articles.CountDocuments(c)
	}()

	go func() {
		defer wg.Done()
		responseData.Culinary.Count = culinary.CountDocuments(c)
	}()

	go func() {
		defer wg.Done()
		responseData.Handcraft.Count = handcraft.CountDocuments(c)
	}()

	go func() {
		defer wg.Done()
		responseData.Lodging.Count = lodging.CountDocuments(c)
	}()

	go func() {
		defer wg.Done()
		responseData.Travel.Count = travel.CountDocuments(c)
	}()

	wg.Wait()

	c.JSON(http.StatusOK, response.SuccessData(responseData))
}
