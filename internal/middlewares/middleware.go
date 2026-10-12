package middlewares

import (
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

type Middleware struct {
	jwt jwt.Service
}

func NewMiddleware(jwt jwt.Service) *Middleware {
	return &Middleware{
		jwt: jwt,
	}
}
