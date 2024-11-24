package store

import (
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type FileStore interface {
	SaveFile(file *types.FileUpload) error
	GetFile(filePath string) ([]byte, error)
}
