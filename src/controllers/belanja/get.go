package belanja

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// GetOne belanja controller
func GetOne(c *gin.Context) {
	var request struct {
		ID *string `form:"id" binding:"required"`
	}

	if err := c.BindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	belanja := new(models.Belanja)

	var found bool
	if request.ID != nil {
		found, _ = belanja.GetByID(c, *request.ID)
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    belanja,
	})
}

// GetMultiple controller
func GetMultiple(c *gin.Context) {
	multiBelanja := new(models.MultipleBelanja)

	if err := multiBelanja.Get(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    multiBelanja.Data(),
	})
}
