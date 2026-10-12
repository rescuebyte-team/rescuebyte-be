package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	GetGoogleAuthURL(ctx context.Context) (dto.GoogleAuthURLResponse, error)
	HandleGoogleCallback(ctx context.Context, query dto.GoogleCallbackQuery) (*dto.AuthResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
}
