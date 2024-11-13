package handler

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

var ErrInvalidEmail = errors.New("invalid email")

func SignupHandler(c *gin.Context) {
	var signupPaylod types.SignupPayload
	err := json.NewDecoder(c.Request.Body).Decode(&signupPaylod)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	if !isEmailValid(signupPaylod.Email) {
		c.JSON(400, gin.H{"errorMessage": ErrInvalidEmail.Error()})
		return
	}
}

func isEmailValid(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
