package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/controller/client"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

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
}
