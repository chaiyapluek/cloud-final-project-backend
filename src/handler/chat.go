package handler

import (
	"dev.chaiyapluek.cloud.final.backend/src/dto"
	appError "dev.chaiyapluek.cloud.final.backend/src/pkg/error"
	"dev.chaiyapluek.cloud.final.backend/src/service"
	"github.com/labstack/echo/v4"
)

type chatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *chatHandler {
	return &chatHandler{
		chatService: chatService,
	}
}

func (h *chatHandler) GetUserChat(e echo.Context) error {
	userId := e.Param("userId")
	locationId := e.QueryParam("locationId")

	if userId == "" || locationId == "" {
		return appError.NewErrBadRequest("userId and locationId is required")
	}

	chats, err := h.chatService.GetUserChat(userId, locationId)
	if err != nil {
		return err
	}

	return e.JSON(200, dto.NewSuccessResponse(200, chats))
}

func (h *chatHandler) Sent(e echo.Context) error {
	var req dto.SendChatRequest
	if err := e.Bind(&req); err != nil {
		return appError.NewErrBadRequest(err.Error())
	}

	if req.UserId == "" || req.LocationId == "" || req.Content == "" {
		return appError.NewErrBadRequest("userId, locationId and content is required")
	}

	resp, err := h.chatService.Sent(req.UserId, req.LocationId, req.Content)
	if err != nil {
		return err
	}

	return e.JSON(200, dto.NewSuccessResponse(200, resp))
}
