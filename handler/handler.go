package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type Handler struct {
	service service.ControllerInterface
}

func NewHandler(service service.ControllerInterface) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/signup", h.SignupHandler)
}

func (h *Handler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	router.POST("/upload", h.FileUploadHandler)
}

func (h *Handler) SignupHandler(c *gin.Context) {
	var payload types.SignupPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	if err := h.service.ValidatePayload(payload); err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	keyPair, err := h.service.GenerateKeyPair()
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	err = h.service.SaveSignupData(payload, keyPair)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	accessKey, secretKey := keyPair.GetKeysPrivate()

	c.JSON(200, gin.H{
		"accessKey": accessKey,
		"secretKey": secretKey,
	})
}

func (h *Handler) FileUploadHandler(c *gin.Context) {
	// reading type from request body
	fileType := c.PostForm("type")

	// reading file from request body
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	// applying validations on file
	err = h.service.ValidateFile(fileHeader.Filename, fileType)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	// save the file to bucket
	err = h.service.SaveUploadedFile(fileHeader)
	if err != nil {
		c.JSON(500, gin.H{"errorMessage": err.Error()})
		return
	}

	// save data in psql

	// replace this with proper uuid
	c.JSON(200, gin.H{
		"message": "uploaded",
	})
}
