package controller

import (
	"github.com/applichic/lynou/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

const codeFileNotFound = "CODE_FILE_NOT_FOUND"
const codeFileNotUploaded = "CODE_FILE_NOT_UPLOADED"

type StorageController struct {
}

// Download specific file respecting roles
func (s *StorageController) DownloadImage(c *gin.Context) {
	path := c.Param("path")
	result, err := util.DownloadObject(path)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to download the file",
			"code":  codeFileNotFound,
		})
		return
	}

	c.Data(http.StatusOK, "image/png", result)
}

// Upload file in the object storage
func (s *StorageController) UploadFile(c *gin.Context) {
	name := c.Param("name")

	// Open the file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeFileNotUploaded,
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeFileNotUploaded,
		})
		return
	}

	content := make([]byte, file.Size)
	_, err = f.Read(content)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeFileNotUploaded,
		})
		return
	}

	// Upload the file
	var contentType = file.Header["Content-Type"][0]
	err = util.UploadObject(name, &content, contentType)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to upload the file",
			"code":  codeFileNotUploaded,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}
