package service

import (
	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailService interface {
	SendFromDefaultSender(to string, subject string, body string) error
	Send(from, to string, subject string, body string) error
}

type emailServiceImpl struct {
	server *mail.SMTPServer
	sender string
}

func NewEmailService(server *mail.SMTPServer, sender string) EmailService {
	return &emailServiceImpl{
		server: server,
		sender: sender,
	}
}

func (s *emailServiceImpl) SendFromDefaultSender(to string, subject string, body string) error {
	return s.Send(s.sender, to, subject, body)
}

func (s *emailServiceImpl) Send(from, to string, subject string, body string) error {
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
	return nil
}
