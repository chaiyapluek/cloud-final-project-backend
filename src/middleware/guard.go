package middleware

import (
	appError "dev.chaiyapluek.cloud.final.backend/src/pkg/error"
	"github.com/labstack/echo/v4"
)

type guardMiddleware struct {
	frontendAPIKey string
}

func NewGuardMiddleware(frontendAPIKey string) *guardMiddleware {
	return &guardMiddleware{
		frontendAPIKey: frontendAPIKey,
	}
}

func (m *guardMiddleware) Guard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKey := c.Request().Header.Get("x-api-key")
		if apiKey != m.frontendAPIKey {
			return appError.NewErrUnauthorized("invalid api key")
		}
		return next(c)
	}
}
