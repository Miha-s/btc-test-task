package email_sender

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/logger"

	"crypto/tls"

	gomail "gopkg.in/gomail.v2"
)

type EmailSenderImpl struct {
	email   string
	subject string
	dialer  gomail.Dialer
}

func (sender *EmailSenderImpl) Init(conf *config.Config) error {
	sender.dialer = *gomail.NewDialer(conf.EmailServiceUrl, conf.EmailServicePort,
		conf.EmailToSendFrom, conf.EmailToSendFromPassword)

	sender.email = conf.EmailToSendFrom
	sender.subject = conf.EmailSubject
	sender.dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sender.subject = conf.EmailSubject
	return nil
}

func (sender *EmailSenderImpl) SendEmail(recipient, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", sender.email)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", sender.subject)
	message.SetBody("text/plain", body)

	if err := sender.dialer.DialAndSend(message); err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

func (sender *EmailSenderImpl) BroadcastEmails(recipients *map[string]struct{}, body string) error {
	for email := range *recipients {
		sender.SendEmail(email, body)
	}
	return nil
}
