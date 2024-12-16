package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"regexp"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/store"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	// store
	dataStore  store.DataStore
	fileStore  store.FileStore
	cacheStore store.CacheStore

	// business logic
	faceMatch  FaceMatcher
	ocrService OCRPerformer

	// task processing
	queue TaskQueue

	// utils
	keyService KeyGenerator
	uuid       UUIDGen
}

type ServiceConfig struct {
	DataStore  store.DataStore
	FileStore  store.FileStore
	CacheStore store.CacheStore
	KeyService KeyGenerator
	FaceMatch  FaceMatcher
	OCR        OCRPerformer
	Queue      TaskQueue
	UUID       UUIDGen
}

func NewService(config *ServiceConfig) Service {
	return Service{
		dataStore:  config.DataStore,
		keyService: config.KeyService,
		fileStore:  config.FileStore,
		cacheStore: config.CacheStore,
		faceMatch:  config.FaceMatch,
		ocrService: config.OCR,
		queue:      config.Queue,
		uuid:       config.UUID,
	}
}

func (c Service) SignupClient(payload types.SignupPayload) (*KeyPair, error) {
	// apply validations on payload
	if err := validateEmail(payload.Email); err != nil {
		return nil, err
	}
	if err := validatePlan(payload.Plan); err != nil {
		return nil, err
	}

	// generate keys
	keyPair, err := c.keyService.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	// save to db
	planId, err := c.dataStore.GetPlanIdFromName(payload.Plan)
	if err != nil {
		return nil, err
	}
	accessKey, _ := keyPair.GetKeysPrivate()
	secretKeyHash := keyPair.GetSecretKeyHash()

	err = c.dataStore.InsertClientData(planId, payload, accessKey, secretKeyHash)
	if err != nil {
		return nil, err
	}

	return keyPair, nil
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
	fileReader, err := fileHeader.Open()
	if err != nil {
		log.Printf("Error while reading the file: %s\n", err.Error())
		return err
	}

	file := &types.FileUpload{
		Name:    uploadMetaData.FilePath,
		Content: fileReader,
		Size:    fileHeader.Size,
		Headers: map[string]string{
			"Content-Type": fileHeader.Header.Get("Content-Type"),
		},
	}

	err = c.fileStore.SaveFile(file)
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

func (c Service) PerformFaceMatch(payload types.FaceMatchPayload, clientID int) (string, error) {
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
	err = c.dataStore.InsertFaceMatchJobCreated(img1Data.Id, img2Data.Id, clientID, jobID)
	if err != nil {
		return "", err
	}

	// push the job onto the queue
	queuePayload := types.FaceMatchQueuePayload{
		Type: types.FACE_MATCH_WORK_TYPE,
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

func (c Service) PerformOCR(payload types.OCRPayload, clientID int) (string, error) {
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
	err = c.dataStore.InsertOCRJobCreated(imgData.Id, clientID, jobID)
	if err != nil {
		return "", err
	}

	// push the job onto the queue
	queuePayload := types.OCRQueuePayload{
		Type: types.OCR_WORK_TYPE,
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
	case types.FACE_TYPE, types.ID_CARD_TYPE:
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
	case types.BASIC_PLAN, types.ADVANCE_PLAN, types.ENTERPRISE_PLAN:
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
	if imgData1.Type != types.FACE_TYPE {
		return ErrNotFaceImg
	}
	if imgData2.Type != types.FACE_TYPE {
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
	if imgData.Type != types.ID_CARD_TYPE {
		return ErrNotIDCardImg
	}

	return nil
}

func (c Service) GetJobDetailsByJobID(jobID, jobType string) (*types.JobRecord, error) {
	switch jobType {
	case types.FACE_MATCH_WORK_TYPE:
		return c.dataStore.GetFaceMatchByJobID(jobID)
	case types.OCR_WORK_TYPE:
		return c.dataStore.GetOCRByJobID(jobID)
	default:
		return nil, fmt.Errorf("invalid job type")
	}
}

func (c Service) FetchDataFromCache(payload interface{}, clientID int, jobType string) (string, bool) {
	var cacheKey string
	switch p := payload.(type) {
	case types.FaceMatchPayload:
		key1 := fmt.Sprintf("%s:%d:%s", jobType, clientID, p.Image1)
		key2 := fmt.Sprintf("%s:%d:%s", jobType, clientID, p.Image2)
		cacheKey = c.makeHash(key1, key2)
	case types.OCRPayload:
		key := fmt.Sprintf("%s:%d:%s", jobType, clientID, p.Image)
		cacheKey = c.makeHash(key)
	}

	// TODO: possible bug
	// this may cause bugs if switch is unable to match the payload with a type
	// & both falsy values will be returned and handler will send empty string in response to client
	val := c.getObjFromCache(cacheKey)
	if val != "" {
		return val, true
	}
	return "", false
}

func (c Service) getObjFromCache(key string) string {
	val, err := c.cacheStore.GetObject(key)
	if err != nil && err != redis.Nil {
		log.Printf("Error while fetching data from cache (%s): %s\n", key, err.Error())
		return ""
	}
	return val
}

func (c Service) setObjInCache(key, val string) {
	err := c.cacheStore.SetObject(key, val)
	if err != nil {
		log.Printf("Error while setting data in cache (%s): %s\n", key, err.Error())
	}
}

func (c Service) SetDataInCache(payload interface{}, clientID int, jobType, jobID string) {
	switch p := payload.(type) {
	case types.FaceMatchPayload:
		key1 := fmt.Sprintf("%s:%d:%s", jobType, clientID, p.Image1)
		key2 := fmt.Sprintf("%s:%d:%s", jobType, clientID, p.Image2)
		cacheKey := c.makeHash(key1, key2)
		c.setObjInCache(cacheKey, jobID)
	case types.OCRPayload:
		cacheKey := c.makeHash(fmt.Sprintf("%s:%d:%s", jobType, clientID, p.Image))
		c.setObjInCache(cacheKey, jobID)
	}
}

func (c Service) makeHash(keys ...string) string {
	// for face match hash
	if len(keys) == 2 {
		key1Bytes := []byte(keys[0])
		key2Bytes := []byte(keys[1])

		// xor op on byte by byte
		xorBytes := make([]byte, len(key1Bytes))
		for i := 0; i < len(key1Bytes); i++ {
			xorBytes[i] = key1Bytes[i] ^ key2Bytes[i]
		}

		// generate SHA256 hash of the xor result
		hash := sha256.Sum256(xorBytes)
		return hex.EncodeToString(hash[:])
	}

	// for ocr key hash
	if len(keys) == 1 {
		keyBytes := []byte(keys[0])
		hash := sha256.Sum256(keyBytes)
		return hex.EncodeToString(hash[:])
	}

	// default case
	return keys[0]
}
