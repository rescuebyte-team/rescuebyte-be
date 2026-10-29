package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/helpers/pagination"
)

type CreateStoreRequest struct {
	Name              string  `json:"name" validate:"required,min=2,max=255"`
	ProfilePictureURL *string `json:"profile_picture_url" validate:"omitempty,url"`
	Description       *string `json:"description" validate:"omitempty"`
	WhatsAppNumber    string  `json:"whatsapp_number" validate:"required,min=8,max=20"`
	Email             string  `json:"email" validate:"required,email"`
	Address           string  `json:"address" validate:"required"`
	City              string  `json:"city" validate:"required,max=100"`
	District          string  `json:"district" validate:"required,max=100"`
	OperationalStart  string  `json:"operational_start" validate:"required"` // e.g. "08:00"
	OperationalEnd    string  `json:"operational_end" validate:"required"`   // e.g. "21:00"
	OperationalDays   []int16 `json:"operational_days" validate:"required,min=1"`
	Latitude          float64 `json:"latitude" validate:"required,latitude"`
	Longitude         float64 `json:"longitude" validate:"required,longitude"`
}

type UpdateStoreProfileRequest struct {
	Name              string  `json:"name" validate:"required,min=2,max=255"`
	ProfilePictureURL *string `json:"profile_picture_url" validate:"omitempty,url"`
	Description       *string `json:"description" validate:"omitempty"`
	WhatsAppNumber    string  `json:"whatsapp_number" validate:"required,min=8,max=20"`
	Email             string  `json:"email" validate:"required,email"`
	Address           string  `json:"address" validate:"required"`
	City              string  `json:"city" validate:"required,max=100"`
	District          string  `json:"district" validate:"required,max=100"`
	Latitude          float64 `json:"latitude" validate:"required,latitude"`
	Longitude         float64 `json:"longitude" validate:"required,longitude"`
}

type UpdateStoreOperationalRequest struct {
	Status           string  `json:"status" validate:"required,oneof=ACTIVE INACTIVE"`
	OperationalStart string  `json:"operational_start" validate:"required"`
	OperationalEnd   string  `json:"operational_end" validate:"required"`
	OperationalDays  []int16 `json:"operational_days" validate:"required,min=1"`
}

type GetAllStoresQuery struct {
	Limit  uint64 `query:"limit" validate:"omitempty,number"`
	Page   uint64 `query:"page" validate:"omitempty,number,gte=1"`
	Offset uint64
	Name   string `query:"name" validate:"omitempty"`
	Nama   string `query:"nama" validate:"omitempty"`
}

type StoreResponse struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	Name              string    `json:"name"`
	ProfilePictureURL *string   `json:"profile_picture_url"`
	Description       *string   `json:"description"`
	WhatsAppNumber    string    `json:"whatsapp_number"`
	Email             string    `json:"email"`
	Address           string    `json:"address"`
	City              string    `json:"city"`
	District          string    `json:"district"`
	Status            string    `json:"status"`
	OperationalStart  string    `json:"operational_start"`
	OperationalEnd    string    `json:"operational_end"`
	OperationalDays   []int16   `json:"operational_days"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type GetStoresResponse struct {
	Stores []StoreResponse     `json:"stores"`
	Meta   pagination.Response `json:"meta"`
}

func FormatToStoreResponse(store entity.Store) StoreResponse {
	return StoreResponse{
		ID:                store.ID,
		UserID:            store.UserID,
		Name:              store.Name,
		ProfilePictureURL: store.ProfilePictureURL,
		Description:       store.Description,
		WhatsAppNumber:    store.WhatsAppNumber,
		Email:             store.Email,
		Address:           store.Address,
		City:              store.City,
		District:          store.District,
		Status:            store.Status,
		OperationalStart:  store.OperationalStart,
		OperationalEnd:    store.OperationalEnd,
		OperationalDays:   store.OperationalDays,
		Latitude:          store.Latitude,
		Longitude:         store.Longitude,
		CreatedAt:         store.CreatedAt,
		UpdatedAt:         store.UpdatedAt,
	}
}
