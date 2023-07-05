// EMAIL SERVICE

func (s *Server) newSendGridClient() {
	s.sendGridClient = sendgrid.NewSendClient(s.cfg.Email.SendGrid.ApiKey)
}