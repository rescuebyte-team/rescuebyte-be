package service

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
	"github.com/rescuebyte/rescuebyte-be/pkg/oauth"
	pkgUUID "github.com/rescuebyte/rescuebyte-be/pkg/uuid"
	"github.com/rescuebyte/rescuebyte-be/pkg/validator"
)

type authService struct {
	transactor  database.Transactor
	userRepo    contracts.UserRepository
	validator   validator.Service
	jwt         jwt.Service
	googleOAuth oauth.GoogleOAuthInterface
	uuid        pkgUUID.Service
}

func NewAuthService(
	transactor database.Transactor,
	userRepo contracts.UserRepository,
	validator validator.Service,
	jwt jwt.Service,
	googleOAuth oauth.GoogleOAuthInterface,
	uuid pkgUUID.Service,
) contracts.AuthService {
	return &authService{
		transactor:  transactor,
		userRepo:    userRepo,
		validator:   validator,
		jwt:         jwt,
		googleOAuth: googleOAuth,
		uuid:        uuid,
	}
}
