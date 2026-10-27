package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
)

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := dto.FormatToUserResponse(*user)
	return &res, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	if valErr := s.validator.Validate(req); valErr != nil {
		return nil, valErr
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name
	if req.AvatarURL != nil {
		if *req.AvatarURL == "" {
			user.AvatarURL = nil
		} else {
			user.AvatarURL = req.AvatarURL
		}
	}
	if req.PhoneNumber != nil {
		if *req.PhoneNumber == "" {
			user.PhoneNumber = nil
		} else {
			user.PhoneNumber = req.PhoneNumber
		}
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}
	if req.DateOfBirth != nil {
		if *req.DateOfBirth == "" {
			user.DateOfBirth = nil
		} else {
			parsedDate, err := time.Parse("2006-01-02", *req.DateOfBirth)
			if err != nil {
				return nil, errx.ErrBadRequest.WithMessage("format date_of_birth tidak valid (harus YYYY-MM-DD)")
			}
			user.DateOfBirth = &parsedDate
		}
	}

	if err := s.userRepo.UpdateUser(ctx, *user); err != nil {
		return nil, err
	}

	res := dto.FormatToUserResponse(*user)
	return &res, nil
}

func (s *userService) UpdateUserRole(ctx context.Context, req dto.UpdateUserRoleRequest) (*dto.UserResponse, error) {
	if valErr := s.validator.Validate(req); valErr != nil {
		return nil, valErr
	}

	user, err := s.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	user.RoleID = req.RoleID

	if err := s.userRepo.UpdateUser(ctx, *user); err != nil {
		return nil, err
	}

	res := dto.FormatToUserResponse(*user)
	return &res, nil
}
