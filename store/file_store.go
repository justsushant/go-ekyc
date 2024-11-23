package store

import "mime/multipart"

type FileStore interface {
	SaveFile(fileHeader *multipart.FileHeader, objectName string) error
	GetFile(filePath string) ([]byte, error)
}
