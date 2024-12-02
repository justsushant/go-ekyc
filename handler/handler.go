package handler

import (
	"encoding/json"
	"log"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type Handler struct {
	service service.ServiceManager
}

func NewHandler(service service.ServiceManager) Handler {
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
	router.POST("/ocr-async", h.OCRHandlerAsync)
	router.GET("/result/:jobType/:jobID", h.ResultHandler)
}

// @Summary Signup
// @Description Signups the client
// @Tags Signup
// @Accept json
// @Produce json
// @Param name body string true "Name of client"
// @Param email body string true "Email of client"
// @Param plan body string true "Name of plan"
// @Success 200 {object} types.SignupResponse "Access & secret keys"
// @Failure 400 {object} types.ErrorResponse "invalid email"
// @Failure 400 {object} types.ErrorResponse "invalid plan, supported plans are basic, advanced, or enterprise"
// @Router /api/v1/signup [post]
func (h *Handler) SignupHandler(c *gin.Context) {
	var payload types.SignupPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(400, types.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	keyPair, err := h.service.SignupClient(payload)
	if err != nil {
		c.JSON(400, types.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	accessKey, secretKey := keyPair.GetKeysPrivate()

	c.JSON(200, types.SignupResponse{
		AccessKey: accessKey,
		SecretKey: secretKey,
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
		FileSizeKB: fileHeader.Size / 1000,
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
		// TODO: what to do when ok is false, or clientID is nil
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

	// fetch data from cache
	jobID, ok := h.service.FetchDataFromCache(payload, clientID.(int), types.FACE_MATCH_WORK_TYPE)
	if ok {
		c.JSON(200, gin.H{
			"id": jobID,
		})
		return
	}

	jobID, err = h.service.PerformFaceMatchAsync(payload, clientID.(int))
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	// set data in cache
	h.service.SetDataInCache(payload, clientID.(int), types.FACE_MATCH_WORK_TYPE, jobID)

	c.JSON(200, gin.H{
		"id": jobID,
	})
}

func (h *Handler) OCRHandlerAsync(c *gin.Context) {
	var payload types.OCRPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	clientID, ok := c.Get("client_id")
	if !ok {
		// TODO: what to do when ok is false, or clientID is nil
	}

	// fetch data from cache
	jobID, ok := h.service.FetchDataFromCache(payload, clientID.(int), types.OCR_WORK_TYPE)
	if ok {
		c.JSON(200, gin.H{
			"id": jobID,
		})
		return
	}

	jobID, err = h.service.PerformOCRAsync(payload, clientID.(int))
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}
	log.Println(jobID)

	// set data in cache
	h.service.SetDataInCache(payload, clientID.(int), types.OCR_WORK_TYPE, jobID)

	c.JSON(200, gin.H{
		"id": jobID,
	})
}

func (h *Handler) ResultHandler(c *gin.Context) {
	jobID := c.Param("jobID")
	jobType := c.Param("jobType")
	clientID, ok := c.Get("client_id")
	if !ok {
		// TODO: what to do when ok is false, or clientID is nil
	}

	data, err := h.service.GetJobDetailsByJobID(jobID, jobType)
	if err != nil {
		log.Println("Error while fetching job details by job id: ", err)
		c.JSON(500, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	// validate if client id of the job and client is same
	if data.ClientID != clientID.(int) {
		c.JSON(400, gin.H{
			"errorMessage": service.ErrInvalidJobId.Error(),
		})
		return
	}

	// filter on the basis of status
	switch data.Status {
	case types.JOB_STATUS_PROCESSING:
		c.JSON(200, gin.H{
			"status":       data.Status,
			"message":      "job is still running",
			"processed_at": data.ProcessedAt,
		})
		return
	case types.JOB_STATUS_CREATED:
		c.JSON(200, gin.H{
			"status":     data.Status,
			"message":    "job is created",
			"created_at": data.CreatedAt,
		})
		return
	case types.JOB_STATUS_FAILED:
		c.JSON(200, gin.H{
			"status":        data.Status,
			"message":       "job is failed",
			"failed_at":     data.FailedAt,
			"failed_reason": data.FailedReason,
		})
		return
	case types.JOB_STATUS_COMPLETED:
		switch data.Type {
		case types.FACE_MATCH_WORK_TYPE:
			c.JSON(200, gin.H{
				"status":       data.Status,
				"message":      "job is completed",
				"completed_at": data.CompletedAt,
				"result":       data.MatchScore,
			})
			return
		case types.OCR_WORK_TYPE:
			c.JSON(200, gin.H{
				"status":       data.Status,
				"message":      "job is completed",
				"completed_at": data.CompletedAt,
				"result":       data.OCRDetails,
			})
			return
		}
	}

	c.JSON(500, gin.H{
		"errorMessage": "Unexpected server error occurred",
	})
}
