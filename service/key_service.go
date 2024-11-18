package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const ACCESS_KEY_LENGTH = 10
const SECRET_KEY_LENGTH = 20

var ErrMissingAccessKey = errors.New("access key not found")
var ErrMissingSecretKey = errors.New("secret key not found")
var ErrGenKey = errors.New("error while generating key")

type KeyPair struct {
	accessKey     string
	secretKey     string
	secretKeyHash string
}

func NewKeyPair(accessKey, secretKey, secretKeyHash string) *KeyPair {
	return &KeyPair{
		accessKey:     accessKey,
		secretKey:     secretKey,
		secretKeyHash: secretKeyHash,
	}
}

func (kp *KeyPair) GetKeysPrivate() (string, string) {
	return kp.accessKey, kp.secretKey
}

func (kp *KeyPair) GetSecretKeyHash() string {
	return kp.secretKeyHash
}

type KeyGenerator interface {
	GenerateKeyPair() (*KeyPair, error)
}

type KeyService struct{}

func NewKeyService() KeyService {
	return KeyService{}
}

func (t KeyService) GenerateKeyPair() (*KeyPair, error) {
	accessKey, err := t.generateRandomString(ACCESS_KEY_LENGTH)
	if err != nil {
		log.Printf("Error while generating access key: %v\n", err)
		return nil, fmt.Errorf("%w: %w", ErrGenKey, err)
	}

	secretKey, err := t.generateRandomString(SECRET_KEY_LENGTH)
	if err != nil {
		log.Printf("Error while generating secret key: %v\n", err)
		return nil, fmt.Errorf("%w: %w", ErrGenKey, err)
	}

	hashedKey, err := bcrypt.GenerateFromPassword([]byte(secretKey), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error while generating secret key hash: %v\n", err)
		return nil, fmt.Errorf("%w: %w", ErrGenKey, err)
	}

	keyPair := &KeyPair{
		accessKey:     accessKey,
		secretKey:     secretKey,
		secretKeyHash: string(hashedKey),
	}

	return keyPair, nil
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
