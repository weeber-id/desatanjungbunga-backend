package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
)

func AdminCheckUsernameIsExists(c *gin.Context) {
	var (
		requestQuery struct {
			Username string `form:"username" binding:"required"`
		}
		responseData struct {
			Status string
		}
		response models.Response
	)

	if err := c.BindQuery(&requestQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	admin := new(models.Admin)
	found, _ := admin.GetByUsername(c, requestQuery.Username)
	if found {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest("username exist"))
		return
	}

	responseData.Status = "username dapat digunakan"
	c.JSON(http.StatusOK, response.SuccessData(responseData))
}

// AdminCreate controller
func AdminCreate(c *gin.Context) {
	var (
		request struct {
			Name                string `json:"name" binding:"required"`
			Email               string `json:"email" binding:"required"`
			Address             string `json:"address" binding:"required"`
			DateofBirth         string `json:"date_of_birth" binding:"required"`
			PhoneNumberWhatsapp string `json:"phone_number_whatsapp" binding:"required"`
			Username            string `json:"username" binding:"required"`
			Password            string `json:"password" binding:"required"`
			Role                int    `json:"role" binding:"required"`
			ProfilePicture      string `json:"profile_picture" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// ====================== Check role =========================
	claims := middlewares.GetClaims(c)
	if claims.Role != 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorForbidden())
		return
	}

	// ================ Check is username exist ? ================
	admin := new(models.Admin)
	found, _ := admin.GetByUsername(c, request.Username)
	if found {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest("username exist"))
		return
	}

	// ================= Create new admin account ================
	newAdmin := &models.Admin{
		Name:                request.Name,
		Email:               request.Email,
		Address:             request.Address,
		DateofBirth:         request.DateofBirth,
		PhoneNumberWhatsapp: request.PhoneNumberWhatsapp,
		Username:            request.Username,
		Role:                request.Role,
		ProfilePicture:      request.ProfilePicture,
	}
	newAdmin.SetPassword(request.Password)

	if err := newAdmin.Create(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessDataCreated(newAdmin))
}
