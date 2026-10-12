package uuid

import (
	"github.com/google/uuid"
)

type Service interface {
	NewV4() (uuid.UUID, error)
	NewV7() (uuid.UUID, error)
	Parse(s string) (uuid.UUID, error)
}

type uuidService struct{}

var UUID Service = &uuidService{}

func (u *uuidService) NewV4() (uuid.UUID, error) {
	return uuid.NewRandom()
}

func (u *uuidService) NewV7() (uuid.UUID, error) {
	return uuid.NewV7()
}

func (u *uuidService) Parse(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
