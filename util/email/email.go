package email

import (
	"strings"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sendgrid/sendgrid-go"
)

func ExtractEmailUsername(email string) string {
	at := strings.LastIndex(email, "@")
	username := email[:at]
	return username
}

func Send(sendGridClient *sendgrid.Client, sendFromName, sendFromAddress, subject, sendToName, sendToAddress, plainTextContent, htmlContent string) error {
	// Prepare email
	from := mail.NewEmail(sendFromName, sendFromAddress)
	to := mail.NewEmail(sendToName, sendToAddress)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	// Send email
	_, err := sendGridClient.Send(message)

	return err
}