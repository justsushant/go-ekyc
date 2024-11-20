package service

import "github.com/google/uuid"

type UUIDGen interface {
	New() string
}

type UuidService struct {}

func (u *UuidService) New() string {
	return uuid.New().String()
}
