package service

import (
	"encoding/json"
	"log"
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
	faceMatch  FaceMatcher
	ocrService OCRPerformer
	queue      TaskQueue
	uuid       UUIDGen
}

func NewService(dataStore store.DataStore, fileStore store.FileStore, keyService KeyGenerator, faceMatch FaceMatcher, ocrService OCRPerformer, queue TaskQueue, uuid UUIDGen) Service {
	return Service{
		dataStore:  dataStore,
		keyService: keyService,
		fileStore:  fileStore,
		faceMatch:  faceMatch,
		ocrService: ocrService,
		queue:      queue,
		uuid:       uuid,
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
	// save the file to filestore
	err := c.fileStore.SaveFileToBucket(fileHeader, uploadMetaData.FilePath)
	if err != nil {
		return err
	}

	// save the file upload metadata to db
	err = c.dataStore.InsertUploadMetaData(uploadMetaData)
	if err != nil {
		return err
	}

	return nil
}

func (c Service) ValidateImage(payload types.FaceMatchPayload, clientID int) error {
	// fetching meta data of images by uuid
	imgData1, err := c.dataStore.GetMetaDataByUUID(payload.Image1)
	if err != nil {
		return err
	}
	imgData2, err := c.dataStore.GetMetaDataByUUID(payload.Image2)
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

	// if client and image have different client id
	if imgData1.ClientID != clientID {
		return ErrInvalidImgId
	}

	// if images are not of faces
	if imgData1.Type != types.FaceType || imgData2.Type != types.FaceType {
		return ErrNotFaceImg
	}

	return nil
}

func (c Service) CalcAndSaveFaceMatchScore(payload types.FaceMatchPayload, clientID int) (int, error) {
	score, err := c.faceMatch.CalcFaceMatchScore(payload)
	if err != nil {
		return 0, err
	}

	// fetching meta data of images by uuid
	imgData1, err := c.dataStore.GetMetaDataByUUID(payload.Image1)
	if err != nil {
		return 0, err
	}
	imgData2, err := c.dataStore.GetMetaDataByUUID(payload.Image2)
	if err != nil {
		return 0, err
	}

	result := &types.FaceMatchData{
		ClientID: clientID,
		ImageID1: imgData1.Id,
		ImageID2: imgData2.Id,
		Score:    score,
	}

	err = c.dataStore.InsertFaceMatchResult(result)
	if err != nil {
		return 0, err
	}

	return score, nil
}

func (c Service) PerformAndSaveOCR(payload types.OCRPayload, clientID int) (*types.OCRResponse, error) {
	data, err := c.ocrService.PerformOCR(payload)
	if err != nil {
		return nil, err
	}

	imgData, err := c.dataStore.GetMetaDataByUUID(payload.Image)
	if err != nil {
		return nil, err
	}

	result := &types.OCRData{
		ImageID:  imgData.Id,
		ClientID: clientID,
		Data:     data.String(),
	}

	err = c.dataStore.InsertOCRResult(result)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c Service) ValidateImageOCR(payload types.OCRPayload, clientID int) error {
	// fetching meta data of image by uuid
	imgData, err := c.dataStore.GetMetaDataByUUID(payload.Image)
	if err != nil {
		return err
	}

	// if image data is nil (for nonexistent uuid case)
	if imgData == nil {
		return ErrInvalidImgId
	}

	// if image belong to different clients
	if imgData.ClientID != clientID {
		return ErrInvalidImgId
	}

	// if image is not of id card
	if imgData.Type != types.IdCardType {
		return ErrNotIDCardImg
	}

	return nil
}

func (c Service) PerformFaceMatchAsync(payload types.FaceMatchPayload, clientID int) (string, error) {
	// make validations for the images
	err := c.validateImagesForFaceMatch(payload, clientID)
	if err != nil {
		return "", err
	}

	// generate the job id
	jobID := c.uuid.New()

	// get metadata of both images
	img1Data, _ := c.dataStore.GetMetaDataByUUID(payload.Image1)
	img2Data, _ := c.dataStore.GetMetaDataByUUID(payload.Image2)

	// mark the job started on the db
	err = c.dataStore.InsertFaceMatchJobCompleted(img1Data.Id, img2Data.Id, clientID, jobID)
	if err != nil {
		return "", err
	}

	// push the job onto the queue
	queuePayload := types.FaceMatchQueuePayload{
		Type: types.FaceMatchWorkType,
		Msg: types.FaceMatchInternalPayload{
			JobID:  jobID,
			Image1: payload.Image1,
			Image2: payload.Image2,
		},
	}
	jsonBytes, err := json.Marshal(queuePayload)
	if err != nil {
		log.Println("Error while marshalling JSON: ", err)
	}
	c.queue.PushJobOnQueue(jsonBytes)

	return jobID, nil
}

func (c Service) PerformOCRAsync(payload types.OCRPayload, clientID int) (string, error) {
	// make validations for the images
	err := c.validateImageForOCR(payload, clientID)
	if err != nil {
		return "", err
	}

	// generate the job id
	jobID := c.uuid.New()

	// get metadata of both images
	imgData, _ := c.dataStore.GetMetaDataByUUID(payload.Image)

	// mark the job started on the db
	err = c.dataStore.InsertOCRJobCompleted(imgData.Id, clientID, jobID)
	if err != nil {
		return "", err
	}

	// push the job onto the queue
	queuePayload := types.OCRQueuePayload{
		Type: types.OCRWorkType,
		Msg: types.OCRInternalPayload{
			JobID: jobID,
			Image: payload.Image,
		},
	}
	jsonBytes, err := json.Marshal(queuePayload)
	if err != nil {
		log.Println("Error while marshalling JSON: ", err)
	}
	c.queue.PushJobOnQueue(jsonBytes)

	return jobID, nil
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

func (c Service) validateImagesForFaceMatch(payload types.FaceMatchPayload, clientID int) error {
	// fetching meta data of images by uuid
	imgData1, err := c.dataStore.GetMetaDataByUUID(payload.Image1)
	if err != nil {
		return err
	}
	imgData2, err := c.dataStore.GetMetaDataByUUID(payload.Image2)
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

	// if client and image have different client id
	if imgData1.ClientID != clientID {
		return ErrInvalidImgId
	}

	// if images are not of faces
	if imgData1.Type != types.FaceType {
		return ErrNotFaceImg
	}
	if imgData2.Type != types.FaceType {
		return ErrNotFaceImg
	}

	return nil
}

func (c Service) validateImageForOCR(payload types.OCRPayload, clientID int) error {
	// fetching meta data of image by uuid
	imgData, err := c.dataStore.GetMetaDataByUUID(payload.Image)
	if err != nil {
		return err
	}

	// if image data is nil (for nonexistent uuid case)
	if imgData == nil {
		return ErrInvalidImgId
	}

	// if image belong to different clients
	if imgData.ClientID != clientID {
		return ErrInvalidImgId
	}

	// if image is not of id card
	if imgData.Type != types.IdCardType {
		return ErrNotIDCardImg
	}

	return nil
}
