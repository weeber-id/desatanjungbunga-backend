package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

// AdminLogin controller
func AdminLogin(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	admin := new(models.Admin)
	found, _ := admin.GetByUsername(c, request.Username)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "data not found"})
		return
	}

	if !admin.IsPasswordMatch(request.Password) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	middlewares.WriteAccessToken2Cookie(c, admin.ID.Hex(), admin.Role)

	c.JSON(http.StatusOK, models.Response{
		Meta: models.Meta{
			Message: "Login Success",
			Status:  "OK",
			Code:    200,
		},
		Data: admin,
	})
}

// AdminLogut controller
func AdminLogut(c *gin.Context) {
	middlewares.DeleteAccessToken2Cookie(c)

	c.JSON(http.StatusOK, gin.H{"message": "admin logout"})
}
