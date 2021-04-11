package account

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
	"github.com/weeber-id/desatanjungbunga-backend/src/services"
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
			Name                string `json:"name"`
			Email               string `json:"email"`
			Address             string `json:"address"`
			DateofBirth         string `json:"date_of_birth"`
			PhoneNumberWhatsapp string `json:"phone_number_whatsapp"`
			Username            string `json:"username"`
			Password            string `json:"password"`
			Role                int    `json:"role"`
			ProfilePicture      string `json:"profile_picture"`
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

	// Send welcome email through email
	email := services.Email{To: request.Email}
	if err := email.SendWelcomeAccount(request.Name, request.Username, request.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, response.ErrorBadRequest(fmt.Sprintf("error in email smtp: %s", err.Error())))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessDataCreated(newAdmin))
}
