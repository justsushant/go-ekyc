package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/controller"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type Handler struct {
	service controller.ControllerInterface
}

func NewHandler(service controller.ControllerInterface) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/signup", h.SignupHandler)
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

	tokenPair, err := h.service.GenerateTokenPair(payload)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	err = h.service.SaveSignupData(payload, tokenPair.RefreshToken)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"accessKey": tokenPair.AccessToken,
		"secretKey": tokenPair.RefreshToken,
	})
}

func (h *Handler) FileUploadHandler(c *gin.Context) {
	// reading type from request body
	fileType := c.PostForm("type")

	// reading file from request body
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	// applying validations on file
	err = h.service.ValidateFile(file.Filename, fileType)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}
}