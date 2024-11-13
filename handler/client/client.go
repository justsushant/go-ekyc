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

	accessToken, err := h.service.GenerateAccessToken(payload, client.AccessTokenExpiry, []byte(TEMP_SECRET))
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	refreshToken, err := h.service.GenerateRefreshToken(payload, client.RefreshTokenExpiry, []byte(TEMP_SECRET))
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"accessKey": accessToken,
		"secretKey": refreshToken,
	})
}
