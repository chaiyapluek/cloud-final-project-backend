package handler

import (
	"log"
	"strconv"

	"dev.chaiyapluek.cloud.final.backend/src/dto"
	appError "dev.chaiyapluek.cloud.final.backend/src/pkg/error"
	"dev.chaiyapluek.cloud.final.backend/src/service"
	"github.com/labstack/echo/v4"
)

type cartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *cartHandler {
	return &cartHandler{
		cartService: cartService,
	}
}

func (h *cartHandler) GetCartByUserId(e echo.Context) error {
	userId := e.Param("userId")
	locationId := e.QueryParam("locationId")

	if userId == "" {
		return appError.NewErrBadRequest("user id is required")
	}

	if locationId == "" {
		return appError.NewErrBadRequest("location id query param is required")
	}

	resp, err := h.cartService.GetCartByUserId(userId, locationId)
	if err != nil {
		return err
	}

	return e.JSON(200, dto.NewSuccessResponse(200, resp))
}

func (h *cartHandler) CreateCart(e echo.Context) error {
	var req dto.CreateCartRequest
	if err := e.Bind(&req); err != nil {
		return appError.NewErrBadRequest(err.Error())
	}

	cartId, err := h.cartService.CreateCart(req.UserId, req.LocationId, req.SessionId)
	if err != nil {
		return err
	}

	return e.JSON(201, dto.NewSuccessResponse(201, dto.CreateCartResponse{CartId: cartId}))
}

func (h *cartHandler) AddCartItem(e echo.Context) error {
	cartId := e.Param("cartId")
	if cartId == "" {
		return appError.NewErrBadRequest("cart id is required")
	}
	var req dto.AddCartItemRequest
	if err := e.Bind(&req); err != nil {
		return appError.NewErrBadRequest(err.Error())
	}

	resp, err := h.cartService.AddCartItem(cartId, &req)
	if err != nil {
		return err
	}

	return e.JSON(200, dto.NewSuccessResponse(200, resp))
}

func (h *cartHandler) DeleteCartItem(e echo.Context) error {
	cartId := e.Param("cartId")
	itemIdStr := e.Param("itemId")
	if cartId == "" {
		return appError.NewErrBadRequest("cart id is required")
	}
	if itemIdStr == "" {
		return appError.NewErrBadRequest("item id is required")
	}

	itemId, err := strconv.Atoi(itemIdStr)
	if err != nil {
		return appError.NewErrBadRequest("invalid item id format")
	}

	resp, err := h.cartService.DeleteCartItem(cartId, itemId)
	if err != nil {
		return err
	}

	return e.JSON(200, dto.NewSuccessResponse(200, resp))
}

func (h *cartHandler) Checkout(e echo.Context) error {
	var req dto.CheckoutRequest
	if err := e.Bind(&req); err != nil {
		log.Println(err)
		return appError.NewErrBadRequest(err.Error())
	}

	if req.CartId == "" || req.UserId == "" || req.Address == "" {
		return appError.NewErrBadRequest("cart id, user id, and address are required")
	}
	log.Println("checkout", req)
	err := h.cartService.Checkout(&req)
	if err != nil {
		return err
	}

	return e.JSON(200, dto.NewSuccessResponse(200, nil))
}
