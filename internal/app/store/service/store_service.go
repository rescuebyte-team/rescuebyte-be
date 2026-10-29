package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/domain/helpers/pagination"
)

func parseAndNormalizeTime(tStr string) (string, error) {
	trimmed := strings.TrimSpace(tStr)
	if trimmed == "" {
		return "", errx.ErrInvalidTimeFormat
	}

	// Normalize dot to colon (e.g. "08.00" -> "08:00")
	normalized := strings.ReplaceAll(trimmed, ".", ":")

	var parsed time.Time
	var err error
	if len(normalized) == 5 {
		parsed, err = time.Parse("15:04", normalized)
	} else if len(normalized) == 8 {
		parsed, err = time.Parse("15:04:05", normalized)
	} else {
		return "", errx.ErrInvalidTimeFormat
	}

	if err != nil {
		return "", errx.ErrInvalidTimeFormat
	}

	return parsed.Format("15:04:00"), nil
}

func (s *storeService) CreateStore(ctx context.Context, userID uuid.UUID, req dto.CreateStoreRequest) (*dto.StoreResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	opStart, err := parseAndNormalizeTime(req.OperationalStart)
	if err != nil {
		return nil, errx.ErrInvalidTimeFormat.WithLocation("body.operational_start")
	}

	opEnd, err := parseAndNormalizeTime(req.OperationalEnd)
	if err != nil {
		return nil, errx.ErrInvalidTimeFormat.WithLocation("body.operational_end")
	}

	for _, day := range req.OperationalDays {
		if day < 1 || day > 7 {
			return nil, errx.ErrBadRequest.WithMessage("Operational days must be between 1 (Monday) and 7 (Sunday)").WithLocation("body.operational_days")
		}
	}

	existingStore, err := s.storeRepo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingStore != nil {
		return nil, errx.ErrStoreAlreadyExists
	}

	now := time.Now()
	storeID, err := s.uuid.NewV7()
	if err != nil {
		return nil, err
	}
	store := entity.Store{
		ID:                storeID,
		UserID:            userID,
		Name:              req.Name,
		ProfilePictureURL: req.ProfilePictureURL,
		Description:       req.Description,
		WhatsAppNumber:    req.WhatsAppNumber,
		Email:             req.Email,
		Address:           req.Address,
		City:              req.City,
		District:          req.District,
		Status:            "ACTIVE",
		OperationalStart:  opStart,
		OperationalEnd:    opEnd,
		OperationalDays:   req.OperationalDays,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.storeRepo.CreateStore(txCtx, store); err != nil {
			return err
		}

		user, err := s.userRepo.GetUserByID(txCtx, userID)
		if err != nil {
			return err
		}
		if user.RoleID == enums.RoleUserIdx {
			user.RoleID = enums.RoleSellerIdx
			if err := s.userRepo.UpdateUser(txCtx, *user); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	res := dto.FormatToStoreResponse(store)
	return &res, nil
}

func (s *storeService) GetMyStore(ctx context.Context, userID uuid.UUID) (*dto.StoreResponse, error) {
	store, err := s.storeRepo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errx.ErrStoreNotFound
	}

	res := dto.FormatToStoreResponse(*store)
	return &res, nil
}

func (s *storeService) GetAllStores(ctx context.Context, query dto.GetAllStoresQuery) (dto.GetStoresResponse, error) {
	if err := s.validator.Validate(query); err != nil {
		return dto.GetStoresResponse{}, err
	}

	pagination.Normalize(&query.Limit, &query.Page, &query.Offset, "store")

	stores, total, err := s.storeRepo.GetAllStores(ctx, query)
	if err != nil {
		return dto.GetStoresResponse{}, err
	}

	responses := make([]dto.StoreResponse, 0, len(stores))
	for _, store := range stores {
		responses = append(responses, dto.FormatToStoreResponse(store))
	}

	return dto.GetStoresResponse{
		Stores: responses,
		Meta:   pagination.NewPaginationResponse(total, query.Page, query.Limit),
	}, nil
}

func (s *storeService) GetStoreByID(ctx context.Context, id uuid.UUID) (*dto.StoreResponse, error) {
	store, err := s.storeRepo.GetStoreByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := dto.FormatToStoreResponse(*store)
	return &res, nil
}

func (s *storeService) UpdateStoreProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateStoreProfileRequest) (*dto.StoreResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	store, err := s.storeRepo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errx.ErrStoreNotFound
	}

	store.Name = req.Name
	store.ProfilePictureURL = req.ProfilePictureURL
	store.Description = req.Description
	store.WhatsAppNumber = req.WhatsAppNumber
	store.Email = req.Email
	store.Address = req.Address
	store.City = req.City
	store.District = req.District
	store.Latitude = req.Latitude
	store.Longitude = req.Longitude
	store.UpdatedAt = time.Now()

	if err := s.storeRepo.UpdateStoreProfile(ctx, *store); err != nil {
		return nil, err
	}

	res := dto.FormatToStoreResponse(*store)
	return &res, nil
}

func (s *storeService) UpdateStoreOperational(ctx context.Context, userID uuid.UUID, req dto.UpdateStoreOperationalRequest) (*dto.StoreResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	opStart, err := parseAndNormalizeTime(req.OperationalStart)
	if err != nil {
		return nil, errx.ErrInvalidTimeFormat.WithLocation("body.operational_start")
	}

	opEnd, err := parseAndNormalizeTime(req.OperationalEnd)
	if err != nil {
		return nil, errx.ErrInvalidTimeFormat.WithLocation("body.operational_end")
	}

	for _, day := range req.OperationalDays {
		if day < 1 || day > 7 {
			return nil, errx.ErrBadRequest.WithMessage("Operational days must be between 1 (Monday) and 7 (Sunday)").WithLocation("body.operational_days")
		}
	}

	store, err := s.storeRepo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errx.ErrStoreNotFound
	}

	store.Status = req.Status
	store.OperationalStart = opStart
	store.OperationalEnd = opEnd
	store.OperationalDays = req.OperationalDays
	store.UpdatedAt = time.Now()

	if err := s.storeRepo.UpdateStoreOperational(ctx, *store); err != nil {
		return nil, err
	}

	res := dto.FormatToStoreResponse(*store)
	return &res, nil
}

func (s *storeService) ToggleStoreStatus(ctx context.Context, userID uuid.UUID) (*dto.StoreResponse, error) {
	store, err := s.storeRepo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errx.ErrStoreNotFound
	}

	newStatus := "INACTIVE"
	if store.Status != "ACTIVE" {
		newStatus = "ACTIVE"
	}

	if err := s.storeRepo.UpdateStoreStatus(ctx, store.ID, newStatus); err != nil {
		return nil, err
	}

	store.Status = newStatus
	res := dto.FormatToStoreResponse(*store)
	return &res, nil
}
