package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID                uuid.UUID      `db:"id" json:"id"`
	UserID            uuid.UUID      `db:"user_id" json:"user_id"`
	Name              string         `db:"name" json:"name"`
	ProfilePictureURL sql.NullString `db:"profile_picture_url" json:"profile_picture_url"`
	Description       string         `db:"description" json:"description"`
	WhatsAppNumber    string         `db:"whatsapp_number" json:"whatsapp_number"`
	Email             string         `db:"email" json:"email"`
	Address           string         `db:"address" json:"address"`
	City              string         `db:"city" json:"city"`
	District          string         `db:"district" json:"district"`
	Status            string         `db:"status" json:"status"`
	OperationalStart  time.Time      `db:"operational_start" json:"operational_start"`
	OperationalEnd    time.Time      `db:"operational_end" json:"operational_end"`
	OperationalDays   []int8         `db:"operational_days" json:"operational_days"`
	Latitude          float64        `db:"latitude" json:"latitude"`
	Longitude         float64        `db:"longitude" json:"longitude"`
	CreatedAt         time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at" json:"updated_at"`
}
