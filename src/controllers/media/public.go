package media

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
	"github.com/weeber-id/desatanjungbunga-backend/src/storages"
)

// UploadPublicFile to minio storages and give the public URL
func UploadPublicFile(c *gin.Context) {
	var (
		req struct {
			FolderName string `form:"folder_name" binding:"required"`
		}
		response models.Response
	)

	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorBadRequest(err.Error()))
		return
	}

	// new public object
	newObject := new(storages.PublicObject)
	newObject.LoadFromFileHeaderRandomName(fileHeader, req.FolderName)
	newObject.Upload(c)

	c.JSON(http.StatusOK, &models.Response{
		Meta: models.Meta{
			Message: "file telah diupload",
			Status:  "success",
			Code:    200,
		},
		Data: gin.H{
			"url": newObject.URL,
		},
	})
}
