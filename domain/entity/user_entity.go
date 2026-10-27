package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

type User struct {
	ID           uuid.UUID          `db:"id" json:"id"`
	Email        string             `db:"email" json:"email"`
	PasswordHash *string            `db:"password_hash" json:"-"`
	Name         string             `db:"name" json:"name"`
	AvatarURL    *string            `db:"avatar_url" json:"avatar_url"`
	RoleID       enums.UserRoleIdx  `db:"role_id" json:"role_id"`
	Provider     enums.AuthProvider `db:"provider" json:"provider"`
	ProviderID   *string            `db:"provider_id" json:"provider_id"`
	PhoneNumber  *string            `db:"phone_number" json:"phone_number"`
	Gender       *enums.GenderIdx   `db:"gender" json:"gender"`
	DateOfBirth  *time.Time         `db:"date_of_birth" json:"date_of_birth"`
	CreatedAt    time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `db:"updated_at" json:"updated_at"`
}
