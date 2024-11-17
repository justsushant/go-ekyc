package service

import (
	"crypto/rand"
	"math/big"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

const ACCESS_KEY_LENGTH = 10
const SECRET_KEY_LENGTH = 20

type KeyPair struct {
	AccessKey string
	SecretKey string
}

type KeyGenerator interface {
	GenerateKeyPair(payload types.SignupPayload) (*KeyPair, error)
}

type KeyService struct{}

func NewKeyService() KeyService {
	return KeyService{}
}

func (t KeyService) GenerateKeyPair(payload types.SignupPayload) (*KeyPair, error) {
	accessKey, err := t.generateRandomString(ACCESS_KEY_LENGTH)
	if err != nil {
		return nil, err
	}

	secretKey, err := t.generateRandomString(SECRET_KEY_LENGTH)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, nil
}

func (t KeyService) generateRandomString(n int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result []byte

	for i := 0; i < n; i++ {
		// generate a random number within the range of the charset
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result = append(result, charset[randomIndex.Int64()])
	}

	return string(result), nil
}
