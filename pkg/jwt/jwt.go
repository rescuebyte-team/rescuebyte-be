package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/env"
)

var ErrInvalidTokenType = errors.New("invalid token type")

type Service interface {
	Create(id uuid.UUID, name string, email string, role enums.UserRole) (string, error)
	CreateRefreshToken(id uuid.UUID, name string, email string, role enums.UserRole) (string, error)
	Decode(tokenString string, claims *Claims) error
	DecodeRefreshToken(tokenString string, claims *Claims) error
}

type Claims struct {
	jwt.RegisteredClaims
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Role      enums.UserRole `json:"role"`
	TokenType string         `json:"token_type,omitempty"`
}

type jwtImpl struct {
	SecretKey      string
	ExpiredTime    time.Duration
	RefreshExpTime time.Duration
}

var JWT = getJWT()

func getJWT() Service {
	return &jwtImpl{
		SecretKey:      env.AppEnv.JwtSecretKey,
		ExpiredTime:    env.AppEnv.JwtExpTime,
		RefreshExpTime: env.AppEnv.JwtRefreshExpTime,
	}
}

func (j *jwtImpl) generateToken(
	id uuid.UUID,
	name string,
	email string,
	role enums.UserRole,
	tokenType string,
	duration time.Duration,
) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "rescuebyte-be",
			Subject:   id.String(),
			Audience:  jwt.ClaimStrings{"rescuebyte-be"},
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		ID:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *jwtImpl) Create(id uuid.UUID, name string, email string, role enums.UserRole) (string, error) {
	return j.generateToken(id, name, email, role, "access", j.ExpiredTime)
}

func (j *jwtImpl) CreateRefreshToken(id uuid.UUID, name string, email string, role enums.UserRole) (string, error) {
	return j.generateToken(id, name, email, role, "refresh", j.RefreshExpTime)
}

func (j *jwtImpl) Decode(tokenString string, claims *Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	},
		jwt.WithIssuer("rescuebyte-be"),
		jwt.WithAudience("rescuebyte-be"),
	)

	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}

func (j *jwtImpl) DecodeRefreshToken(tokenString string, claims *Claims) error {
	if err := j.Decode(tokenString, claims); err != nil {
		return err
	}

	if claims.TokenType != "refresh" {
		return ErrInvalidTokenType
	}

	return nil
}
