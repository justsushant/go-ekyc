package store

import "mime/multipart"

type FileStore interface {
	SaveFile(fileHeader *multipart.FileHeader) error
}
