package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByProviderID(ctx context.Context, provider enums.AuthProvider, providerID string) (*entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) error

	// Refresh token operations
	CreateRefreshToken(ctx context.Context, token entity.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}

type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UserResponse, error)
	UpdateUserRole(ctx context.Context, req dto.UpdateUserRoleRequest) (*dto.UserResponse, error)
}
