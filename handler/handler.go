package handler

import (
	"encoding/json"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	router.POST("/face-match", h.FaceMatchHandler)
	router.POST("/ocr", h.OCRHandler)
	router.POST("/face-match-async", h.FaceMatchHandlerAsync)
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

	// fetching client_id from request scoped variables
	clientID, ok := c.Get("client_id")
	if !ok {
		// TODO: what to do when ok is false, or clientID is nil
	}

	objectName := uuid.NewString()
	uploadMetaData := &types.UploadMetaData{
		Type:       fileType,
		ClientID:   clientID.(int),
		FilePath:   strconv.Itoa(clientID.(int)) + "/" + objectName + filepath.Ext(fileHeader.Filename), // filepath is saved like, clientID/uuid.extension
		FileSizeKB: fileHeader.Size,
	}

	// save the file to bucket and psql
	err = h.service.SaveFile(fileHeader, uploadMetaData)
	if err != nil {
		c.JSON(500, gin.H{"errorMessage": err.Error()})
		return
	}

	// replace this with proper uuid
	c.JSON(200, gin.H{
		"id": objectName,
	})
}

func (h *Handler) FaceMatchHandler(c *gin.Context) {
	var payload types.FaceMatchPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	// fetching client_id from request scoped variables
	clientID, ok := c.Get("client_id")
	if !ok {
		// TODO: what to do when ok is false, or clientID is nil
	}

	if err := h.service.ValidateImage(payload, clientID.(int)); err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	score, err := h.service.CalcAndSaveFaceMatchScore(payload, clientID.(int))
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"score": score,
	})
}

func (h *Handler) OCRHandler(c *gin.Context) {
	var payload types.OCRPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	// generating UUID for file name
	clientID, ok := c.Get("client_id")
	if !ok {
		// TODO: fetch clientID here using
	}

	if err := h.service.ValidateImageOCR(payload, clientID.(int)); err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	resp, err := h.service.PerformAndSaveOCR(payload, clientID.(int))
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	c.JSON(200, resp)
}

func (h *Handler) FaceMatchHandlerAsync(c *gin.Context) {
	var payload types.FaceMatchPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	// fetching client_id from request scoped variables
	clientID, ok := c.Get("client_id")
	if !ok {
		// TODO: what to do when ok is false, or clientID is nil
	}

	id, err := h.service.PerformFaceMatchAsync(payload, clientID.(int))
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"id": id,
	})
}

func (h *Handler) OCRHandlerAsync(c *gin.Context) {
	var payload types.OCRPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	// generating UUID for file name
	clientID, ok := c.Get("client_id")
	if !ok {
		// TODO: fetch clientID here using
	}

	id, err := h.service.PerformOCRAsync(payload, clientID.(int))
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"id": id,
	})
}
