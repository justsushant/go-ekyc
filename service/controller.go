package service

import (
	"mime/multipart"
	"path/filepath"
	"regexp"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/store"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type Service struct {
	dataStore  store.DataStore
	fileStore  store.FileStore
	keyService KeyGenerator
}

func NewService(dataStore store.DataStore, fileStore store.FileStore, keyService KeyGenerator) Service {
	return Service{
		dataStore:  dataStore,
		keyService: keyService,
		fileStore:  fileStore,
	}
}

func (c Service) ValidatePayload(payload types.SignupPayload) error {
	if err := validateEmail(payload.Email); err != nil {
		return err
	}
	if err := validatePlan(payload.Plan); err != nil {
		return err
	}

	return nil
}

func (c Service) GenerateKeyPair(payload types.SignupPayload) (*KeyPair, error) {
	return nil, c.keyService.GenerateKeyPair(payload)
}

func (c Service) SaveSignupData(payload types.SignupPayload, refreshToken string) error {
	planId, err := c.dataStore.GetPlanIdFromName(payload.Plan)
	if err != nil {
		return err
	}

	err = c.dataStore.InsertClientData(planId, payload)
	if err != nil {
		return err
	}

	return nil
}

func (c Service) ValidateFile(fileName, fileType string) error {
	err := validateFileType(fileType)
	if err != nil {
		return err
	}

	err = validateFileExt(fileName)
	if err != nil {
		return err
	}

	return nil
}

func (c Service) SaveUploadedFile(fileHeader *multipart.FileHeader) error {
	err := c.fileStore.SaveFile(fileHeader)
	if err != nil {
		return err
	}
	return nil
}

func validateFileType(fileType string) error {
	switch fileType {
	case types.FaceType, types.IdCardType:
		return nil
	default:
		return ErrInvalidFileType
	}
}

func validateEmail(email string) error {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	if re.MatchString(email) {
		return nil
	}
	return ErrInvalidEmail
}

func validateFileExt(fileName string) error {
	switch filepath.Ext(fileName) {
	case types.VALID_FORMAT_PNG, types.VALID_FORMAT_JPEG:
		return nil
	default:
		return ErrInvalidFileFormat
	}
}

func validatePlan(plan string) error {
	switch plan {
	case types.BasicPlan, types.AdvancePlan, types.EnterprisePlan:
		return nil
	default:
		return ErrInvalidPlan
	}
}
