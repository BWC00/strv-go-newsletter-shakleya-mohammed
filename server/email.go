package server

import (
	"github.com/sendgrid/sendgrid-go"
)

// EMAIL SERVICE

// newSendGridClient initializes a new SendGrid client for email services.
// It creates a new SendGrid client using the API key provided in the server's configuration settings.
// This client can be used to send emails using the SendGrid service.
func (s *Server) newSendGridClient() {
	s.sendGridClient = sendgrid.NewSendClient(s.cfg.Email.SendGrid.ApiKey)
}