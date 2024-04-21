package handler

import (
	"dev.chaiyapluek.cloud.final.backend/src/dto"
	appError "dev.chaiyapluek.cloud.final.backend/src/pkg/error"
	"dev.chaiyapluek.cloud.final.backend/src/service"
	"github.com/labstack/echo/v4"
)

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

func (h *authHandler) GetMe(e echo.Context) error {
	userId := e.Param("userId")
	resp, err := h.authService.GetMe(userId)
	if err != nil {
		return err
	}
	return e.JSON(200, dto.NewSuccessResponse(200, resp))
}

func (h *authHandler) Login(e echo.Context) error {
	var req dto.LoginRequest
	if err := e.Bind(&req); err != nil {
		return appError.NewErrBadRequest()
	}
	if req.RequestAt == nil {
		return appError.NewErrBadRequest("requestAt not exist")
	}
	resp, err := h.authService.Login(req.AttemptId, req.Email, req.Code, *req.RequestAt)
	if err != nil {
		return err
	}
	return e.JSON(200, dto.NewSuccessResponse(200, resp, "Login success"))
}

func (h *authHandler) LoginAttempt(e echo.Context) error {
	var req dto.LoginAttemptRequest
	if err := e.Bind(&req); err != nil {
		return appError.NewErrBadRequest()
	}
	if req.Email == "" {
		return appError.NewErrBadRequest()
	}
	resp, err := h.authService.LoginAttempt(req.Email, req.Password)
	if err != nil {
		return err
	}
	return e.JSON(200, dto.NewSuccessResponse(200, resp, "Login code send"))
}

func (h *authHandler) RegisterAttempt(e echo.Context) error {
	var req dto.RegisterAttemptRequest
	if err := e.Bind(&req); err != nil {
		return appError.NewErrBadRequest()
	}
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return appError.NewErrBadRequest()
	}
	resp, err := h.authService.RegisterAttempt(req.Email, req.Password, req.Name)
	if err != nil {
		return err
	}
	return e.JSON(200, dto.NewSuccessResponse(200, resp, "Register code send"))
}

func (h *authHandler) Register(e echo.Context) error {
	var req dto.RegisterRequest
	if err := e.Bind(&req); err != nil {
		return appError.NewErrBadRequest()
	}
	if req.RequestAt == nil {
		return appError.NewErrBadRequest("requestAt not exist")
	}
	resp, err := h.authService.Register(req.AttemptId, req.Email, req.Code, *req.RequestAt)
	if err != nil {
		return err
	}
	return e.JSON(200, dto.NewSuccessResponse(200, resp, "Register success"))
}
