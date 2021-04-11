package article

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
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

	if err := c.BindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

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

	articles := new(models.Articles)

	claims := middlewares.GetClaims(c)
	if claims.ID == "" {
		articles.FilterOnlyActive()
	}

	if request.Search != nil {
		articles.FilterBySearch(*request.Search)
	}
	if claims.Role != 0 {
		articles.FilterByAuthorID(claims.ID)
	}
	articles.SortByRecommendation()
	if request.SortTitle != nil {
		articles.SortByTitle(*request.SortTitle)
	}
	if request.SortDate != nil {
		articles.SortByDate(*request.SortDate)
	} else {
		articles.SortByDate("desc")
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

// GetSearchInline controller
// used for admin get list all content when choose content related to article
func GetSearchInline(c *gin.Context) {
	type responseRowData struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Source string `json:"source"`
	}

	var (
		wg      sync.WaitGroup
		request struct {
			Search         *string `form:"search"`
			SortTitle      *string `form:"sort_title"`
			SortDate       *string `form:"sort_date"`
			Page           *int    `form:"page"`
			ContentPerPage *int    `form:"content_per_page"`
		}
		responseData struct {
			Data []responseRowData `json:"data"`
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

	if request.Search != nil {
		articles.FilterBySearch(*request.Search)
		culinaries.FilterBySearch(*request.Search)
		handcrafts.FilterBySearch(*request.Search)
		lodgings.FilterBySearch(*request.Search)
		travels.FilterBySearch(*request.Search)
	}

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

	var (
		articlesData  []models.Article
		culinaryData  []models.Culinary
		handcraftData []models.Handcraft
		lodgingData   []models.Lodging
		travelData    []models.Travel
	)

	wg.Add(5)
	go func() {
		defer wg.Done()
		articles.Get(c)
		articlesData = articles.Data()
	}()

	go func() {
		defer wg.Done()
		culinaries.Get(c)
		culinaryData = culinaries.Data()
	}()

	go func() {
		defer wg.Done()
		handcrafts.Get(c)
		handcraftData = handcrafts.Data()
	}()

	go func() {
		defer wg.Done()
		lodgings.Get(c)
		lodgingData = lodgings.Data()
	}()

	go func() {
		defer wg.Done()
		travels.Get(c)
		travelData = travels.Data()
	}()

	wg.Wait()

	for _, row := range articlesData {
		responseData.Data = append(responseData.Data, responseRowData{
			ID:     row.ID.Hex(),
			Name:   row.Title,
			Source: variables.Collection.Article,
		})
	}
	for _, row := range culinaryData {
		responseData.Data = append(responseData.Data, responseRowData{
			ID:     row.ID.Hex(),
			Name:   row.Name,
			Source: variables.Collection.Culinary,
		})
	}
	for _, row := range handcraftData {
		responseData.Data = append(responseData.Data, responseRowData{
			ID:     row.ID.Hex(),
			Name:   row.Name,
			Source: variables.Collection.Handcraft,
		})
	}
	for _, row := range lodgingData {
		responseData.Data = append(responseData.Data, responseRowData{
			ID:     row.ID.Hex(),
			Name:   row.Name,
			Source: variables.Collection.Lodging,
		})
	}
	for _, row := range travelData {
		responseData.Data = append(responseData.Data, responseRowData{
			ID:     row.ID.Hex(),
			Name:   row.Name,
			Source: variables.Collection.Travel,
		})
	}

	c.JSON(http.StatusOK, response.SuccessData(responseData))
}
