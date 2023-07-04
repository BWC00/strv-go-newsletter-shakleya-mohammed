package config

type EmailConfig struct {
	SendGrid SendGrid
}

type SendGrid struct {
	SendFromName    string `env:"SEND_FROM_NAME,required"`
	SendFromAddress string `env:"SEND_FROM_ADDRESS,required"`
	ApiKey  string `env:"SENDGRID_API_KEY,required"`
}