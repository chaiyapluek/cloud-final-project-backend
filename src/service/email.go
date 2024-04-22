package service

import (
	"time"

	"dev.chaiyapluek.cloud.final.backend/src/entity"
	appError "dev.chaiyapluek.cloud.final.backend/src/pkg/error"
	"dev.chaiyapluek.cloud.final.backend/src/repository"
	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailService interface {
	SendFromDefaultSender(to string, subject string, body string) error
	Send(from, to string, subject string, body string) error
}

type emailServiceImpl struct {
	emailRepo       repository.EmailRepository
	limitSendPerDay int
	server          *mail.SMTPServer
	sender          string
}

func NewEmailService(emailRepo repository.EmailRepository, limitSendPerDay int, server *mail.SMTPServer, sender string) EmailService {
	return &emailServiceImpl{
		emailRepo:       emailRepo,
		limitSendPerDay: limitSendPerDay,
		server:          server,
		sender:          sender,
	}
}

func (s *emailServiceImpl) SendFromDefaultSender(to string, subject string, body string) error {
	return s.Send(s.sender, to, subject, body)
}

func (s *emailServiceImpl) Send(from, to string, subject string, body string) error {

	number, err := s.emailRepo.GetNumberOfEmailSendWithInADay(to)
	if err != nil {
		return err
	}
	if number >= s.limitSendPerDay {
		return appError.NewErrUnprocessableEntity("limit send email per day")
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject).
		SetBody(mail.TextHTML, body)
	if email.Error != nil {
		return email.Error
	}
	client, err := s.server.Connect()
	if err != nil {
		return err
	}
	err = email.Send(client)
	if err != nil {
		return err
	}

	s.emailRepo.Save(&entity.Email{
		To:     to,
		SendAt: time.Now(),
	})

	return nil
}
