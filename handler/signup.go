package handler

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidPlan = errors.New("invalid plan, supported plans are basic, advanced, or enterprise")

type Plan string

const BasicPlan = "basic"
const AdvancePlan = "advance"
const EnterprisePlan = "enterprise"

func getPlanFromString(plan string) (Plan, error) {
	switch plan {
	case BasicPlan:
		return BasicPlan, nil
	case AdvancePlan:
		return AdvancePlan, nil
	case EnterprisePlan:
		return EnterprisePlan, nil
	default:
		return "", ErrInvalidPlan
	}
}

func SignupHandler(c *gin.Context) {
	var signupPaylod types.SignupPayload
	err := json.NewDecoder(c.Request.Body).Decode(&signupPaylod)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	_, err = isEmailValid(signupPaylod.Email)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

	_, err = getPlanFromString(signupPaylod.Plan)
	if err != nil {
		c.JSON(400, gin.H{"errorMessage": err.Error()})
		return
	}

}

func isEmailValid(email string) (bool, error) {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	if re.MatchString(email) {
		return true, nil
	}
	return false, ErrInvalidEmail
}
