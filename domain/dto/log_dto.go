package dto

import "github.com/google/uuid"

type LogActor struct {
	ID    uuid.UUID `validate:"omitempty,uuid"`
	Name  string
	Email string `validate:"required,email"`
}
