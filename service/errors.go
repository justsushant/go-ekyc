package service

import (
	"errors"
)

var (
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPlan       = errors.New("invalid plan, supported plans are basic, advanced, or enterprise")
	ErrInvalidFileType   = errors.New("invalid type, supported types are face or id_card")
	ErrInvalidFileFormat = errors.New("invalid file format, supported formats are png or jpeg")
	ErrInvalidImgId      = errors.New("invalid or missing image id")
	ErrNotFaceImg        = errors.New("not a face image")
	ErrNotIDCardImg      = errors.New("not an id card image")
)
