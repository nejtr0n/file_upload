package http

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"test/storage/proto"
)

type router struct {
	*gin.Engine
	client proto.StorageService
}

type List []File

type File struct {
	Name string `json:"name" binding:"required"`
	Size int `json:"size" binding:"required"`
	Type string `json:"type" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func NewRouter(client proto.StorageService) *router {
	r := new(router)

	r.client = client

	r.Engine = gin.Default()
	r.Engine.LoadHTMLGlob("storage/ui/http/templates/*.html") // ToDo: move to separate repo


	api := r.Engine.Group("/storage")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Test form
	api.GET("/test", func(c *gin.Context) {
		c.HTML(200, "test.html", struct {}{})
	})

	// Upload file route
	api.POST("/upload", func(c *gin.Context) {
		// Get multipart form
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		files := form.File["images[]"]

		for _, file := range files {
			// Get file body
			reader, err := file.Open()
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("could not open file to read: %s", err.Error()))
				return
			}
			content, err := ioutil.ReadAll(reader)
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("could not read file: %s", err.Error()))
				return
			}

			// Get file mimetype
			mime := http.DetectContentType(content)

			// Create store file rpc request
			request := &proto.File{
				Name: filepath.Base(file.Filename),
				Size: file.Size,
				Content: content,
				Type: mime,
			}

			// Send rpc request to save file
			_, err = r.client.Save(context.Background(), request)
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("could not send rpc request: %s", err.Error()))
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "Success",
		})
	})

	api.POST("/upload/link", func(c *gin.Context) {
		url := c.PostForm("url")
		file, err := http.Get(url)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("error downloading file: %s", err.Error()))
			return
		}
		defer file.Body.Close()

		content, err := ioutil.ReadAll(file.Body)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("could not get file body: %s", err.Error()))
			return
		}
		// Get file mimetype
		mime := http.DetectContentType(content)

		// Create store file rpc request
		request := &proto.File{
			Name: path.Base(url),
			Size: int64(len(content)),
			Content: content,
			Type: mime,
		}

		// Send rpc request to save file
		resp, err := r.client.Save(context.Background(), request)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("could not send rpc request: %s", err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": resp.Msg,
		})
	})

	// Get json array
	api.POST("/upload/json", func(c *gin.Context) {
		// Unpack json to list of structures
		data := new(List)
		err := c.BindJSON(data)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("could not unmarshal files: %s", err.Error()))
			return
		}
		// Save each file
		for _, file := range *data {
			// Get base64 value
			b64data := file.Content[strings.IndexByte(file.Content, ',')+1:]
			data, err := base64.StdEncoding.DecodeString(b64data)
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("could not decode base64 file string: %s", err.Error()))
				return
			}

			// Create store file rpc request
			request := &proto.File{
				Name: file.Name,
				Size: int64(file.Size),
				Content: data,
				Type: file.Type,
			}

			// Send rpc request to save file
			_, err = r.client.Save(context.Background(), request)
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("could not send rpc request: %s", err.Error()))
				return
			}

		}

		c.JSON(http.StatusOK, gin.H{
			"status": "Success",
		})
	})

	return r
}

