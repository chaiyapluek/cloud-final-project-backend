package middleware

import (
	"log"
	"net/http"
	"time"

	"dev.chaiyapluek.cloud.final.backend/src/dto"
	appError "dev.chaiyapluek.cloud.final.backend/src/pkg/error"
	"github.com/labstack/echo/v4"
)

func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.Println("recover", r)
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Code:      http.StatusInternalServerError,
					Message:   http.StatusText(http.StatusInternalServerError),
					Path:      c.Request().URL.Path,
					Timestamp: time.Now(),
				})
				return
			}
		}()
		err := next(c)
		if err != nil {
			log.Println("error", err)
			switch err := err.(type) {
			case *echo.HTTPError:
				var errMsg string
				errMsg, ok := err.Message.(string)
				if !ok {
					errMsg = http.StatusText(err.Code)
				}
				return c.JSON(err.Code, dto.ErrorResponse{
					Code:      err.Code,
					Message:   errMsg,
					Path:      c.Request().URL.Path,
					Timestamp: time.Now(),
				})
			case *appError.AppError:
				return c.JSON(err.Code, dto.ErrorResponse{
					Code:      err.Code,
					Message:   err.Message,
					Path:      c.Request().URL.Path,
					Timestamp: time.Now(),
				})
			default:
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Code:      http.StatusInternalServerError,
					Message:   "any error",
					Path:      c.Request().URL.Path,
					Timestamp: time.Now(),
				})
			}
		}
		return nil
	}
}
