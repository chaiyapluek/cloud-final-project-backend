package service

import (
	"bytes"
	"context"

	"dev.chaiyapluek.cloud.final.backend/src/dto"
	"dev.chaiyapluek.cloud.final.backend/src/entity"
	appError "dev.chaiyapluek.cloud.final.backend/src/pkg/error"
	"dev.chaiyapluek.cloud.final.backend/src/repository"
	"dev.chaiyapluek.cloud.final.backend/template"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartService interface {
	GetCartByUserId(userId, locationId string) (*dto.CartResponse, error)
	CreateCart(userId, locationId, sessionId *string) (string, error)
	AddCartItem(cid string, req *dto.AddCartItemRequest) (*dto.AddCartItemResponse, error)
	DeleteCartItem(cartId string, itemId int) (*dto.DeleteCartItemResponse, error)
	Checkout(req *dto.CheckoutRequest) error
}

type cartServiceImpl struct {
	cartRepo     repository.CartRepository
	userRepo     repository.UserRepository
	chatRepo     repository.ChatRepository
	emailService EmailService
}

func NewCartService(cartRepo repository.CartRepository, userRepo repository.UserRepository, chatRepo repository.ChatRepository, emailService EmailService) CartService {
	return &cartServiceImpl{
		cartRepo:     cartRepo,
		userRepo:     userRepo,
		chatRepo:     chatRepo,
		emailService: emailService,
	}
}

func (s *cartServiceImpl) GetCartByUserId(userId, locationId string) (*dto.CartResponse, error) {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid user id format")
	}

	oLocationId, err := primitive.ObjectIDFromHex(locationId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid location id format")
	}

	cart, err := s.cartRepo.GetCartByUserId(&oid, &oLocationId)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}
	if cart == nil {
		cartId, err := s.CreateCart(&userId, &locationId, nil)
		if err != nil {
			return nil, appError.NewErrInternalServerError(err.Error())
		}
		return &dto.CartResponse{
			CartId:     cartId,
			UserId:     userId,
			LocationId: locationId,
			Items:      []*dto.CartItem{},
		}, nil
	}

	cartItemResp := []*dto.CartItem{}
	for idx := range cart.Items {
		if cart.Items[idx] == nil || cart.Items[idx].MenuId == nil {
			continue
		}
		steps := []*dto.ItemStep{}
		for _, step := range cart.Items[idx].Steps {
			steps = append(steps, &dto.ItemStep{
				Step:    step.Step,
				Options: step.Options,
			})
		}
		cartItemResp = append(cartItemResp, &dto.CartItem{
			MenuId:     cart.Items[idx].MenuId.Hex(),
			MenuName:   cart.Items[idx].MenuName,
			ItemId:     cart.Items[idx].ItemId,
			Quantity:   cart.Items[idx].Quantity,
			TotalPrice: cart.Items[idx].TotalPrice,
			Steps:      steps,
		})
	}

	return &dto.CartResponse{
		CartId:       cart.Id.Hex(),
		UserId:       userId,
		LocationId:   locationId,
		LocationName: cart.LocationName,
		Items:        cartItemResp,
	}, nil
}

func (s *cartServiceImpl) CreateCart(userId, locationId, sessionId *string) (string, error) {
	var oUserId *primitive.ObjectID
	if userId != nil {
		oid, err := primitive.ObjectIDFromHex(*userId)
		if err != nil {
			return "", appError.NewErrBadRequest("invalid user id format")
		}
		oUserId = &oid
	}

	var oLocationId *primitive.ObjectID
	if locationId != nil {
		oid, err := primitive.ObjectIDFromHex(*locationId)
		if err != nil {
			return "", appError.NewErrBadRequest("invalid location id format")
		}
		oLocationId = &oid
	}

	if oUserId != nil {
		newCart := &entity.Cart{
			UserId:     oUserId,
			LocationId: oLocationId,
			Items:      []*entity.CartItem{},
		}
		err := s.cartRepo.CreateCart(newCart)
		if err != nil {
			return "", appError.NewErrInternalServerError(err.Error())
		}
		return newCart.Id.Hex(), nil
	} else if sessionId != nil {
		// create cart with session id
		newCart := &entity.Cart{
			SessionId:  sessionId,
			LocationId: oLocationId,
			Items:      []*entity.CartItem{},
		}
		err := s.cartRepo.CreateCart(newCart)
		if err != nil {
			return "", appError.NewErrInternalServerError(err.Error())
		}
		return newCart.Id.Hex(), nil
	} else {
		return "", appError.NewErrBadRequest("either user id or session id must be provided	")
	}
}

