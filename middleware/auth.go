package middleware

import (
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/store"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthMiddleware struct {
	store store.DataStore
}

func NewAuthMiddleware(store store.DataStore) *AuthMiddleware {
	return &AuthMiddleware{store: store}
}

func (am *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract keys from headers
		accessKey := c.GetHeader("accessKey")
		secretKey := c.GetHeader("secretKey")

		if len(accessKey) == 0 || len(secretKey) == 0 {
			c.JSON(401, gin.H{"errorMessage": "invalid access or secret key"})
			return
		}

		// get user details on the basis of access key
		clientData, err := am.store.GetClientFromAccessKey(accessKey)
		if err != nil {
			c.JSON(500, gin.H{"errorMessage": err.Error()})
			return
		}

		// match the hash of the key
		err = bcrypt.CompareHashAndPassword([]byte(clientData.SecretKeyHash), []byte(secretKey))
		if err != nil {
			if err != nil {
				c.JSON(401, gin.H{"errorMessage": "invalid access or secret key"})
				return
			}
		}

		// set the client id on gin.Context
		c.Set("client_id", clientData.Id)

		// call the next handler
		c.Next()
	}
}
