package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/env"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	if valErr := s.validator.Validate(req); valErr != nil {
		return nil, valErr
	}

	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errx.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errx.ErrInternalServer.WithError(err)
	}
	passHashStr := string(hashedPassword)

	userID, err := s.uuid.NewV7()
	if err != nil {
		userID = uuid.New()
	}

	now := time.Now()
	user := entity.User{
		ID:           userID,
		Email:        req.Email,
		PasswordHash: &passHashStr,
		Name:         req.Name,
		Role:         enums.RoleUser,
		Provider:     enums.ProviderLocal,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.Create(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.CreateRefreshToken(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	rtID, _ := s.uuid.NewV7()
	rtEntity := entity.RefreshToken{
		ID:        rtID,
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: now.Add(env.AppEnv.JwtRefreshExpTime),
		CreatedAt: now,
	}

	if err := s.userRepo.CreateRefreshToken(ctx, rtEntity); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			Provider:  user.Provider,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	if valErr := s.validator.Validate(req); valErr != nil {
		return nil, valErr
	}

	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errx.ErrInvalidCredentials
	}

	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return nil, errx.ErrPasswordNotSet
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errx.ErrInvalidCredentials
	}

	accessToken, err := s.jwt.Create(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.CreateRefreshToken(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	rtID, _ := s.uuid.NewV7()
	rtEntity := entity.RefreshToken{
		ID:        rtID,
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: now.Add(env.AppEnv.JwtRefreshExpTime),
		CreatedAt: now,
	}

	if err := s.userRepo.CreateRefreshToken(ctx, rtEntity); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
			Role:      user.Role,
			Provider:  user.Provider,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	if valErr := s.validator.Validate(req); valErr != nil {
		return nil, valErr
	}

	var claims jwt.Claims
	if err := s.jwt.DecodeRefreshToken(req.RefreshToken, &claims); err != nil {
		return nil, errx.ErrInvalidRefreshToken.WithError(err)
	}

	storedToken, err := s.userRepo.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if storedToken.RevokedAt != nil {
		return nil, errx.ErrRefreshTokenRevoked
	}
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, errx.ErrInvalidRefreshToken.WithMessage("Refresh token has expired.")
	}

	user, err := s.userRepo.GetUserByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, err
	}

	// Revoke old token & generate new tokens atomically
	newAccessToken, err := s.jwt.Create(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwt.CreateRefreshToken(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	newRtID, _ := s.uuid.NewV7()
	newRtEntity := entity.RefreshToken{
		ID:        newRtID,
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: now.Add(env.AppEnv.JwtRefreshExpTime),
		CreatedAt: now,
	}

	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.userRepo.RevokeRefreshToken(txCtx, req.RefreshToken); err != nil {
			return err
		}
		return s.userRepo.CreateRefreshToken(txCtx, newRtEntity)
	})
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
			Role:      user.Role,
			Provider:  user.Provider,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return s.userRepo.RevokeRefreshToken(ctx, refreshToken)
}

func (s *authService) GetGoogleAuthURL(ctx context.Context) (dto.GoogleAuthURLResponse, error) {
	state := s.googleOAuth.GenerateState()
	authURL := s.googleOAuth.GetAuthURL(state)

	return dto.GoogleAuthURLResponse{
		AuthURL: authURL,
		State:   state,
	}, nil
}

func (s *authService) HandleGoogleCallback(ctx context.Context, query dto.GoogleCallbackQuery) (*dto.AuthResponse, error) {
	if valErr := s.validator.Validate(query); valErr != nil {
		return nil, valErr
	}

	oauthToken, err := s.googleOAuth.Exchange(ctx, query.Code)
	if err != nil {
		return nil, errx.ErrOAuthExchangeFailed.WithError(err)
	}

	googleUser, err := s.googleOAuth.GetUserInfo(ctx, oauthToken)
	if err != nil {
		return nil, errx.ErrOAuthGetUserInfoFailed.WithError(err)
	}

	// Check if user already exists with Google ProviderID or with Email
	var user *entity.User
	user, err = s.userRepo.GetUserByProviderID(ctx, enums.ProviderGoogle, googleUser.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if user == nil {
		// Check by email
		existingByEmail, err := s.userRepo.GetUserByEmail(ctx, googleUser.Email)
		if err != nil {
			return nil, err
		}

		if existingByEmail != nil {
			// Link user
			user = existingByEmail
			user.ProviderID = &googleUser.ID
			if user.AvatarURL == nil && googleUser.Picture != "" {
				user.AvatarURL = &googleUser.Picture
			}
			if err := s.userRepo.UpdateUser(ctx, *user); err != nil {
				return nil, err
			}
		} else {
			// Create brand new user
			newUserID, _ := s.uuid.NewV7()
			newUser := entity.User{
				ID:         newUserID,
				Email:      googleUser.Email,
				Name:       googleUser.Name,
				AvatarURL:  &googleUser.Picture,
				Role:       enums.RoleUser,
				Provider:   enums.ProviderGoogle,
				ProviderID: &googleUser.ID,
				CreatedAt:  now,
				UpdatedAt:  now,
			}
			if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
				return nil, err
			}
			user = &newUser
		}
	}

	accessToken, err := s.jwt.Create(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.CreateRefreshToken(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	rtID, _ := s.uuid.NewV7()
	rtEntity := entity.RefreshToken{
		ID:        rtID,
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: now.Add(env.AppEnv.JwtRefreshExpTime),
		CreatedAt: now,
	}

	if err := s.userRepo.CreateRefreshToken(ctx, rtEntity); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
			Role:      user.Role,
			Provider:  user.Provider,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) GetMe(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
