package controller

import (
	"github.com/applichic/lynou/config"
	"github.com/applichic/lynou/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

const codeFileNotFound = "CODE_FILE_NOT_FOUND"
const codeFileNotUploaded = "CODE_FILE_NOT_UPLOADED"

type StorageController struct {
}

// Download Video file
func (s *StorageController) DownloadVideoFile(c *gin.Context) {
	name := c.Param("name")
	key := c.DefaultQuery("key", "")

	// Check if the key is sent
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to download the video",
			"code":  codeFileNotFound,
		})
		return
	}

	// Parse the token
	token, err := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.JwtSecret), nil
	})

	// Check if the token is correct and valid
	if err != nil || token == nil || !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to download the video",
			"code":  codeFileNotFound,
		})
		return
	}

	header, result, err := util.DownloadObject(name)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to download the video",
			"code":  codeFileNotFound,
		})
		return
	}

	c.Header("Last-Modified", header.Get("Last-Modified"))
	c.Header("Etag", header.Get("Etag"))
	c.Header("X-Trans-Id", header.Get("X-Trans-Id"))
	c.Header("Content-Length", header.Get("Content-Length"))
	c.Header("Accept-Ranges", header.Get("Accept-Ranges"))
	c.Header("X-Timestamp", header.Get("X-Timestamp"))
	c.Data(http.StatusOK, header.Get("Content-Type"), result)
}

// Download specific file respecting roles
func (s *StorageController) DownloadFile(c *gin.Context) {
	name := c.Param("name")
	header, result, err := util.DownloadObject(name)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to download the file",
			"code":  codeFileNotFound,
		})
		return
	}

	c.Data(http.StatusOK, header.Get("Content-Type"), result)
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
