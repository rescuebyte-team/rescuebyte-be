package dto

import (
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

type UpdateProfileRequest struct {
	Name        string           `json:"name" validate:"required,min=2,max=100"`
	AvatarURL   *string          `json:"avatar_url" validate:"omitempty,url"`
	PhoneNumber *string          `json:"phone_number" validate:"omitempty,min=8,max=20"`
	Gender      *enums.GenderIdx `json:"gender" validate:"omitempty,gender"`
	DateOfBirth *string          `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02"`
}

type UpdateUserRoleRequest struct {
	UserID uuid.UUID         `json:"user_id" validate:"required"`
	RoleID enums.UserRoleIdx `json:"role_id" validate:"required,user_role"`
}

type GetUsersQuery struct {
	Page   int                `query:"page"`
	Limit  int                `query:"limit"`
	RoleID *enums.UserRoleIdx `query:"role_id" validate:"omitempty,user_role"`
}
