package email

import (
	"strings"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sendgrid/sendgrid-go"
)

// ExtractEmailUsername extracts the username portion from an email address.
// It takes an email address as input and returns the username.
func ExtractEmailUsername(email string) string {
	at := strings.LastIndex(email, "@")
	username := email[:at]
	return username
}

// Send sends an email using the provided SendGrid client.
// It takes the SendGrid client, sender information, recipient information,
// subject, plain text content, and HTML content as input.
// It returns an error if the email sending fails.
func Send(sendGridClient *sendgrid.Client, sendFromName, sendFromAddress, subject, sendToName, sendToAddress, plainTextContent, htmlContent string) error {
	// Prepare email
	from := mail.NewEmail(sendFromName, sendFromAddress)
	to := mail.NewEmail(sendToName, sendToAddress)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	// Send email
	_, err := sendGridClient.Send(message)

	return err
}