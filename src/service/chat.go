package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"dev.chaiyapluek.cloud.final.backend/src/dto"
	"dev.chaiyapluek.cloud.final.backend/src/entity"
	appError "dev.chaiyapluek.cloud.final.backend/src/pkg/error"
	"dev.chaiyapluek.cloud.final.backend/src/repository"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/bedrock"
	"github.com/tmc/langchaingo/prompts"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService interface {
	GetUserChat(userId string, locationId string) ([]*dto.ChatResponse, error)
	Sent(userId string, locationId string, content string) ([]*dto.ChatResponse, error)
}

type chatService struct {
	chatRepo     repository.ChatRepository
	locationRepo repository.LocationRepository
	llm          *bedrock.LLM
}

func NewChatService(chatRepo repository.ChatRepository, locationRepo repository.LocationRepository, llm *bedrock.LLM) ChatService {
	return &chatService{
		chatRepo:     chatRepo,
		locationRepo: locationRepo,
		llm:          llm,
	}
}

func extractInfo(raw string) (*entity.AIAnswer, string) {
	idx1 := strings.Index(raw, "{")
	idx2 := strings.Index(raw, "}")
	if idx1 == -1 || idx2 == -1 {
		return nil, ""
	}
	raw = raw[idx1 : idx2+1]
	var aiAnswer entity.AIAnswer
	err := json.Unmarshal([]byte(raw), &aiAnswer)
	if err != nil {
		return nil, ""
	}

	return &aiAnswer, raw
}

func findOption(str string, options []*entity.Option) []string {
	result := []string{}
	for _, v := range options {
		// regex case insensitive
		search := strings.ReplaceAll(v.Name, "+", "\\+")
		re := regexp.MustCompile("(?i)" + search)
		if re.MatchString(str) {
			result = append(result, v.Value)
		}
	}
	return result
}

func extractItemInfoAndEncode(raw string, steps []*entity.Step) (string, error) {
	m := map[string][]string{}
	for idx, v := range steps {
		switch idx {
		case 0:
			m["step-0"] = findOption(raw, v.Options)
		case 1:
			m["step-1"] = findOption(raw, v.Options)
		case 2:
			m["step-2"] = findOption(raw, v.Options)
		case 3:
			m["step-3"] = findOption(raw, v.Options)
		case 4:
			m["step-4"] = findOption(raw, v.Options)
		case 5:
			m["step-5"] = findOption(raw, v.Options)
		case 6:
			m["step-6"] = findOption(raw, v.Options)
		}
	}

	anyOnMap := false
	for _, v := range m {
		if len(v) > 0 {
			anyOnMap = true
			break
		}
	}

	if !anyOnMap {
		return "", nil
	}

	// base64 encode
	b, err := json.Marshal(m)
	if err != nil {
		return "", appError.NewErrInternalServerError(err.Error())
	}
	encoded := base64.StdEncoding.EncodeToString(b)

	return encoded, nil
}

func (s *chatService) GetUserChat(userId string, locationId string) ([]*dto.ChatResponse, error) {
	oUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid user id format")
	}

	oLocationId, err := primitive.ObjectIDFromHex(locationId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid location id format")
	}

	chats, err := s.chatRepo.GetUserChat(&oUserId, &oLocationId)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	resp := []*dto.ChatResponse{}
	for _, v := range chats {
		resp = append(resp, &dto.ChatResponse{
			Type:    v.Type,
			Sender:  v.Sender,
			Content: v.Content,
		})
	}

	return resp, nil
}

func (s *chatService) Sent(userId string, locationId string, content string) ([]*dto.ChatResponse, error) {

	oUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid user id format")
	}

	oLocationId, err := primitive.ObjectIDFromHex(locationId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid location id format")
	}

	chats, err := s.chatRepo.GetUserChat(&oUserId, &oLocationId)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	history := ""
	input := ""

	if len(chats) == 0 {
		input = fmt.Sprintf("%s [/INST]", content)
	} else {
		for _, v := range chats {
			if v.Sender == 0 && v.Type == 0 {
				history += v.RawContent
			} else {
				history += v.RawContent
			}
		}
		input = fmt.Sprintf("<s>[INST] %s (reponse no longer than 150 words and do not answer other than sandwich) [/INST]</s>", content)
	}

	userSendTime := time.Now()

	p := prompts.NewPromptTemplate(promptTemplate, []string{"history", "input"})
	c := chains.NewLLMChain(s.llm, p,
		chains.WithTemperature(0.25),
		chains.WithMaxTokens(256),
		chains.WithTopP(1),
	)
	out, err := chains.Call(context.Background(), c, map[string]interface{}{
		"history": history,
		"input":   input,
	})
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	_, menus, err := s.locationRepo.GetLocationById(&oLocationId)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	menu, err := s.locationRepo.GetMenuItmes(&oLocationId, menus[0].Id)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}
	if menu == nil {
		return nil, appError.NewErrInternalServerError("menu not found")
	}

	encoded, err := extractItemInfoAndEncode(out[c.OutputKey].(string), menu.Steps)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	err = s.chatRepo.InsertMany([]*entity.Chat{
		{
			LocationId: &oLocationId,
			UserId:     &oUserId,
			Type:       0,
			Sender:     1,
			RawContent: input,
			Content:    content,
			SentAt:     userSendTime,
		},
		{
			LocationId: &oLocationId,
			UserId:     &oUserId,
			Type:       0,
			Sender:     0,
			RawContent: fmt.Sprintf("%s </s>", out[c.OutputKey].(string)),
			Content:    out[c.OutputKey].(string),
			SentAt:     time.Now(),
		},
		{
			LocationId: &oLocationId,
			UserId:     &oUserId,
			Type:       1,
			Sender:     0,
			Content:    fmt.Sprintf("%s?preference=%s", menus[0].Id.Hex(), encoded),
			SentAt:     time.Now().Add(time.Second),
		},
	})
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	return []*dto.ChatResponse{
		{
			Type:    0,
			Sender:  1,
			Content: content,
		},
		{
			Type:    0,
			Sender:  0,
			Content: out[c.OutputKey].(string),
		},
		{
			Type:    1,
			Sender:  0,
			Content: fmt.Sprintf("%s?preference=%s", menus[0].Id.Hex(), encoded),
		},
	}, nil
}

