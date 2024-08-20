package helper

import (
	"bytes"
	"html/template"

	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper/email"
	"gopkg.in/gomail.v2"
)

func SendEmailInviteUser(payload dto.EmailDataHtmlDto, recipientEmail string, configAppPass string) error {
	tmpl, err := template.New("emailTemplate").Parse(email.EmailTemplate)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, payload); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", constant.BusinessHypayEmail)
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "Welcome to Hypay!")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, constant.BusinessHypayEmail, configAppPass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
