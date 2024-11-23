package store

import "mime/multipart"

type FileStore interface {
	SaveFileToBucket(fileHeader *multipart.FileHeader, objectName string) error
	GetFile(filePath string) ([]byte, error)
}
