package dto

type ChatResponse struct {
	Type    int    `json:"type"`
	Sender  int    `json:"sender"`
	Content string `json:"content"`
}

type SendChatRequest struct {
	UserId     string `json:"userId"`
	LocationId string `json:"locationId"`
	Content    string `json:"content"`
}
