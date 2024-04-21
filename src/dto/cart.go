package dto

type CreateCartRequest struct {
	UserId     *string `json:"userId"`
	LocationId *string `json:"locationId"`
	SessionId  *string `json:"sessionId"`
}

type CreateCartResponse struct {
	CartId string `json:"cartId"`
}

type AddCartItemRequest struct {
	MenuId     string      `json:"menuId"`
	Quantity   int         `json:"quantity"`
	TotalPrice int         `json:"totalPrice"`
	Steps      []*ItemStep `json:"steps"`
}

type AddCartItemResponse struct {
	CartId    string `json:"cartId"`
	ItemId    int    `json:"itemId"`
	TotalItem int    `json:"totalItem"`
}

type ItemStep struct {
	Step    string   `json:"step"`
	Options []string `json:"options"`
}

type CartItem struct {
	MenuId     string      `json:"menuId"`
	MenuName   string      `json:"menuName"`
	ItemId     int         `json:"itemId"`
	Quantity   int         `json:"quantity"`
	TotalPrice int         `json:"totalPrice"`
	Steps      []*ItemStep `json:"steps"`
}

type CartResponse struct {
	CartId       string      `json:"cartId"`
	UserId       string      `json:"userId"`
	LocationId   string      `json:"locationId"`
	LocationName string      `json:"locationName"`
	Items        []*CartItem `json:"items"`
}

type DeleteCartItemResponse struct {
	CartId     string `json:"cartId"`
	TotalItem  int    `json:"totalItem"`
	TotalPrice int    `json:"totalPrice"`
}

type CheckoutRequest struct {
	CartId  string `json:"cartId"`
	UserId  string `json:"userId"`
	Address string `json:"address"`
}
