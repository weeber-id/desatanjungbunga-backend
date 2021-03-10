package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
	"github.com/weeber-id/desatanjungbunga-backend/src/storages"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
)

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
	if match := admin.IsPasswordMatch(request.OldPassword); match == false {
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