var promptTemplate = `
<s>[INST] <<SYS>>
You are sandwich restaurant cashier in sandwich restaurant.
Your task is recommend sandwich order to user and calculate the price.
You are talkative and provides sandwich order recommendation. 
If you do not know the answer to a question, you truthfully say you do not know. 
If customer ask about something that out of sandwich order recommendation topic, you will say you cannot answer.
Price is in baht.
You will recommend order base on following step.
Your recommendation message have to complete 6 steps.
Here is 7 steps to create subway order.

Step 1. Choose bread.
"name": "bread"
"description": "Choose type of bread"
"type": "radio" (Choose 1)
"require": true
"items":
name, price
"Wheat", 0
"Honey Oat", 0
"Italian", 0
"Parmesan oregano", 0
"Flatbread", 0

Step 2. Choose bread size.
"name": "bread_size"
"description": "Choose bread size"
"type": "radio" (Choose 1)
"require": true
"items":
name, price
"Six", 0
"Footlong", 0

Step 3. Choose menu.
"name": "menu"
"description": "Choose favorite menu"
"type": "radio" (Choose 1)
"require": true
"items":
name, six_price, footlong_price
"Teriyaki Chicken", 145, 259
"Roasted Beef", 145, 259
"BBQ Chicken", 129, 229
"Roasted Chicken", 129, 229
"Veggie Delux", 109, 209
"Tuna", 99, 189
"Sliced Chicken", 99, 189
"Veggie Delite", 89, 169

Step 4. Choose vegetable.
"name": "vegetable"
"description": "Choose vegetable"
"type": "checkbox" (Choose 1 or more)
"require": true
"items":
name, price
"Lettuce", 0
"Tomatoes", 0
"Cucumbers", 0
"Pickles", 0
"Green Peppers", 0
"Olives", 0
"Onion", 0
"Jelapenos", 0

Step 5. Choose sauce.
"name": "sauce"
"description": "Choose sauce"
"type": "checkbox" (Choose between 0 to 3)
"require": true
"items":
name, price
"Honey Mustard", 0
"Sweet Onion", 0
"Chipotle Southwest", 0
"Mayonnaise", 0
"BBQ Sauce", 0
"Tomato Sauce", 0
"Thousand Island Dressing", 0
"Hot chilli sauce", 0

Step 6. Choose add-ons.
"name": "add-on"
"description": "Choose your add-ons"
"type": "checkbox" (Choose 0 or more)
"require": false
"items":
name, price
"Mozzarella Cheese", 15
"Cheddar Cheese", 15
"Bacon", 40
"Egg", 15
"Double Meat", 60
"Avocado", 40
"Chopped Mushroom", 20

Step 7. Choose meal.
"name": "meal"
"description": "Make it a meal"
"type": "checkbox" (Choose 1)
"require": false
"items":
name, price
"No meal", 0
"Cookie + 22 Oz. Pepsi", 59

Recommendation template:
Your recommendation...
Summary:
Choose type of bread: your recommendation about type of bread\n
Choose bread size: your recommendation about bread size\n
Choose menu: your recommendation about menu\n
Choose vegetables: your recommendation about vegetables\n
Choose sauces: your recommendation about sauces\n
Meal: your recommendation about meal\n
Total price: total price


Your recommendation is not longer than 100 words.
Your recommendation message have to complete 6 steps. 
You must pick something to your customer. Not ask customer to choose.
Make each step you reccomend as bullet points.
Only generate for Cashier part. Don't generate Customer part.
<</SYS>>
{{.history}}
{{.input}}
`
