package dto

import (
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

type UpdateProfileRequest struct {
	Name      string  `json:"name" validate:"required,min=2,max=100"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url"`
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
