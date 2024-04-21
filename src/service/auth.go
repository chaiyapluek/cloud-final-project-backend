package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"dev.chaiyapluek.cloud.final.backend/src/dto"
	"dev.chaiyapluek.cloud.final.backend/src/entity"
	appError "dev.chaiyapluek.cloud.final.backend/src/pkg/error"
	"dev.chaiyapluek.cloud.final.backend/src/repository"
	"dev.chaiyapluek.cloud.final.backend/template"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	GetMe(userId string) (*dto.UserResponse, error)
	Login(attemptId string, email string, code string, requestAt time.Time) (*dto.UserResponse, error)
	LoginAttempt(email string, password string) (*dto.LoginAttemptResponse, error)
	Register(attemptId string, email string, code string, requestAt time.Time) (*dto.UserResponse, error)
	RegisterAttempt(email string, password string, name string) (*dto.RegisterAttemptResponse, error)
}

type authService struct {
	emailService EmailService
	authRepo     repository.AuthRepository
	userRepo     repository.UserRepository
}

func NewAuthService(emailService EmailService, authRepo repository.AuthRepository, userRepo repository.UserRepository) AuthService {
	return &authService{
		emailService: emailService,
		authRepo:     authRepo,
		userRepo:     userRepo,
	}
}

func generateCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10))
	}
	return code
}

func (s *authService) GetMe(userId string) (*dto.UserResponse, error) {
	oUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid user id format")
	}

	user, err := s.userRepo.GetById(&oUserId)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}
	if user == nil {
		return nil, appError.NewErrUnauthorized("user not found")
	}

	return &dto.UserResponse{
		UserId: user.Id.Hex(),
		Name:   user.Name,
	}, nil
}

func (s *authService) Login(attemptId string, email string, code string, requestAt time.Time) (*dto.UserResponse, error) {

	oAttemptId, err := primitive.ObjectIDFromHex(attemptId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid attempt id format")
	}

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}
	if user == nil {
		return nil, appError.NewErrUnauthorized("user not found")
	}

	attempt, err := s.authRepo.GetLoginAttemptById(&oAttemptId)
	if err != nil {
		return nil, err
	}
	if attempt == nil || attempt.Email != email || attempt.Code != code {
		return nil, appError.NewErrUnauthorized("invalid credential")
	}
	if requestAt.After(attempt.Expire) {
		return nil, appError.NewErrUnauthorized("code expire")
	}

	return &dto.UserResponse{
		UserId: user.Id.Hex(),
		Name:   user.Name,
	}, nil
}

func (s *authService) LoginAttempt(email string, password string) (*dto.LoginAttemptResponse, error) {

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}
	if user == nil {
		return nil, appError.NewErrUnauthorized("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println("compare password error", err)
		return nil, appError.NewErrUnauthorized("invalid credential")
	}

	code := generateCode()

	formattedCode := ""
	for _, c := range code {
		formattedCode += string(c) + " "
	}
	writer := bytes.NewBufferString("")
	err = template.Body(formattedCode).Render(context.Background(), writer)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	err = s.emailService.SendFromDefaultSender(email, "SAYWUB Login Code", writer.String())
	if err != nil {
		log.Println("send mail error", err)
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	newAttempt := entity.LoginAttempt{
		Email:  email,
		Code:   code,
		Expire: time.Now().Add(time.Minute * time.Duration(2)),
	}
	err = s.authRepo.CreateLoginAttempt(&newAttempt)
	if err != nil {
		log.Println("create login attempt", err)
		return nil, appError.NewErrInternalServerError(err.Error())
	}
	return &dto.LoginAttemptResponse{
		AttemptId: newAttempt.Id.Hex(),
		Code:      newAttempt.Code,
	}, nil
}

func (s *authService) Register(attemptId string, email string, code string, requestAt time.Time) (*dto.UserResponse, error) {

	oAttemptId, err := primitive.ObjectIDFromHex(attemptId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid attempt id format")
	}

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}
	if user != nil {
		return nil, appError.NewErrBadRequest("email already exist")
	}

	attempt, err := s.authRepo.GetRegisterAttemptById(&oAttemptId)
	if err != nil {
		return nil, err
	}
	if attempt == nil || attempt.Email != email || attempt.Code != code {
		return nil, appError.NewErrUnauthorized("invalid credential")
	}
	if requestAt.After(attempt.Expire) {
		return nil, appError.NewErrUnauthorized("code expire")
	}

	newUser := entity.User{
		Email:    email,
		Password: attempt.Password,
		Name:     attempt.Name,
	}
	err = s.userRepo.CreateUser(&newUser)
	if err != nil {
		log.Println("create user", err)
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	return &dto.UserResponse{
		UserId: newUser.Id.Hex(),
		Name:   newUser.Name,
	}, nil
}

func (s *authService) RegisterAttempt(email string, password string, name string) (*dto.RegisterAttemptResponse, error) {

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}
	if user != nil {
		return nil, appError.NewErrBadRequest("email already exist")
	}

	code := generateCode()

	formattedCode := ""
	for _, c := range code {
		formattedCode += string(c) + " "
	}
	writer := bytes.NewBufferString("")
	err = template.Body(formattedCode).Render(context.Background(), writer)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	err = s.emailService.SendFromDefaultSender(email, "SAYWUB Register Code", writer.String())
	if err != nil {
		log.Println("send mail error", err)
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("hash password error", err)
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	newAttempt := entity.RegisterAttempt{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
		Code:     code,
		Expire:   time.Now().Add(time.Minute * time.Duration(2)),
	}
	err = s.authRepo.CreateRegisterAttempt(&newAttempt)
	if err != nil {
		log.Println("create register attempt", err)
		return nil, appError.NewErrInternalServerError(err.Error())
	}
	return &dto.RegisterAttemptResponse{
		AttemptId: newAttempt.Id.Hex(),
		Code:      newAttempt.Code,
	}, nil
}
