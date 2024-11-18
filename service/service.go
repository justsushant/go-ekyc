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

func (c Service) GenerateKeyPair() (*KeyPair, error) {
	return c.keyService.GenerateKeyPair()
}

func (c Service) SaveSignupData(payload types.SignupPayload, keyPair *KeyPair) error {
	planId, err := c.dataStore.GetPlanIdFromName(payload.Plan)
	if err != nil {
		return err
	}

	accessKey, _ := keyPair.GetKeysPrivate()
	secretKeyHash := keyPair.GetSecretKeyHash()

	err = c.dataStore.InsertClientData(planId, payload, accessKey, secretKeyHash)
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

func (c Service) SaveFile(fileHeader *multipart.FileHeader, uploadMetaData *types.UploadMetaData) error {
	err := c.fileStore.SaveFileToBucket(fileHeader, uploadMetaData.FilePath)
	if err != nil {
		return err
	}

	err = c.dataStore.InsertUploadMetaData(uploadMetaData)
	if err != nil {
		return err
	}

	return nil
}

func (c Service) ValidateImage(payload types.FaceMatchPayload) error {
	// fetching meta data of images by uuid
	imgData1, err := c.dataStore.GetMetaDataByUUID(payload.ImageID1)
	if err != nil {
		return err
	}
	imgData2, err := c.dataStore.GetMetaDataByUUID(payload.ImageID2)
	if err != nil {
		return err
	}

	// if image data is nil (for nonexistent uuid case)
	if imgData1 == nil || imgData2 == nil {
		return ErrInvalidImgId
	}

	// if image belong to different clients
	if imgData1.ClientID != imgData2.ClientID {
		return ErrInvalidImgId
	}

	// if images are not of faces
	if imgData1.Type != types.FaceType || imgData2.Type != types.FaceType {
		return ErrNotFaceImg
	}

	return nil
}

func (c Service) CalcFaceMatchScore(payload types.FaceMatchPayload) (int, error) {
	return 0, nil
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