func (s *cartServiceImpl) AddCartItem(cid string, req *dto.AddCartItemRequest) (*dto.AddCartItemResponse, error) {
	cartId, err := primitive.ObjectIDFromHex(cid)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid cart id format")
	}

	cart, err := s.cartRepo.GetByCartId(&cartId)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	if cart == nil {
		return nil, appError.NewErrBadRequest("cart not found")
	}

	menuId, err := primitive.ObjectIDFromHex(req.MenuId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid menu id format")
	}

	maxItemId := 0
	for _, v := range cart.Items {
		if v.ItemId > maxItemId {
			maxItemId = v.ItemId
		}
	}

	cartItem := &entity.CartItem{
		MenuId:     &menuId,
		ItemId:     maxItemId + 1,
		Quantity:   req.Quantity,
		TotalPrice: req.TotalPrice,
	}

	for _, step := range req.Steps {
		cartItem.Steps = append(cartItem.Steps, entity.ItemStep{
			Step:    step.Step,
			Options: step.Options,
		})
	}

	cart, err = s.cartRepo.AddCartItem(&cartId, cartItem)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	return &dto.AddCartItemResponse{
		CartId:    cartId.Hex(),
		ItemId:    cartItem.ItemId,
		TotalItem: len(cart.Items),
	}, nil
}

func (s *cartServiceImpl) DeleteCartItem(cartId string, itemId int) (*dto.DeleteCartItemResponse, error) {
	cid, err := primitive.ObjectIDFromHex(cartId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid cart id format")
	}

	cart, err := s.cartRepo.DeleteCartItem(&cid, itemId)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	totalPrice := 0
	for _, v := range cart.Items {
		totalPrice += v.TotalPrice * v.Quantity
	}

	return &dto.DeleteCartItemResponse{
		CartId:     cartId,
		TotalItem:  len(cart.Items),
		TotalPrice: totalPrice,
	}, nil
}

func (s *cartServiceImpl) Checkout(req *dto.CheckoutRequest) error {
	cartId, err := primitive.ObjectIDFromHex(req.CartId)
	if err != nil {
		return appError.NewErrBadRequest("invalid cart id format")
	}

	userId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return appError.NewErrBadRequest("invalid user id format")
	}

	cart, err := s.cartRepo.GetByCartId(&cartId)
	if err != nil {
		return appError.NewErrInternalServerError(err.Error())
	}
	if cart == nil {
		return appError.NewErrBadRequest("cart not found")
	}

	user, err := s.userRepo.GetById(&userId)
	if err != nil {
		return appError.NewErrInternalServerError(err.Error())
	}
	if user == nil {
		return appError.NewErrBadRequest("user not found")
	}

	// validate cart
	if cart.UserId != nil && *cart.UserId != userId {
		return appError.NewErrBadRequest("cart not belong to user")
	}

	cart, err = s.cartRepo.GetCartByUserId(&userId, cart.LocationId)
	if err != nil {
		return appError.NewErrInternalServerError(err.Error())
	}

	totalPrice := 0
	for _, v := range cart.Items {
		totalPrice += v.TotalPrice * v.Quantity
	}

	items := []*template.ReceiptItem{}
	for _, v := range cart.Items {
		steps := []*template.Step{}
		for _, step := range v.Steps {
			steps = append(steps, &template.Step{
				Name:    step.Step,
				Options: step.Options,
			})
		}
		items = append(items, &template.ReceiptItem{
			MenuName:   v.MenuName,
			Quantity:   v.Quantity,
			TotalPrice: v.TotalPrice,
			Steps:      steps,
		})
	}

	writer := bytes.NewBufferString("")
	err = template.Receipt(template.ReceiptProps{
		Name:         user.Name,
		Address:      req.Address,
		LocationName: cart.LocationName,
		TotalPrice:   totalPrice,
		Items:        items,
	}).Render(context.Background(), writer)
	if err != nil {
		return appError.NewErrInternalServerError(err.Error())
	}

	err = s.emailService.SendFromDefaultSender(user.Email, "SAYWUB Receipt", writer.String())
	if err != nil {
		return appError.NewErrInternalServerError(err.Error())
	}

	err = s.cartRepo.DeleteCartById(&cartId)
	if err != nil {
		return appError.NewErrInternalServerError(err.Error())
	}

	err = s.chatRepo.DeleteChat(&userId, cart.LocationId)
	if err != nil {
		return appError.NewErrInternalServerError(err.Error())
	}

	return nil
}
