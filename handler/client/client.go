package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/controller/client"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

// TODO: Use the secret from config type
const TEMP_SECRET = "xyz"

type ClientHandler struct {
	service client.ClientServiceInterface
}

func NewHandler(service client.ClientServiceInterface) ClientHandler {
	return ClientHandler{
		service: service,
	}
}

func (h *ClientHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/signup", h.SignupHandler)
}

func (h *ClientHandler) SignupHandler(c *gin.Context) {
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

	c.JSON(200, gin.H{
		"accessKey": tokenPair.AccessToken,
		"secretKey": tokenPair.RefreshToken,
	})
}
