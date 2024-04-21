package handler

import (
	"dev.chaiyapluek.cloud.final.backend/src/dto"
	"dev.chaiyapluek.cloud.final.backend/src/service"
	"github.com/labstack/echo/v4"

	appError "dev.chaiyapluek.cloud.final.backend/src/pkg/error"
)

type locationHandler struct {
	locationService service.LocationService
}

func NewLocationHandler(locationService service.LocationService) *locationHandler {
	return &locationHandler{
		locationService: locationService,
	}
}

func (h *locationHandler) GetAllLocation(c echo.Context) error {
	resp, err := h.locationService.GetAllLocation()
	if err != nil {
		return err
	}

	return c.JSON(200, dto.NewSuccessResponse(200, resp))
}

func (h *locationHandler) GetLocationById(c echo.Context) error {
	idParam := c.Param("id")

	if idParam == "" {
		return appError.NewErrBadRequest("location id is required")
	}
	resp, err := h.locationService.GetLocationById(idParam)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.NewSuccessResponse(200, resp))
}

func (h *locationHandler) GetLocationMenu(c echo.Context) error {
	locationId := c.Param("locationId")
	menuId := c.Param("menuId")

	if locationId == "" || menuId == "" {
		return appError.NewErrBadRequest("location id and menu id is required")
	}

	resp, err := h.locationService.GetLocationMenu(locationId, menuId)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.NewSuccessResponse(200, resp))
}
