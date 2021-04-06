package account

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
	"github.com/weeber-id/desatanjungbunga-backend/src/services"
	"github.com/weeber-id/desatanjungbunga-backend/src/storages"
	"github.com/weeber-id/desatanjungbunga-backend/src/tools"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
)

// AdminUpdateSellerAccount super admin update seller account
func AdminUpdateSellerAccount(c *gin.Context) {
	var (
		requestQuery struct {
			UserID string `form:"user_id" binding:"required"`
		}
		requestBody struct {
			Name                string `json:"name" binding:"required"`
			Address             string `json:"address" binding:"required"`
			DateofBirth         string `json:"date_of_birth" binding:"required"`
			PhoneNumberWhatsapp string `json:"phone_number_whatsapp" binding:"required"`
			Email               string `json:"email" binding:"required"`
			ProfilePicture      string `json:"profile_picture"`
		}
		response models.Response
	)

	if err := c.BindQuery(&requestQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// ====================== Check role =========================
	claims := middlewares.GetClaims(c)
	if claims.Role != 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorForbidden())
		return
	}

	admin := new(models.Admin)
	found, _ := admin.GetByID(c, requestQuery.UserID)
	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorDataNotFound())
		return
	}

	// ================= Create new admin account ================
	admin.Name = requestBody.Name
	admin.Email = requestBody.Email
	admin.Address = requestBody.Address
	admin.DateofBirth = requestBody.DateofBirth
	admin.PhoneNumberWhatsapp = requestBody.PhoneNumberWhatsapp
	admin.ProfilePicture = requestBody.ProfilePicture

	if err := admin.Update(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessDataCreated(admin))
}

// ========================================================================================

func AdminUpdateSellerResetPassword(c *gin.Context) {
	var (
		requestQuery struct {
			UserID string `form:"user_id" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindQuery(&requestQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// Check role, only super admin to seller
	claims := middlewares.GetClaims(c)
	if claims.Role != 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorForbidden())
		return
	}

	// find target seller
	seller := new(models.Admin)
	found, _ := seller.GetByID(c, requestQuery.UserID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}
	if seller.Role == 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorForbidden())
		return
	}

	// Generate new password
	newPassword := tools.RandStringRunes(10)

	// Send generated password through email
	email := services.Email{To: seller.Email}
	if err := email.SendNewPasswordForReset(seller.Name, seller.Username, newPassword); err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, response.ErrorBadRequest(fmt.Sprintf("error in email smtp: %s", err.Error())))
		return
	}

	// Update seller account with new password
	seller.SetPassword(newPassword)
	if err := seller.Update(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(nil))
}

// ========================================================================================

func AdminUpdate(c *gin.Context) {
	var (
		requestBody struct {
			Name                string `json:"name" binding:"required"`
			Address             string `json:"address" binding:"required"`
			DateofBirth         string `json:"date_of_birth" binding:"required"`
			PhoneNumberWhatsapp string `json:"phone_number_whatsapp" binding:"required"`
			Email               string `json:"email" binding:"required"`
			ProfilePicture      string `json:"profile_picture"`
		}
		response models.Response
	)

	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	claims := middlewares.GetClaims(c)

	admin := new(models.Admin)
	found, _ := admin.GetByID(c, claims.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorDataNotFound())
		return
	}

	// ================= Create new admin account ================
	admin.Name = requestBody.Name
	admin.Email = requestBody.Email
	admin.Address = requestBody.Address
	admin.DateofBirth = requestBody.DateofBirth
	admin.PhoneNumberWhatsapp = requestBody.PhoneNumberWhatsapp
	admin.ProfilePicture = requestBody.ProfilePicture

	if err := admin.Update(c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessDataCreated(admin))
}

// AdminChangePassword controller
func AdminChangePassword(c *gin.Context) {
	var (
		request struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}
		response models.Response
	)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	claims := middlewares.GetClaims(c)

	admin := new(models.Admin)
	found, _ := admin.GetByID(c, claims.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	// =================== Check old password ===================
	if match := admin.IsPasswordMatch(request.OldPassword); !match {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest("wrong old password"))
		return
	}

	admin.SetPassword(request.NewPassword)
	if err := admin.Update(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(nil))
}

// AdminChangeProfilePicture controller
func AdminChangeProfilePicture(c *gin.Context) {
	var response models.Response

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	claims := middlewares.GetClaims(c)

	admin := new(models.Admin)
	found, _ := admin.GetByID(c, claims.ID)
	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, response.ErrorDataNotFound())
		return
	}

	newObject := new(storages.PublicObject)
	newObject.LoadFromFileHeaderRandomName(fileHeader, variables.FolderName.ProfilePicture)
	_, err = newObject.Upload(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	admin.ProfilePicture = newObject.URL
	if err := admin.Update(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(nil))
}

// AdminDeleteProfilePicture controller
func AdminDeleteProfilePicture(c *gin.Context) {
	var response models.Response

	admin := middlewares.GetAdmin(c)
	admin.ProfilePicture = ""

	if err := admin.Update(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorInternalServer(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessData(admin))
}
